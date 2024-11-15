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
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		cfg *config.Config

		// Domain
		target internal.Domain

		// Metrics
		progress *progressbar.ProgressBar
		started  time.Time
		elapsed  time.Duration
		total    int64
		live     int64
	}
	FilterFunc func() func(<-chan internal.Domain, *config.Config) <-chan internal.Domain
)

// NewUrlinsane ...
func New(conf *config.Config) (u *Urlinsane) {
	return &Urlinsane{
		total:    0,
		cfg:      conf,
		started:  time.Now(),
		progress: progressbar.DefaultSilent(1000),
	}
}

// Init
func (u *Urlinsane) Init() <-chan internal.Domain {
	out := make(chan internal.Domain)

	u.target = domain.New(u.cfg.Target())

	log := log.WithFields(
		log.Fields{"fqdn": u.target.String(), "prefix": u.target.Prefix(),
			"name": u.target.Name(), "suffix": u.target.Suffix()})

	// Initialize database plugins if needed
	if db, ok := u.cfg.Database().(internal.Initializer); ok {
		log.Debug("Init database:", u.cfg.Database().Id())
		db.Init(u.cfg)
	}

	// Initialize collector plugins if needed
	log.Debug("Collectors:", len(u.cfg.Collectors()))
	for _, info := range u.cfg.Collectors() {
		if inf, ok := info.(internal.Initializer); ok {
			log.Debug("Init collector:", info.Id())
			inf.Init(u.cfg)
		}
	}

	// Initialize algorithm plugins if needed
	log.Debug("Algorithms:", len(u.cfg.Algorithms()))
	for _, algorithm := range u.cfg.Algorithms() {
		if al, ok := algorithm.(internal.Initializer); ok {
			log.Debug("Init algorithm: ", algorithm.Id())
			al.Init(u.cfg)
		}
	}

	// Initialize analyzer plugins if needed
	log.Debug("Analyzers:", len(u.cfg.Analyzers()))
	for _, alz := range u.cfg.Analyzers() {
		if anz, ok := alz.(internal.Initializer); ok {
			log.Debug("Init analyzer:", alz.Id())
			anz.Init(u.cfg)
		}
	}

	// Initialize output plugin if needed
	if out, ok := u.cfg.Output().(internal.Initializer); ok {
		log.Debug("Init output: ", u.cfg.Output().Id())
		out.Init(u.cfg)
	}

	// Send original domain down to get data collected about it
	go func() {
		if u.target.Valid() {
			log.Debug("Target: ", u.target.String())
			out <- u.target
		} else {
			log.Debugf("domain %s not valid.", u.target.String())
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

		for origin := range in {
			// Send origianl domain
			out <- origin

			// Send a domain for each algorithm to apply
			for _, algorithm := range u.cfg.Algorithms() {
				log.Debugf("Adding %s algo to %s for typo gen", algorithm.Id(), origin.String())
				out <- domain.Variant(algorithm, origin.String(), origin.String())
			}
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
					if domain.Algorithm() != nil {
						acc := NewAccumulator(out, u.target, u.cfg)
						log.Debugf("Executing %s algo", domain.Algorithm().Id())
						if err := domain.Algorithm().Exec(u.target, acc); err != nil {
							log.Error("Algorithm failed: ", err)
						}
					} else {
						out <- domain
					}
				}
			}(w, in, out)
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
	return in
}

func (u *Urlinsane) Constraints(in <-chan internal.Domain, Filters ...FilterFunc) <-chan internal.Domain {
	for _, fn := range Filters {
		in = fn()(in, u.cfg)
	}
	// // Show optional progress bar
	// if u.cfg.Progress() {
	// 	u.progress = progressbar.Default(int64(u.cfg.Count()))
	// }
	return in
}

// func (u *Urlinsane) PreFilters(in <-chan internal.Domain) <-chan internal.Domain {
// 	out := make(chan internal.Domain)
// 	variants := make(map[string]bool)

// 	go func() {
// 		for typo := range in {
// 			// Removing duplicates
// 			if _, ok := variants[typo.String()]; !ok {
// 				variants[typo.String()] = true

// 				// Make sure the variant does not match the original
// 				if typo != nil {

// 					// Set Levenshtein distance
// 					//   https://en.wikipedia.org/wiki/Levenshtein_distance
// 					dist := fuzzy.Levenshtein(typo.String(), u.target.String())
// 					logger := log.WithFields(log.Fields{"min": u.cfg.Distance(), "dist": dist, "d": typo.String()})
// 					logger.Debug("Levenshtein distance")
// 					typo.Ld(dist)

// 					if dist <= u.cfg.Distance() {
// 						out <- typo
// 					} else {
// 						logger.Warn("Does not have the levenshtein minimum distance")
// 					}
// 				}
// 			}
// 		}

// 		// Update domain count
// 		u.total = int64(len(variants))

// 		// Show optional progress bar
// 		if u.cfg.Progress() {
// 			u.progress = progressbar.Default(u.total)
// 		}
// 		close(out)
// 	}()

// 	return out
// }

// func (u *Urlinsane) ReadCache(in <-chan internal.Domain) <-chan internal.Domain {
// 	out := make(chan internal.Domain)

// 	go func() {
// 		for domain := range in {
// 			log.Debugf("Reading from cache: %s", domain.String())
// 			if data, _ := u.cfg.Database().Read(domain.String()); data != "" {
// 				domain.Json(data)
// 				log.Debug("From cache", data)
// 			}
// 			out <- domain
// 		}
// 		close(out)
// 	}()

// 	return out
// }

// func (u *Urlinsane) Collectors(in <-chan internal.Domain) <-chan internal.Domain {
// 	out := make(chan internal.Domain)

// 	go func() {
// 		for domain := range in {
// 			if domain.Cached() {
// 				return u.PostFilters(in)
// 			}
// 			return

// 		}
// 		close(out)
// 	}()

//		return u.CollectorWorkers(in)
//	}
func (u *Urlinsane) Collectors(in <-chan internal.Domain) <-chan internal.Domain {
	if len(u.cfg.Collectors()) > 0 {
		out := make(chan internal.Domain)
		var wg sync.WaitGroup

		for w := 1; w <= u.cfg.Workers(); w++ {
			wg.Add(1)
			go func(in <-chan internal.Domain, out chan<- internal.Domain) {
				defer wg.Done()
				for c := range u.CollectorChain(u.cfg.Collectors(), in) {
					log.Debugf("Collection chain completed for %s", c.String())
					out <- c
				}
			}(in, out)
		}
		go func() {
			wg.Wait()
			close(out)
		}()
		return out
	}
	log.Debug("No collectors !")
	return in
}

func (u *Urlinsane) WriteCache(in <-chan internal.Domain) <-chan internal.Domain {
	out := make(chan internal.Domain)

	go func() {
		kv := u.cfg.Database()
		for domain := range in {
			if !domain.Cached() {
				json := domain.Json()
				if err := kv.Write(json, domain.String()); err != nil {
					log.Error(err)
				}
				log.Debug("Saved to cache:", json)
			}

			out <- domain
		}
		close(out)
	}()

	return out
}

func (u *Urlinsane) PostFilters(in <-chan internal.Domain) <-chan internal.Domain {
	out := make(chan internal.Domain)
	logger := log.WithFields(log.Fields{"reg": u.cfg.Registered(), "unreg": u.cfg.Unregistered()})

	go func() {
		for typo := range in {
			if typo.Live() {
				u.live++
			}
			logger = logger.WithFields(log.Fields{"d": typo.String(), "live": typo.Live(), "ld": typo.Ld()})
			logger.Debug("Filtering domain")

			if u.cfg.Registered() {
				if typo.Live() {
					logger.Debug("Filter allowing registered")
					out <- typo
				}
			}

			if u.cfg.Unregistered() {
				if !typo.Live() {
					logger.Debug("Filter allowing unregistered")
					out <- typo
				}
			}

			if u.cfg.Unregistered() == u.cfg.Registered() {
				logger.Debug("Filter allowing all")
				out <- typo

			}
		}
		close(out)
	}()

	return out
}

// CollectorChain creates a chain of information-gathering functions
func (u *Urlinsane) CollectorChain(funcs []internal.Collector, in <-chan internal.Domain) <-chan internal.Domain {
	if len(funcs) == 0 {
		log.Debug("No collectors to chain !")
		return in
	}
	var xfunc internal.Collector
	out := make(chan internal.Domain)
	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for variant := range in {
			if fn, ok := xfunc.(internal.Initializer); ok {
				fn.Init(u.cfg)
			}
			// Timing options
			time.Sleep(u.cfg.Random() * u.cfg.Delay())

			var ctx context.Context
			var cancel context.CancelFunc
			// Execute the collector and timeout if it takes too long
			if u.cfg.Timeout() != 0 {
				ctx, cancel = context.WithTimeout(context.Background(), u.cfg.Timeout())
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}

			// acc := NewAccumulator(out)
			u.runner(ctx, xfunc, variant, out)
			cancel()
		}
		close(out)
	}()

	if len(funcs) > 0 {
		return u.CollectorChain(funcs, out)
	}

	return out
}

func (u *Urlinsane) runner(ctx context.Context, fn internal.Collector, domain internal.Domain, out chan internal.Domain) {
	logger := log.WithFields(log.Fields{"c": fn.Id(), "d": domain.String()})

	if domain.Cached() {
		out <- domain
	} else {
		fn.Exec(NewAccumulator(out, domain, u.cfg))
	}

	select {
	case <-time.After(1 * time.Second):
		logger.Infof("Collector %s completed", fn.Id())
	case <-ctx.Done():
		logger.Error("Collector timed out:", ctx.Err())
		out <- domain
	}
}

func (u *Urlinsane) Analyzers(in <-chan internal.Domain) <-chan internal.Domain {
	if len(u.cfg.Analyzers()) == 0 {
		log.Debug("No analyzers to run !")
		return in
	}
	return in
}

// Progress adds a visible progesssbar if -p flag is set
func (u *Urlinsane) Progress(in <-chan internal.Domain) <-chan internal.Domain {
	if u.cfg.Progress() {
		log.Debug("Total count ", u.cfg.Count())

		out := make(chan internal.Domain)
		go func(in <-chan internal.Domain, out chan<- internal.Domain) {
			// Show optional progress bar
			u.progress = progressbar.Default(int64(u.cfg.Count()))
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
	logger := log.WithFields(log.Fields{"output": u.cfg.Output().Id()})
	output := u.cfg.Output()

	for c := range in {
		logger.WithFields(log.Fields{"d": c.String()}).
			Debug("Sending to output")

		// Send domain typos to the output plugin
		if c.Algorithm() != nil {
			output.Read(c)
		}

		// Save domain to directory
		if u.cfg.AssetDir() != "" {
			if err := c.Save(u.cfg.AssetDir()); err != nil {
				logger.Error("Saving: ", err)
			}
		}
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

	// os.Exit(1)
}

func (u *Urlinsane) Execute() (err error) {
	typos := u.Init()
	typos = u.Target(typos)
	typos = u.Algorithms(typos)
	typos = u.Constraints(typos,
		DedupFilter,
		RegexFilter,
		LevenshteinFilter,
		ReadCacheFilter,
	)
	typos = u.Collectors(typos)
	typos = u.WriteCache(typos)
	typos = u.PostFilters(typos)
	typos = u.Analyzers(typos)
	typos = u.Progress(typos)
	u.Output(typos)
	u.Close()

	return
}

func (u *Urlinsane) Stream() <-chan internal.Domain {
	typos := u.Init()
	typos = u.Target(typos)
	typos = u.Algorithms(typos)
	typos = u.Constraints(typos,
		DedupFilter,
		RegexFilter,
		LevenshteinFilter,
		ReadCacheFilter,
	)
	typos = u.Collectors(typos)
	typos = u.WriteCache(typos)
	typos = u.PostFilters(typos)
	typos = u.Analyzers(typos)
	return typos
}

func Banner(cfg *config.Config) {
	var lang, board, algo, collectors []string
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
	for _, c := range cfg.Collectors() {
		collectors = append(collectors, c.Id())
	}
	fmt.Printf(
		internal.BANNER,
		internal.VERSION,
		name,
		strings.Join(lang, ","),
		strings.Join(board, ","),
		strings.Join(algo, ","),
		strings.Join(collectors, ","),
		cfg.Output().Id(),
		timestamp,
	)
}
