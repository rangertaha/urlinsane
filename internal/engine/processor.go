// Copyright (C) 2024 Rangertaha
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package engine

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/domain"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		cfg config.Config

		// Domain
		target internal.Domain

		// Metrics
		progress *progressbar.ProgressBar
		started  time.Time
		elapsed  time.Duration
		total    int64
		live     int64
	}
)

// NewUrlinsane ...
func New(conf config.Config) (u Urlinsane) {
	return Urlinsane{
		total:    0,
		cfg:      conf,
		started:  time.Now(),
		progress: progressbar.DefaultSilent(1000),
	}
}

// Init
func (u *Urlinsane) Init() <-chan internal.Domain {
	out := make(chan internal.Domain)

	// Create application directory for the target domains
	// used to store files and images we collect
	if err := u.Mkdir(u.cfg.Target()); err != nil {
		log.Error("Creating target directory", err)
		// time.Sleep(1 * time.Second)
	}

	go func() {
		// Initialize database plugins if needed
		if db, ok := u.cfg.Database().(internal.Initializer); ok {
			log.Debug("Init database:", u.cfg.Database().Id())
			db.Init(&u.cfg)
		}

		// Initialize collector plugins if needed
		log.Debug("Collectors:", len(u.cfg.Collectors()))
		for _, info := range u.cfg.Collectors() {
			if inf, ok := info.(internal.Initializer); ok {
				log.Debug("Init collector:", info.Id())
				inf.Init(&u.cfg)
			}
		}

		// Initialize algorithm plugins if needed
		log.Debug("Algorithms:", len(u.cfg.Algorithms()))
		for _, algorithm := range u.cfg.Algorithms() {
			if al, ok := algorithm.(internal.Initializer); ok {
				log.Debug("Init algorithm: ", algorithm.Id())
				al.Init(&u.cfg)
			}
		}

		// Initialize analyzer plugins if needed
		log.Debug("Analyzers:", len(u.cfg.Analyzers()))
		for _, alz := range u.cfg.Analyzers() {
			if anz, ok := alz.(internal.Initializer); ok {
				log.Debug("Init analyzer:", alz.Id())
				anz.Init(&u.cfg)
			}
		}

		// Initialize output plugin if needed
		if out, ok := u.cfg.Output().(internal.Initializer); ok {
			log.Debug("Init output: ", u.cfg.Output().Id())
			out.Init(&u.cfg)
		}

		// Send original domain down to get data collected about it
		if original := domain.New(u.cfg.Target()); original.Valid() {
			log.Debug("Target: ", original.String())
			out <- original
		} else {
			log.Debugf("domain %s not valid.", original.String())
			u.Close()
		}

		if u.cfg.Banner() {
			log.Debug("Show banner !")
			Banner(u.cfg)
		}

		close(out)
	}()
	return out
}

// Target collects the same info on the target domain
func (u *Urlinsane) Target(in <-chan internal.Domain) <-chan internal.Domain {
	out := make(chan internal.Domain)

	go func() {
		// Collect info on target domain
		u.target = <-u.CollectorChain(u.cfg.Collectors(), in)
		log.Debug("Target domain: ", u.target.String())

		// Print report of target domain
		log.Debug("Generate domain report: ", u.target.String())

		// Initialize algorithm plugins if needed
		for _, algorithm := range u.cfg.Algorithms() {
			out <- domain.NewVariant(algorithm, u.cfg.Target())
		}

		close(out)
	}()

	return out
}

// Algorithms generate typo variations using the algorithm plugins
func (u *Urlinsane) Algorithms(in <-chan internal.Domain) <-chan internal.Domain {
	if len(u.cfg.Algorithms()) > 0 {
		out := make(chan internal.Domain)
		var wg sync.WaitGroup

		for w := 1; w <= u.cfg.Workers(); w++ {
			wg.Add(1)
			go func(id int, in <-chan internal.Domain, out chan<- internal.Domain) {
				defer wg.Done()
				for domain := range in {

					acc := NewAccumulator(out)
					if err := domain.Algorithm().Exec(u.target, acc); err != nil {
						log.Error("Algorithm failed: ", err)
					}
				}
			}(w, in, out)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return u.PreFilters(out)
	}
	return u.PreFilters(in)
}

func (u *Urlinsane) PreFilters(in <-chan internal.Domain) <-chan internal.Domain {
	out := make(chan internal.Domain)
	variants := make(map[string]bool)

	go func() {
		for typo := range in {
			// Removing duplicates
			if _, ok := variants[typo.String()]; !ok {
				variants[typo.String()] = true

				// Make sure the variant does not match the original
				if typo.String() != u.target.String() {

					// Set Levenshtein distance
					//   https://en.wikipedia.org/wiki/Levenshtein_distance
					dist := fuzzy.Levenshtein(typo.String(), u.target.String())
					typo.Ld(dist)

					if dist >= u.cfg.Distance() {
						out <- typo
					}
				}
			}
		}

		// Update domain count
		u.total = int64(len(variants))

		// Show optional progress bar
		if u.cfg.Progress() {
			u.progress = progressbar.Default(u.total)
		}
		close(out)
	}()

	return out
}

func (u *Urlinsane) Collectors(in <-chan internal.Domain) <-chan internal.Domain {
	if len(u.cfg.Collectors()) > 0 {
		out := make(chan internal.Domain)
		var wg sync.WaitGroup

		for w := 1; w <= u.cfg.Workers(); w++ {
			wg.Add(1)
			go func(in <-chan internal.Domain, out chan<- internal.Domain) {
				defer wg.Done()
				for c := range u.CollectorChain(u.cfg.Collectors(), in) {
					out <- c
				}
			}(in, out)
		}
		go func() {
			wg.Wait()
			close(out)
		}()
		return u.PostFilters(out)
	}
	log.Debug("No collectors !")
	return u.PostFilters(in)
}

func (u *Urlinsane) PostFilters(in <-chan internal.Domain) <-chan internal.Domain {
	out := make(chan internal.Domain)

	go func() {
		for typo := range in {
			if typo.Live() {
				u.live++
			}
			if typo.Live() && u.cfg.Registered() {
				out <- typo
			} else if !typo.Live() && u.cfg.Unregistered() {
				out <- typo
			} else if !u.cfg.Registered() && !u.cfg.Unregistered() {
				out <- typo
			}
		}
		close(out)
	}()

	return out
}

// InfoChain creates a chain of information-gathering functions
func (u *Urlinsane) CollectorChain(funcs []internal.Collector, in <-chan internal.Domain) <-chan internal.Domain {
	if len(funcs) == 0 {
		log.Debug("No collectors to chain !")
		return in
	}
	var xfunc internal.Collector
	out := make(chan internal.Domain)
	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			if fn, ok := xfunc.(internal.Initializer); ok {
				fn.Init(&u.cfg)
			}
			// Timing options
			time.Sleep(u.cfg.Random() * u.cfg.Delay())

			// Execute the collector and timeout if it takes too long
			ctx, cancel := context.WithTimeout(context.Background(), u.cfg.Timeout())
			acc := NewAccumulator(out)
			u.runner(ctx, xfunc, i, acc)
			cancel()
		}
		close(out)
	}()

	if len(funcs) > 0 {
		return u.CollectorChain(funcs, out)
	}

	return out
}

func (u *Urlinsane) runner(ctx context.Context, fn internal.Collector, domain internal.Domain, acc internal.Accumulator) {
	logger := log.WithFields(log.Fields{"c": fn.Id(), "d": domain.String()})
	fn.Exec(domain, acc)
	select {
	case <-time.After(1 * time.Second):
		logger.Info("Collector completed")
	case <-ctx.Done():
		logger.Error("Collector timed out:", ctx.Err())
	}
}

func (u *Urlinsane) Analyzers(in <-chan internal.Domain) <-chan internal.Domain {

	return in
}

// Progress adds a visible progesssbar if -p flag is set
func (u *Urlinsane) Progress(in <-chan internal.Domain) <-chan internal.Domain {
	if u.cfg.Progress() {
		out := make(chan internal.Domain)
		go func(in <-chan internal.Domain, out chan<- internal.Domain) {
			for t := range in {
				u.progress.Add(1)
				out <- t
			}
			// Clear/hide the progress bar after all typos have passed through
			u.progress.Clear()
			close(out)

		}(in, out)
		return out
	}
	log.Debug("No progress bar !")
	return in
}

func (u *Urlinsane) Output(in <-chan internal.Domain) {
	output := u.cfg.Output()

	// Send domain typos to the output plugin
	for c := range in {
		output.Read(c)
	}

	// Writes output if it can't stream
	output.Write()

	// Save typo records collected by the output plugin
	if fname := u.cfg.File(); fname != "" {
		output.Save(fname)
	}

	// Print summary
	if u.cfg.Summary() {
		u.elapsed = time.Since(u.started)
		summary := map[string]string{
			"  TIME:":                             u.elapsed.String(),
			text.FgGreen.Sprintf("%s", "  LIVE:"): fmt.Sprintf("%d", u.live),
			text.FgRed.Sprintf("%s", "  OFFLINE"): fmt.Sprintf("%d", u.total-u.live),
			"  TOTAL:":                            fmt.Sprintf("%d", u.total),
		}
		output.Summary(summary)
	}
}

func (u *Urlinsane) Close() {
	// Initialize information plugins if needed
	for _, info := range u.cfg.Collectors() {
		if inf, ok := info.(internal.Closer); ok {
			inf.Close()
		}
	}

	// Close db
	u.cfg.Database().Close()
}

func (u *Urlinsane) Mkdir(dir string) (err error) {

	return
}

func (u *Urlinsane) Execute() (err error) {
	typos := u.Init()
	typos = u.Target(typos)
	typos = u.Algorithms(typos)
	typos = u.Collectors(typos)
	typos = u.Analyzers(typos)
	typos = u.Progress(typos)
	u.Output(typos)
	u.Close()

	return
}

func (u *Urlinsane) Stream() <-chan internal.Domain {
	typos := u.Init()
	typos = u.Algorithms(typos)
	typos = u.Collectors(typos)
	typos = u.Analyzers(typos)
	return typos
}

func Banner(cfg config.Config) {
	var lang, board, algo []string
	t := time.Now()
	timestamp := t.Format("2006-01-02 15:04:05")
	name := text.FgRed.Sprint(cfg.Target())
	for _, l := range cfg.Languages() {
		lang = append(lang, l.Id())
	}
	for _, b := range cfg.Keyboards() {
		board = append(board, b.Id())
	}
	for _, a := range cfg.Algorithms() {
		algo = append(algo, a.Id())
	}
	fmt.Printf(
		internal.BANNER,
		internal.VERSION,
		name,
		strings.Join(lang, ","),
		strings.Join(board, ","),
		strings.Join(algo, ","),
		cfg.Output().Id(),
		timestamp,
	)
}
