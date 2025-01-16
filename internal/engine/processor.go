// Copyright 2024 Rangertaha. All Rights Reserved.
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
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		cfg *config.Config

		// Domain
		target db.Domain
		scan   db.Scan

		// Metrics
		progress *progressbar.ProgressBar
		started  time.Time
		elapsed  time.Duration
		total    int64
		live     int64
	}
	FilterFunc func() func(<-chan *db.Domain, *config.Config) <-chan *db.Domain
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
func (u *Urlinsane) Init() <-chan *db.Domain {
	out := make(chan *db.Domain)

	// u.target = &db.Domain{Name: u.cfg.Target()}
	db.DB.FirstOrInit(&u.target, db.Domain{Name: u.cfg.Target()})
	db.DB.Preload("Results").FirstOrInit(&u.scan, db.Scan{Query: u.target.Name})

	log := log.WithFields(
		log.Fields{"domain": u.target.Name})

	// // Initialize database plugins if needed
	// if db, ok := u.cfg.Database().(internal.Initializer); ok {
	// 	log.Debug("Init database:", u.cfg.Database().Id())
	// 	db.Init(u.cfg)
	// }

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

	go func() {
		// Send original domain
		out <- &u.target

		if u.cfg.Banner() {
			log.Debug("Show banner !")
			Banner(u.cfg)
		}

		close(out)
	}()
	return out
}

// Algorithms generate typo variations using the algorithm plugins
func (u *Urlinsane) Algorithms(in <-chan *db.Domain) <-chan *db.Domain {
	if len(u.cfg.Algorithms()) > 0 {
		out := make(chan *db.Domain)
		var wg sync.WaitGroup

		for domain := range in {
			for _, algo := range u.cfg.Algorithms() {
				wg.Add(1)
				go func(algo internal.Algorithm, in <-chan *db.Domain, out chan<- *db.Domain) {
					defer wg.Done()

					domains, err := algo.Exec(domain)
					if err != nil {
						log.Errorf("Algorithm %s failed: %s", algo.Name(), err.Error())
					}
					for _, domain := range domains {
						out <- domain
					}

				}(algo, in, out)
			}
		}

		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
	return in
}

// Constrains apply pre-processing filters to exclude domain names from processing
func (u *Urlinsane) Constraints(in <-chan *db.Domain, Filters ...FilterFunc) <-chan *db.Domain {
	for _, fn := range Filters {
		in = fn()(in, u.cfg)
	}
	return u.Load(in)
}

// Load gets the domain from the database
func (u *Urlinsane) Load(in <-chan *db.Domain) <-chan *db.Domain {
	out := make(chan *db.Domain)
	go func() {
		for d := range in {
			var domain *db.Domain
			result := db.DB.Preload("Dns").Preload("IPs").Preload("Redirect").FirstOrInit(&domain, db.Domain{Name: d.Name})
			if result.Error != nil {
				log.Errorf("Loading %s failed: %s", d.Name, result.Error.Error())
			}

			out <- d
		}
		close(out)
	}()

	return out
}

func (u *Urlinsane) Collectors(in <-chan *db.Domain) <-chan *db.Domain {
	if len(u.cfg.Collectors()) > 0 {
		out := make(chan *db.Domain)
		var wg sync.WaitGroup

		for w := 1; w <= u.cfg.Workers(); w++ {
			wg.Add(1)
			go func(in <-chan *db.Domain, out chan<- *db.Domain) {
				defer wg.Done()
				for c := range u.CollectorChain(u.cfg.Collectors(), in) {
					log.Debugf("Collection chain completed for %s", c.Name)
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

// CollectorChain creates a chain of information-gathering functions
func (u *Urlinsane) CollectorChain(funcs []internal.Collector, in <-chan *db.Domain) <-chan *db.Domain {
	if len(funcs) == 0 {
		log.Debug("No collectors to chain !")
		return in
	}
	var xfunc internal.Collector
	out := make(chan *db.Domain)
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

func (u *Urlinsane) runner(ctx context.Context, fn internal.Collector, domain *db.Domain, out chan *db.Domain) {
	logger := log.WithFields(log.Fields{"c": fn.Id(), "d": domain.Name})

	var err error
	domain, err = fn.Exec(domain)
	if err != nil {
		logger.Errorf("Collector err: %s", err.Error())
	}

	select {
	case <-time.After(5 * time.Second):
		logger.Infof("Collector %s completed", fn.Id())
		out <- domain
	case <-ctx.Done():
		logger.Error("Collector timed out:", ctx.Err())
		out <- domain
	}
}

func (u *Urlinsane) Analyzers(in <-chan *db.Domain) <-chan *db.Domain {
	if len(u.cfg.Analyzers()) == 0 {
		log.Debug("No analyzers to run !")
		return in
	}
	return in
}

func (u *Urlinsane) Output(in <-chan *db.Domain) {
	output := u.cfg.Output()

	for c := range in {
		// Stream or collect domains
		output.Read(c)

		// Save domain to database
		if c.Live() {
			c.Save()
		}

		// Collect scan results
		u.scan.Results = append(u.scan.Results, c)
	}

	// Optionally, writes collected domains
	output.Write()

	// Optionally print summary
	if u.cfg.Summary() {
		output.Report()
	}

	// Save scans
	db.DB.Save(&u.scan)
}

func (u *Urlinsane) Close() {
	// Initialize information plugins if needed
	for _, info := range u.cfg.Collectors() {
		if inf, ok := info.(internal.Closer); ok {
			inf.Close()
		}
	}

	// Close db
	// u.cfg.Database().Close()

	// os.Exit(1)
}

func (u *Urlinsane) Execute() (err error) {
	typos := u.Init()
	typos = u.Algorithms(typos)
	typos = u.Constraints(typos, Dedup, Regex, Levenshtein)
	typos = u.Collectors(typos)
	typos = u.Analyzers(typos)
	u.Output(typos)
	u.Close()
	return
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
