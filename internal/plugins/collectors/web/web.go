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
package web

import (
	"fmt"
	"path/filepath"

	"github.com/gocolly/colly/v2"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	"github.com/rangertaha/urlinsane/pkg/fuzzy/ssdeep"
	log "github.com/sirupsen/logrus"
)

const (
	ORDER       = 8
	CODE        = "web"
	DESCRIPTION = "Web request and hasing hashing content"
)

type Plugin struct {
	conf   internal.Config
	db     internal.Database
	client *colly.Collector
	dir    string
}

func (p *Plugin) Id() string {
	return CODE
}

func (p *Plugin) Order() int {
	return ORDER
}

func (i *Plugin) Init(c internal.Config) {
	i.db = c.Database()
	i.conf = c
	i.client = colly.NewCollector()
	i.dir = filepath.Join(c.Dir(), "domains")

}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (p *Plugin) Headers() []string {
	return []string{"STATUS"}
}

func (p *Plugin) Exec(domain internal.Domain, acc internal.Accumulator) (err error) {
	if domain.Live() {
		assetDir, err := acc.Mkdir(p.dir, domain.String())
		if err != nil {
			log.Error("AssetDir: ", err)
		}
		p.client.OnRequest(func(r *colly.Request) {
			r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		})

		// p.client.OnResponse(func(r *colly.Response) {
		// 	if strings.Contains(r.Headers.Get("Content-Type"), "video/mp4") {
		// 		log.Println("dowloading")
		// 		path := "./download/" + uuid.New().String() + ".mp4"
		// 		err := r.Save(path)
		// 		if err != nil {
		// 			log.Println("dowload video error")
		// 			log.Println(err)
		// 			return
		// 		}
		// 		log.Println("dowloaded")
		// 	}
		// })

		// // On every a element which has href attribute call callback
		// p.client.OnHTML("img", func(e *colly.HTMLElement) {
		// 	// Get the URL of the image
		// 	imgSrc := e.Attr("src")

		// 	// Use absolute URL for the image
		// 	imgSrc = e.Request.AbsoluteURL(imgSrc)

		// 	// Download the image
		// 	fmt.Printf("Image found: %s\n", imgSrc)
		// 	// Create a new folder to store images if it does not exist
		// 	// os.MkdirAll("images", os.ModePerm)
		// 	// Download the image and save to the file
		// 	// fileName := "images/" + e.Attr("alt") + ".jpg"
		// 	err := p.client.Visit(imgSrc)
		// 	if err != nil {
		// 		log.Printf("Failed to download image %s: %s", imgSrc, err)
		// 		return
		// 	}
		// 	filename := filepath.Join(assetDir, "image.gif")
		// 	p.client.sa Save(filename)

		// })

		p.client.OnError(func(r *colly.Response, err error) {
			domain.SetMeta("STATUS", fmt.Sprint(r.StatusCode))
		})

		// attach callbacks after login
		p.client.OnResponse(func(r *colly.Response) {
			acc.Mkfile(assetDir, "index.html", r.Body)
			domain.SetMeta("STATUS", fmt.Sprint(r.StatusCode))
			domain.Live(true)

			hash, err := ssdeep.FuzzyBytes(r.Body)
			if err != nil {
				// fmt.Println(err)
				log.Error("SSDeep: ", err)
			}
			domain.SetMeta("HASH", hash)

		})

		p.client.Visit(fmt.Sprintf("http://%s", domain.String()))
		p.client.Visit(fmt.Sprintf("https://%s", domain.String()))

	}
	acc.Add(domain)
	return
}

// Register the plugin
func init() {
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{}
	})
}
