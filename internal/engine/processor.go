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
	"sync"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/domain"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		cfg config.Config

		// Domains
		target internal.Domain

		// Metrics
		progress *progressbar.ProgressBar
		started  time.Time
		elapsed  time.Duration
		total    int64
		// online   int64
		// filtered int64
		// scanned  int64
	}
)

// NewUrlinsane ...
func New(conf config.Config) (u Urlinsane) {
	return Urlinsane{
		total:    0,
		cfg:      conf,
		started:  time.Now(),
		progress: progressbar.DefaultSilent(0),
	}
}

// Init
func (u *Urlinsane) Init() <-chan internal.Domain {
	out := make(chan internal.Domain)

	// Create application directory for the target domains
	// used to store files and images we collect
	if err := u.Mkdir(u.cfg.Target()); err != nil {
		log.Error("Creating target directory", err)
		time.Sleep(1 * time.Second)
	}

	go func() {

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

		// Initialize database plugins if needed
		if db, ok := u.cfg.Database().(internal.Initializer); ok {
			log.Debug("Init database:", u.cfg.Database().Id())
			db.Init(&u.cfg)
		}

		if u.cfg.Banner() {
			log.Debug("Show banner !")
			internal.Banner(u.cfg.Target())
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
	out := make(chan internal.Domain)
	var wg sync.WaitGroup

	for w := 1; w <= u.cfg.Concurrency(); w++ {
		wg.Add(1)
		go func(id int, in <-chan internal.Domain, out chan<- internal.Domain) {
			defer wg.Done()
			for domain := range in {

				acc := NewAccumulator(out)
				domain.Algorithm().Exec(u.target, acc)
			}
		}(w, in, out)
	}

	go func() {
		wg.Wait()

		// Update total after all algorithms complete procducing typos
		// u.total = int64(len(u.typos))
		// if u.cfg.Progress() {
		// 	u.progress = progressbar.Default(u.total)
		// }
		close(out)
	}()

	return out
}

func (u *Urlinsane) Filters(in <-chan internal.Domain) <-chan internal.Domain {
	out := make(chan internal.Domain)
	variants := make(map[string]bool)

	go func() {
		for typo := range in {
			// Removing duplicates
			if _, ok := variants[typo.String()]; !ok {
				variants[typo.String()] = true

				// Make sure the variant does not match the original
				if typo.String() != u.target.String() {
					out <- typo
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

		for w := 1; w <= u.cfg.Concurrency(); w++ {
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
		return out
	}
	log.Debug("No collectors !")
	return in
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
	logger := log.WithFields(log.Fields{"collector": fn.Id(), "domain": domain.String()})
	fn.Exec(domain, acc)
	select {
	case <-time.After(1 * time.Second):
		logger.Info("Function completed successfully")
	case <-ctx.Done():
		logger.Error("Function timed out:", ctx.Err())
	}
}

// Progress adds a visible progesssbar if -p flag is set
func (u *Urlinsane) Progress(in <-chan internal.Domain) <-chan internal.Domain {
	if u.cfg.Progress() {
		out := make(chan internal.Domain)
		go func(in <-chan internal.Domain, out chan<- internal.Domain) {
			for t := range in {
				u.progress.Add(1)
				out <- t

				log.WithFields(log.Fields{
					"domain": t.String(),
				}).Debug("Progress(<-)")

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

	// Stream typo records to the output plugin
	for c := range in {
		u.cfg.Output().Write(c)
		// if u.cfg.ShowAll() {
		// 	u.cfg.Output().Write(c)

		// } else if vari.Live {
		// 	u.cfg.Output().Write(c)
		// }

	}

	// Save typo records collected by the output plugin
	u.cfg.Output().Save()

	// Print summary
	u.elapsed = time.Since(u.started)
	summary := map[string]int{
		"ELAPSED": int(u.elapsed),
		"TOTAL":   int(u.total),
	}
	u.cfg.Output().Summary(summary)

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
	// time.Sleep(5 * time.Second)

	return
}

func (u *Urlinsane) Execute() (err error) {
	typos := u.Init()
	typos = u.Target(typos)
	typos = u.Algorithms(typos)
	typos = u.Filters(typos)
	typos = u.Collectors(typos)
	typos = u.Progress(typos)
	u.Output(typos)
	u.Close()

	return
}

// func (u *Urlinsane) Stream() <-chan internal.Domain {
// 	typos := u.Init()
// 	typos = u.Algorithms(typos)
// 	typos = u.Filters(typos)
// 	typos = u.Collectors(typos)
// 	return typos
// }
