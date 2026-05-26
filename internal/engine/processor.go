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
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/engine/dag"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		cfg internal.Config

		// cols is the dependency-ordered collector set actually executed
		// (includes any dependencies auto-included by the DAG). Set when the
		// collector stage runs; used by Close.
		cols []internal.Collector

		// Domain
		target db.Domain

		// Metrics
		progress *progressbar.ProgressBar
		started  time.Time
		elapsed  time.Duration
		total    int64
		live     int64
	}
	FilterFunc func() func(<-chan *db.Domain, internal.Config) <-chan *db.Domain
)

// NewUrlinsane ...
func New(conf internal.Config) (u *Urlinsane) {
	return &Urlinsane{
		total:    0,
		cfg:      conf,
		started:  time.Now(),
		progress: progressbar.DefaultSilent(1000),
	}
}

// Init seeds the pipeline with the original target domain and initializes the
// algorithm, analyzer and output plugins.
func (u *Urlinsane) Init(ctx context.Context) <-chan *db.Domain {
	out := make(chan *db.Domain)

	// u.target = &db.Domain{Name: u.cfg.Target()}
	// Seed the target, hydrating from the content-addressed store if a prior
	// scan cached it; otherwise start fresh.
	u.target = db.Domain{Name: u.cfg.Target(), Type: u.cfg.EntityType()}
	if s := u.cfg.Store(); s != nil {
		if cached, ok, err := s.Get(u.cfg.EntityType(), u.cfg.Target()); err == nil && ok {
			u.target = *cached
		}
	}

	log := log.WithFields(
		log.Fields{"domain": u.target.Name})

	// // Initialize database plugins if needed
	// if db, ok := u.cfg.Database().(internal.Initializer); ok {
	// 	log.Debug("Init database:", u.cfg.Database().Id())
	// 	db.Init(u.cfg)
	// }

	// Collector plugins are initialized in the collector stage (Collectors),
	// once the dependency DAG has resolved the full set to run.

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
		defer close(out)
		// Send original domain
		select {
		case out <- &u.target:
		case <-ctx.Done():
			return
		}

		if u.cfg.Banner() {
			log.Debug("Show banner !")
			Banner(u.cfg)
		}
	}()
	return out
}

// Algorithms generate typo variations using the algorithm plugins
func (u *Urlinsane) Algorithms(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
	if len(u.cfg.Algorithms()) > 0 {
		out := make(chan *db.Domain)
		var wg sync.WaitGroup

		for domain := range in {
			for _, algo := range u.cfg.Algorithms() {
				wg.Add(1)
				go func(algo internal.Algorithm, origin *db.Domain) {
					defer wg.Done()

					domains, err := algo.Exec(origin)
					if err != nil {
						log.Errorf("Algorithm %s failed: %s", algo.Name(), err.Error())
					}
					for _, variant := range domains {
						// Carry the origin so the Analyzers stage can pair
						// each variant with the domain it was derived from.
						variant.Origin = origin
						select {
						case out <- variant:
						case <-ctx.Done():
							return
						}
					}

				}(algo, domain)
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
func (u *Urlinsane) Constraints(ctx context.Context, in <-chan *db.Domain, Filters ...FilterFunc) <-chan *db.Domain {
	for _, fn := range Filters {
		in = fn()(in, u.cfg)
	}
	return u.Load(ctx, in)
}

// Load preloads each variant's cached record (Dns/IPs/Redirect) from the
// database so collectors can reuse it, and forwards the loaded record rather
// than the bare in-flight variant. Pipeline-only metadata that a DB load drops
// (Algorithm, Levenshtein, Origin) is carried over from the in-flight variant.
func (u *Urlinsane) Load(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
	out := make(chan *db.Domain)
	go func() {
		defer close(out)
		for d := range in {
			loaded := d // default: forward the in-flight variant
			if s := u.cfg.Store(); s != nil {
				if cached, ok, err := s.Get(d.EntityType(), d.Name); err != nil {
					log.Errorf("Loading %s failed: %s", d.Name, err.Error())
				} else if ok {
					// Reuse cached result data; carry pipeline-only metadata.
					cached.Algorithm = d.Algorithm
					cached.Levenshtein = d.Levenshtein
					cached.Origin = d.Origin
					loaded = cached
				}
			}

			select {
			case out <- loaded:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Collectors enriches each variant by running the collector plugins in
// dependency order. The execution order is a DAG built from every collector's
// declared Dependencies() (e.g. geo/ptr/wi/bn run after ip). Variants are
// processed in parallel across a worker pool; within a single variant the
// collectors run sequentially, level by level, so a collector always sees the
// results of the collectors it depends on — and no two collectors mutate the
// same variant concurrently.
func (u *Urlinsane) Collectors(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
	if len(u.cfg.Collectors()) == 0 {
		log.Debug("No collectors !")
		return in
	}

	// Build the dependency DAG. Dependencies the user did not explicitly
	// select are auto-included (closure) via collectorResolver.
	levels, err := dag.Levels(u.cfg.Collectors(), collectorResolver)
	if err != nil {
		log.Errorf("Collector dependency error: %s; skipping collection", err)
		return in
	}
	u.cols = dag.Flatten(levels)

	// Single Init pass over the full resolved set (fixes the previous
	// double-initialization and covers DAG-included dependencies).
	for _, c := range u.cols {
		if init, ok := c.(internal.Initializer); ok {
			log.Debug("Init collector:", c.Id())
			init.Init(u.cfg)
		}
	}

	out := make(chan *db.Domain)
	var wg sync.WaitGroup
	for w := 0; w < u.cfg.Workers(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for variant := range in {
				u.runCollectorLevels(ctx, levels, variant)
				select {
				case out <- variant:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// runCollectorLevels runs all collectors for one variant in dependency order:
// level by level, sequentially within the variant. Sequential execution keeps
// the shared *db.Domain race-free; parallelism comes from many variants in
// flight across the worker pool.
func (u *Urlinsane) runCollectorLevels(ctx context.Context, levels [][]internal.Collector, variant *db.Domain) {
	for _, level := range levels {
		for _, c := range level {
			if ctx.Err() != nil {
				return
			}
			u.execCollector(ctx, c, variant)
		}
	}
}

// execCollector runs a single collector against a single variant, applying the
// configured throttle and a per-collector timeout. The timeout-bearing context
// is passed to Exec so context-aware collectors honor the deadline.
func (u *Urlinsane) execCollector(ctx context.Context, c internal.Collector, variant *db.Domain) {
	logger := log.WithFields(log.Fields{"c": c.Id(), "d": variant.Name})

	// Optional throttle between collector calls.
	if d := u.cfg.Random() * u.cfg.Delay(); d > 0 {
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return
		}
	}

	cctx := ctx
	var cancel context.CancelFunc = func() {}
	if t := u.cfg.Timeout(); t > 0 {
		cctx, cancel = context.WithTimeout(ctx, t)
	}
	defer cancel()

	if _, err := c.Exec(cctx, variant); err != nil {
		logger.Errorf("collector err: %s", err.Error())
	}
	logger.Debugf("collector %s completed", c.Id())
}

// collectorResolver constructs a collector by id so the DAG can auto-include
// dependencies that were not explicitly selected.
func collectorResolver(id string) (internal.Collector, bool) {
	creator, err := collectors.Get(id)
	if err != nil {
		return nil, false
	}
	return creator(), true
}

// Analyzers runs the analyzer plugins on each variant, pairing it with the
// origin domain it was derived from (carried on variant.Origin). Analyzers
// enrich/score the variant in place.
func (u *Urlinsane) Analyzers(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
	if len(u.cfg.Analyzers()) == 0 {
		log.Debug("No analyzers to run !")
		return in
	}

	out := make(chan *db.Domain)
	go func() {
		defer close(out)
		for variant := range in {
			origin := variant.Origin
			if origin == nil {
				origin = &u.target
			}
			for _, az := range u.cfg.Analyzers() {
				if _, err := az.Exec(ctx, origin, variant); err != nil {
					log.Errorf("Analyzer %s failed for %s: %s", az.Id(), variant.Name, err.Error())
				}
			}
			select {
			case out <- variant:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func (u *Urlinsane) Output(ctx context.Context, in <-chan *db.Domain) {
	if output := u.cfg.Output(); output != nil {
		var live []*db.Domain
		for c := range in {
			// Stream or collect domains
			output.Read(c)

			// Collect live variants for persistence
			if c.Live() {
				live = append(live, c)
			}
		}

		// Optionally, writes collected domains
		output.Write()

		// Optionally print summary
		if u.cfg.Summary() {
			output.Report()
		}

		// Persist the scan to the content-addressed store (the source of truth).
		if s := u.cfg.Store(); s != nil {
			if _, err := s.PutScan(u.cfg.Target(), live); err != nil {
				log.Errorf("store PutScan failed: %s", err.Error())
			}
		}
	} else {
		fmt.Println("Invalid output formater")
	}
}

func (u *Urlinsane) Close() {
	cols := u.cols
	if cols == nil {
		cols = u.cfg.Collectors()
	}
	for _, info := range cols {
		if inf, ok := info.(internal.Closer); ok {
			inf.Close()
		}
	}
}

// stages returns the ordered pipeline of producer/transform stages. The final
// sink (Output) is run separately by Execute since it produces no stream.
func (u *Urlinsane) stages() []Stage {
	return []Stage{
		stageFunc{"init", func(ctx context.Context, _ <-chan *db.Domain) <-chan *db.Domain {
			return u.Init(ctx)
		}},
		stageFunc{"algorithms", func(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
			return u.Algorithms(ctx, in)
		}},
		stageFunc{"constraints", func(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
			return u.Constraints(ctx, in, Dedup, Regex, Levenshtein)
		}},
		stageFunc{"collectors", func(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
			return u.Collectors(ctx, in)
		}},
		stageFunc{"analyzers", func(ctx context.Context, in <-chan *db.Domain) <-chan *db.Domain {
			return u.Analyzers(ctx, in)
		}},
	}
}

// Execute runs the full scan pipeline as a composition of stages, threading a
// single cancellable context through every stage and into the collectors.
func (u *Urlinsane) Execute(ctx context.Context) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var stream <-chan *db.Domain
	for _, st := range u.stages() {
		stream = st.Run(ctx, stream)
	}
	u.Output(ctx, stream)
	u.Close()
	return ctx.Err()
}

func Banner(cfg internal.Config) {
	var lang, board, algo, collectors []string
	t := time.Now()
	timestamp := t.Format("2006-01-02 15:04:05")
	name := fmt.Sprintf("%s  (%s)", text.FgRed.Sprint(cfg.Target()), cfg.EntityType())
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
