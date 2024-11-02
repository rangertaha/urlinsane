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
		Typos  map[string]*internal.Typo

		algoWG sync.WaitGroup
		infoWG sync.WaitGroup

		progress *progressbar.ProgressBar
		// stats <-chan Statser
		// errs <-chan interface{}
	}
)

// NewUrlinsane ...
func New(conf config.Config) Urlinsane {
	return Urlinsane{
		Config: conf,
	}
}

// Init typo config options
func (t *Urlinsane) Init() {
	// Used for deduping and updating the count
	t.Typos = make(map[string]*internal.Typo)

	internal.Banner()
}

// GenOptions typo config options
func (t *Urlinsane) GenOptions() <-chan internal.Typo {
	out := make(chan internal.Typo)
	go func() {
		for _, algorithm := range t.Config.Algorithms() {

			if al, ok := algorithm.(internal.Initializer); ok {
				al.Init(&t.Config)
			}

			// // fmt.Println("GenOptions: ", algorithm)
			// domain := domain.New(t.Config.Target())

			out <- &Typo{
				// languages: t.Config.Languages(), // remove
				// keyboards: t.Config.Keyboards(), // remove
				algorithm: algorithm,
				original:  t.Config.Target(), // remove and add to config: config.Domain()
				variant:   &target.Target{},
			}
			// fmt.Println("GenOptions: ", algorithm)
		}
		close(out)
	}()
	return out
}

// Algorithms generate typos using the algorithm plugins
func (ts *Urlinsane) Algorithms(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)

	for w := 1; w <= ts.Config.Concurrency(); w++ {
		ts.algoWG.Add(1)
		go func(id int, in <-chan internal.Typo, out chan<- internal.Typo) {
			defer ts.algoWG.Done()
			for typo := range in {
				algo := typo.Algorithm()

				// Execute typo algorithm returning typos
				if al, ok := algo.(internal.Initializer); ok {
					al.Init(&ts.Config)
				}

				// if ts.Config.IsMode(internal.DOMAIN) {
				// 	if fn, ok := algo.(internal.DomainAlgo); ok {
				// 		exec = fn.Domain
				// 	}

				// }
				// if ts.Config.IsMode(internal.USERNAME) {

				// }
				// if ts.Config.IsMode(internal.NAME) {

				// }

				for _, typ := range algo.Exec(typo) {

					// Dedup typo variants by checking and adding to a map
					if variant, ok := ts.Typos[typ.Variant().Name()]; !ok {
						ts.Typos[typ.Variant().Name()] = variant

						// Make sure the variant does not match the original
						if typ.Variant().Name() != typ.Original().Name() {

							// Cache and or reuse
							out <- ts.Cache(typ)
						}
					}
				}
			}
		}(w, in, out)
	}
	go func() {
		ts.algoWG.Wait()
		ts.Config.Count(int64(len(ts.Typos)))
		close(out)
	}()

	return out
}

func (t *Urlinsane) Cache(typo internal.Typo) internal.Typo {

	return typo
}

func (t *Urlinsane) Information(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)
	for w := 1; w <= t.Config.Concurrency(); w++ {
		t.infoWG.Add(1)
		go func(in <-chan internal.Typo, out chan<- internal.Typo) {
			defer t.infoWG.Done()
			output := t.InfoChain(t.Config.Information(), in)
			for c := range output {
				out <- c
			}
		}(in, out)
	}
	go func() {
		t.infoWG.Wait()
		close(out)
	}()
	return out
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

func (t *Urlinsane) Storage(in <-chan internal.Typo) <-chan internal.Typo {

	return in
}

func (t *Urlinsane) Progress(in <-chan internal.Typo) <-chan internal.Typo {
	if t.Config.Progress() {
		out := make(chan internal.Typo)
		go func(in <-chan internal.Typo, out chan<- internal.Typo) {
			for c := range in {
				if t.Config.Count() != 0 && t.progress == nil {
					t.progress = progressbar.Default(t.Config.Count())
				}
				out <- c
				t.progress.Add(1)
			}
			close(out)
		}(in, out)
		return out
	}
	return in
}

func (t *Urlinsane) Output(in <-chan internal.Typo) {
	// Initialize output plugin and provide config
	if out, ok := t.Config.Output().(internal.Initializer); ok {
		out.Init(&t.Config)
	}

	// Stream typo records to the output plugin
	for c := range in {
		t.Config.Output().Write(c)
	}

	// Save typo records collected by the output plugin
	t.Config.Output().Save()
}

func (t *Urlinsane) Start() {
	t.Init()
	typos := t.GenOptions()
	typos = t.Algorithms(typos)
	typos = t.Information(typos)
	typos = t.Storage(typos)
	typos = t.Progress(typos)
	t.Output(typos)
}
