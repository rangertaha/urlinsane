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
package urlinsane

import (
	"sync"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/pkg/target"
	"github.com/schollz/progressbar/v3"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		Config config.Config
		Typos  map[string]internal.Typo

		// algoWG sync.WaitGroup
		infoWG sync.WaitGroup

		// ..
		progress *progressbar.ProgressBar

		// Metrics
		total     int64
		online    int64
		filtered  int64
		duplicate int64
		processed int64
	}
)

// func init() {
// 	// Log as JSON instead of the default ASCII formatter.
// 	// log.SetFormatter(&log.JSONFormatter{})

// 	// Output to stdout instead of the default stderr
// 	// Can be any io.Writer, see below for File example
// 	logrus.SetOutput(os.Stdout)

// 	// Only log the warning severity or above.
// 	logrus.SetLevel(logrus.DebugLevel)

// 	// contextLogger := log.WithFields(log.Fields{
// 	// 	"common": "this is a common field",
// 	// 	"other": "I also should be logged always",
// 	//   })

// }

// NewUrlinsane ...
func New(conf config.Config) (u Urlinsane) {
	return Urlinsane{
		total:     0,
		online:    0,
		filtered:  0,
		duplicate: 0,
		processed: 0,
		Config:    conf,
		Typos:     make(map[string]internal.Typo),
		progress:  progressbar.DefaultSilent(0),
	}
}

// Init
func (u *Urlinsane) Init() {
	internal.Banner()
}

// GenOptions typo config options
func (t *Urlinsane) Start() <-chan internal.Typo {
	out := make(chan internal.Typo)
	go func() {
		// Initialize information plugins if needed
		for _, info := range t.Config.Information() {
			if inf, ok := info.(internal.Initializer); ok {
				inf.Init(&t.Config)
			}
		}
		for _, algorithm := range t.Config.Algorithms() {

			// Initialize algorithm plugins if needed
			if al, ok := algorithm.(internal.Initializer); ok {
				al.Init(&t.Config)
			}

			out <- &Typo{
				algorithm: algorithm,
				original:  t.Config.Target(),
				variant:   &target.Target{},
			}
		}
		close(out)
	}()
	return out
}

// Algorithms generate typo variations using the algorithm plugins
func (u *Urlinsane) Algorithms(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)
	var wg sync.WaitGroup
	var ttype = u.Config.Type() == internal.DOMAIN

	for w := 1; w <= u.Config.Concurrency(); w++ {
		wg.Add(1)
		go func(id int, in <-chan internal.Typo, out chan<- internal.Typo) {
			defer wg.Done()
			for typo := range in {
				algo := typo.Algorithm()
				typos := []internal.Typo{}

				if al, ok := algo.(internal.DomainAlgorithm); ok && ttype {
					typos = append(typos, al.Domain(typo)...)
				}
				if al, ok := algo.(internal.PackageAlgorithm); ok && ttype {
					typos = append(typos, al.Package(typo)...)
				}
				if al, ok := algo.(internal.EmailAlgorithm); ok && ttype {
					typos = append(typos, al.Email(typo)...)
				}
				if al, ok := algo.(internal.UserAlgorithm); ok && ttype {
					typos = append(typos, al.Username(typo)...)
				}
				if al, ok := algo.(internal.ExecAlgorithm); ok {
					typos = append(typos, al.Exec(typo)...)
				}

				// Execute typo algorith returning typos
				for _, variant := range typos {
					if variant != nil {
						out <- variant
					}
				}
			}
		}(w, in, out)
	}

	go func() {
		wg.Wait()

		// Update total after all algorithms complete procducing typos
		u.total = int64(len(u.Typos))
		if u.Config.Progress() {
			// Add a visible progesssbar if -p flag is set
			u.progress = progressbar.Default(u.total)
		}
		close(out)
	}()

	return out
}

func (u *Urlinsane) Filters(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)
	distance := u.Config.Dist()

	go func(dist int) {
		for typo := range in {
			var filtered = false
			orig := typo.Original()
			vari := typo.Variant()

			if _, ok := u.Typos[vari.Name()]; !ok {
				u.Typos[vari.Name()] = typo

				// Make sure the variant does not match the original
				if vari.Name() != orig.Name() {
					filtered = false
				}

				if distance < typo.Ld() {
					filtered = true
				}

				if !filtered {
					out <- typo
				} else {
					u.filtered++
				}
			} else {
				u.duplicate++
			}
		}
		close(out)
	}(distance)

	return out
}

// func (t *Urlinsane) Analysis(in <-chan internal.Typo) <-chan internal.Typo {
// 	logrus.Debug("Analysis()")
// 	// out := make(chan internal.Typo)

// 	// for typo := range in {
// 	// 	// Dedup typo variants by checking and adding to a map
// 	// 	if variant, ok := t.Typos[typo.Variant().Name()]; !ok {
// 	// 		t.Typos[typo.Variant().Name()] = variant

// 	// 		// Make sure the variant does not match the original
// 	// 		if typo.Variant().Name() != typo.Original().Name() {
// 	// 			t.count++
// 	// 			out <- typo
// 	// 		}
// 	// 	}
// 	// }

// 	// for _, typo := range t.Typos {
// 	// 	out <- typo
// 	// }

// 	return in
// }

// func (t *Urlinsane) Cache(in <-chan internal.Typo) <-chan internal.Typo {
// 	logrus.Debug("Cache()")
// 	// out := make(chan internal.Typo)

// 	// for typo := range in {
// 	// 	// Dedup typo variants by checking and adding to a map
// 	// 	if variant, ok := t.Typos[typo.Variant().Name()]; !ok {
// 	// 		t.Typos[typo.Variant().Name()] = variant

// 	// 		// Make sure the variant does not match the original
// 	// 		if typo.Variant().Name() != typo.Original().Name() {
// 	// 			t.count++
// 	// 			out <- typo
// 	// 		}
// 	// 	}
// 	// }

// 	// for _, typo := range t.Typos {
// 	// 	out <- typo
// 	// }

// 	return in
// }

func (t *Urlinsane) Information(in <-chan internal.Typo) <-chan internal.Typo {
	if len(t.Config.Information()) > 0 {
		out := make(chan internal.Typo)
		var wg sync.WaitGroup

		for w := 1; w <= t.Config.Concurrency(); w++ {
			wg.Add(1)
			go func(in <-chan internal.Typo, out chan<- internal.Typo) {
				defer wg.Done()

				for c := range t.InfoChain(t.Config.Information(), in) {
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
	return in
}

// InfoChain creates a chain of information-gathering functions
func (t *Urlinsane) InfoChain(funcs []internal.Information, in <-chan internal.Typo) <-chan internal.Typo {
	if len(funcs) == 0 {
		return in
	}
	var xfunc internal.Information
	out := make(chan internal.Typo)
	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			if fn, ok := xfunc.(internal.Initializer); ok {
				fn.Init(&t.Config)
			}
			time.Sleep(t.Config.Random() * t.Config.Delay())
			out <- xfunc.Exec(i)
		}
		close(out)
	}()

	if len(funcs) > 0 {
		return t.InfoChain(funcs, out)
	}

	return out
}

func (t *Urlinsane) Progress(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)
	// var wg sync.WaitGroup
	// wg.Add(1)
	go func(in <-chan internal.Typo, out chan<- internal.Typo) {
		// defer wg.Done()
		for c := range in {
			out <- c
			t.progress.Add(1)
		}
		// Clear/hide the progress bar after all typos have passed through
		t.progress.Clear()
		close(out)

	}(in, out)

	// go func() {
	// 	wg.Wait()
	// 	close(out)
	// }()

	return out
}

func (t *Urlinsane) Output(in <-chan internal.Typo) {
	// Initialize output plugin if needed and provide config
	if out, ok := t.Config.Output().(internal.Initializer); ok {
		out.Init(&t.Config)
	}

	// Stream typo records to the output plugin
	for c := range in {

		if c.Variant().Live() {
			t.online++
		}

		if t.Config.All() {
			t.Config.Output().Write(c)
		} else {
			if c.Variant().Live() {
				t.Config.Output().Write(c)
			}
		}
	}

	// Save typo records collected by the output plugin
	t.Config.Output().Save()

	// Print summary
	report := map[string]int64{
		"TOTAL:":   t.total,
		"LIVE:":    t.online,
		"FILTERED": t.filtered,
	}
	t.Config.Output().Summary(report)

}

func (u *Urlinsane) Close() {
	// Initialize information plugins if needed
	for _, info := range u.Config.Information() {
		if inf, ok := info.(internal.Closer); ok {
			inf.Close()
		}
	}
}

func (t *Urlinsane) Execute() {
	t.Init()
	typos := t.Start()
	typos = t.Algorithms(typos)
	typos = t.Filters(typos)
	typos = t.Information(typos)
	typos = t.Progress(typos)
	t.Output(typos)
}
