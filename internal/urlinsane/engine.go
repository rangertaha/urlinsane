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
		Typos  map[string]Typo

		algoWG sync.WaitGroup
		infoWG sync.WaitGroup

		progress *progressbar.ProgressBar
		live     int64
		count    int64
	}
)

// NewUrlinsane ...
func New(conf config.Config) Urlinsane {
	return Urlinsane{
		live:   0,
		count:  0,
		Config: conf,
	}
}

// Init typo config options
func (t *Urlinsane) Init() {
	// Used for deduping and updating the count
	t.Typos = make(map[string]Typo)

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

// Algorithms generate typos using the algorithm plugins
func (ts *Urlinsane) Algorithms(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)

	for w := 1; w <= ts.Config.Concurrency(); w++ {
		ts.algoWG.Add(1)
		go func(id int, in <-chan internal.Typo, out chan<- internal.Typo) {
			defer ts.algoWG.Done()
			for typo := range in {
				algo := typo.Algorithm()

				//
				if al, ok := algo.(internal.Initializer); ok {
					al.Init(&ts.Config)
				}
				// Execute typo algorithm returning typos
				for _, newtypo := range algo.Exec(typo) {

					// Dedup typo variants by checking and adding to a map
					if variant, ok := ts.Typos[newtypo.Variant().Name()]; !ok {
						ts.Typos[newtypo.Variant().Name()] = variant

						// Make sure the variant does not match the original
						if newtypo.Variant().Name() != newtypo.Original().Name() {
							out <- newtypo
						}
					}
				}
			}
		}(w, in, out)
	}
	go func() {
		ts.algoWG.Wait()
		close(out)
	}()

	return out
}

func (t *Urlinsane) Cache(in <-chan internal.Typo) <-chan internal.Typo {
	// out := make(chan internal.Typo)

	// for typo := range in {
	// 	// Dedup typo variants by checking and adding to a map
	// 	if variant, ok := t.Typos[typo.Variant().Name()]; !ok {
	// 		t.Typos[typo.Variant().Name()] = variant

	// 		// Make sure the variant does not match the original
	// 		if typo.Variant().Name() != typo.Original().Name() {
	// 			t.count++
	// 			out <- typo
	// 		}
	// 	}
	// }

	// for _, typo := range t.Typos {
	// 	out <- typo
	// }

	return in
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
		// close(out)
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
				if t.count != 0 && t.progress == nil {
					t.progress = progressbar.Default(t.count)
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
	t.Config.Output().Summary(t.count, t.live)

	// Save typo records collected by the output plugin
	t.Config.Output().Save()
}

func (t *Urlinsane) Start() {
	t.Init()
	typos := t.GenOptions()
	typos = t.Algorithms(typos)
	typos = t.Cache(typos)
	typos = t.Information(typos)
	typos = t.Storage(typos)
	typos = t.Progress(typos)
	t.Output(typos)
}
