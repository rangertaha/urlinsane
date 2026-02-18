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
package web

import (
	"github.com/gocolly/colly/v2"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

type Plugin struct {
	collectors.Plugin
	client *colly.Collector
}

func (p *Plugin) Init(c internal.Config) {
	p.Plugin.Init(c)
	p.client = colly.NewCollector()
}

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	// res := &Response{
	// 	HTML: HTML{
	// 		Meta: []Metatags{},
	// 	},
	// }
	// if acc.Domain().Live() {

	// 	p.client.OnRequest(func(r *colly.Request) {
	// 		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	// 	})

	// 	p.client.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 		href := e.Attr("href")
	// 		if len(href) > len(acc.Domain().String())+10 {
	// 			res.HTML.Links = append(res.HTML.Links, href)
	// 		}
	// 	})

	// 	p.client.OnHTML("img", func(e *colly.HTMLElement) {
	// 		src := e.Attr("src")
	// 		if len(src) > len(acc.Domain().String())+10 {
	// 			res.HTML.Images = append(res.HTML.Images, src)
	// 		}
	// 	})

	// 	p.client.OnHTML("title", func(e *colly.HTMLElement) {
	// 		res.HTML.Title = e.Text
	// 	})

	// 	p.client.OnHTML("meta", func(e *colly.HTMLElement) {
	// 		res.HTML.Title = e.Text
	// 		meta := Metatags{
	// 			Property: e.Attr("property"),
	// 			Name:     e.Attr("name"),
	// 			Value:    e.Attr("content"),
	// 		}
	// 		res.HTML.Meta = append(res.HTML.Meta, meta)

	// 	})

	// 	p.client.OnError(func(r *colly.Response, err error) {
	// 		acc.SetMeta("STATUS", fmt.Sprint(r.StatusCode))
	// 		res.StatusCode = r.StatusCode
	// 	})

	// 	// attach callbacks after login
	// 	p.client.OnResponse(func(r *colly.Response) {
	// 		acc.SetMeta("STATUS", fmt.Sprint(r.StatusCode))
	// 		res.StatusCode = r.StatusCode
	// 		res.Headers = Header(*r.Headers)
	// 		res.URL = r.Request.URL.String()

	// 		if p.conf.AssetDir() != "" {
	// 			fpath := filepath.Join(p.dir, acc.Domain().String(), "/index.html")
	// 			log.Debugf("Save file: %s", fpath)

	// 			if err := r.Save(fpath); err != nil {
	// 				log.Error(err)
	// 			}
	// 		}

	// 		// SSDeep
	// 		hash, err := ssdeep.FuzzyBytes(r.Body)
	// 		if err != nil {
	// 			log.Error("SSDeep: ", err)
	// 		}
	// 		res.SSDeep = hash

	// 		// Keyword Extraction

	// 	})

	// 	// p.client.Visit(fmt.Sprintf("http://%s/robot.txt", acc.Domain().String()))
	// 	// p.client.Visit(fmt.Sprintf("https://%s/robot.txt", acc.Domain().String()))
	// 	p.client.Visit(fmt.Sprintf("http://%s", acc.Domain().String()))
	// 	p.client.Visit(fmt.Sprintf("https://%s", acc.Domain().String()))

	// 	acc.SetJson("WEB", res.Json())

	// }
	return domain, err
}

// Register the plugin
func init() {
	var CODE = "web"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       8,
				Code:      CODE,
				Title:     "Web Request",
				Summary:   "Gets the website content",
				DependsOn: []string{},
			},
		}
	})
}
