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
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/domain"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/typo"
	"github.com/schollz/progressbar/v3"
)

type (

	// Urlinsane ...
	Urlinsane struct {
		Config config.Config
		Typos  map[string]internal.Typo

		// Processing
		// infoWG   sync.WaitGroup
		progress *progressbar.ProgressBar

		// Metrics
		total    int64
		online   int64
		filtered int64
		scanned  int64
	}
)

// NewUrlinsane ...
func New(conf config.Config) (u Urlinsane) {
	return Urlinsane{
		total:    0,
		online:   0,
		filtered: 0,
		scanned:  0,
		Config:   conf,
		Typos:    make(map[string]internal.Typo),
		progress: progressbar.DefaultSilent(0),
	}
}

// Init
func (u *Urlinsane) Init() <-chan internal.Typo {
	out := make(chan internal.Typo)
	go func() {
		// Initialize information plugins if needed
		for _, info := range u.Config.Information() {
			if inf, ok := info.(internal.Initializer); ok {
				inf.Init(&u.Config)
			}
		}

		// Initialize algorithm plugins if needed
		for _, algorithm := range u.Config.Algorithms() {
			if al, ok := algorithm.(internal.Initializer); ok {
				al.Init(&u.Config)
			}
			if u.Config.Target() != "" {
				// fmt.Println(domain.Parse(u.Config.Target()))
				out <- &typo.Typo{
					Algorithm: algorithm,
					Original:  domain.Parse(u.Config.Target()),
					Variant:   models.Domain{},
				}
			} else {
				fmt.Println("Need to proved a domain name")
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
	// var ttype = u.Config.Type() == internal.DOMAIN

	for w := 1; w <= u.Config.Concurrency(); w++ {
		wg.Add(1)
		go func(id int, in <-chan internal.Typo, out chan<- internal.Typo) {
			defer wg.Done()
			for typo := range in {
				algo := typo.Algo()
				for _, typ := range algo.Exec(typo) {
					if typ != nil {
						out <- typ
					}
				}

				// Execute typo algorith returning typos
				// for _, variant := range typos {
				// 	if variant != nil {
				// 		out <- variant
				// 	}
				// }
			}
		}(w, in, out)
	}

	go func() {
		wg.Wait()

		// Update total after all algorithms complete procducing typos
		u.total = int64(len(u.Typos))
		// if u.Config.Progress() {
		// 	u.progress = progressbar.Default(u.total)
		// }
		close(out)
	}()

	return out
}

func (u *Urlinsane) Filters(in <-chan internal.Typo) <-chan internal.Typo {
	out := make(chan internal.Typo)

	go func() {
		for typo := range in {
			orig, vari := typo.Get()
			// fmt.Println(orig, vari.FQDN)

			// Removing duplicates
			if _, ok := u.Typos[vari.Name]; !ok {
				u.Typos[vari.Name] = typo

				// Make sure the variant does not match the original
				if vari.Fqdn() != orig.Fqdn() {

					// Only allow variants with a minimum levenshtein distance
					if u.Config.Dist() >= typo.Dist() {
						// fmt.Println(orig, vari.FQDN)
						out <- typo
					} else {
						u.filtered++
					}
				}
			}
		}
		if u.Config.Progress() {
			u.progress = progressbar.Default(u.total - u.filtered)
		}
		close(out)
	}()

	return out
}

func (u *Urlinsane) Information(in <-chan internal.Typo) <-chan internal.Typo {
	if len(u.Config.Information()) > 0 {
		out := make(chan internal.Typo)
		var wg sync.WaitGroup

		for w := 1; w <= u.Config.Concurrency(); w++ {
			wg.Add(1)
			go func(in <-chan internal.Typo, out chan<- internal.Typo) {
				defer wg.Done()
				for c := range u.InfoChain(u.Config.Information(), in) {
					u.scanned++
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
func (u *Urlinsane) InfoChain(funcs []internal.Information, in <-chan internal.Typo) <-chan internal.Typo {
	if len(funcs) == 0 {
		return in
	}
	var xfunc internal.Information
	out := make(chan internal.Typo)
	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			if fn, ok := xfunc.(internal.Initializer); ok {
				fn.Init(&u.Config)
			}
			time.Sleep(u.Config.Random() * u.Config.Delay())

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			func(ctx context.Context) {
				select {
				case <-time.After(1 * time.Second):

					out <- xfunc.Exec(i)
					// fmt.Println("Function completed successfully")
				case <-ctx.Done():
					// fmt.Println("Function timed out:", ctx.Err())
				}
			}(ctx)

		}
		close(out)
	}()

	if len(funcs) == 1 {
		u.scanned++
	}

	if len(funcs) > 0 {
		return u.InfoChain(funcs, out)
	}

	return out
}

// Progress adds a visible progesssbar if -p flag is set
func (u *Urlinsane) Progress(in <-chan internal.Typo) <-chan internal.Typo {
	if u.Config.Progress() {
		out := make(chan internal.Typo)
		go func(in <-chan internal.Typo, out chan<- internal.Typo) {
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

	return in
}

func (u *Urlinsane) Output(in <-chan internal.Typo) {
	// Initialize output plugin if needed and provide config
	if out, ok := u.Config.Output().(internal.Initializer); ok {
		out.Init(&u.Config)
	}

	// Stream typo records to the output plugin
	for c := range in {
		_, vari := c.Get()
		if vari.Live {
			u.online++
		}
		// u.Config.Output().Write(c)
		if u.Config.ShowAll() {
			u.Config.Output().Write(c)

		} else if vari.Live {
			u.Config.Output().Write(c)
		}

	}

	// Save typo records collected by the output plugin
	u.Config.Output().Save()

	// Print summary
	report := map[string]int64{
		"TOTAL:":    u.total,
		"LIVE:":     u.online,
		"FILTERED:": u.filtered,
		"SCANNED:":  u.scanned,
	}
	u.Config.Output().Summary(report)

}

func (u *Urlinsane) Close() {
	// Initialize information plugins if needed
	for _, info := range u.Config.Information() {
		if inf, ok := info.(internal.Closer); ok {
			inf.Close()
		}
	}
}

func (u *Urlinsane) Execute() {
	typos := u.Init()
	typos = u.Algorithms(typos)
	typos = u.Filters(typos)
	typos = u.Information(typos)
	typos = u.Progress(typos)
	u.Output(typos)
}

func (u *Urlinsane) Stream() <-chan internal.Typo {
	typos := u.Init()
	typos = u.Algorithms(typos)
	typos = u.Filters(typos)
	typos = u.Information(typos)
	return typos
}
