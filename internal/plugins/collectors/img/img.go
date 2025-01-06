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
package img

import (
	"context"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

// const (
// 	ORDER       = 8
// 	CODE        = "img"
// 	DESCRIPTION = "Download screeshot of domains"
// )

type Plugin struct {
	collectors.Plugin
	// conf   internal.Config
	dir    string
	ctx    context.Context
	cancel context.CancelFunc
}

// func (p *Plugin) Id() string {
// 	return CODE
// }

// func (p *Plugin) Order() int {
// 	return ORDER
// }

func (i *Plugin) Init(c internal.Config) {
	// i.db = c.Database()
	i.Conf = c
	i.dir = filepath.Join(c.AssetDir(), "domains")
	i.ctx, i.cancel = chromedp.NewContext(context.Background())
	// defer cancel()

	i.ctx, i.cancel = context.WithTimeout(i.ctx, 2*time.Minute)
	// defer cancel()
}

// func (n *Plugin) Description() string {
// 	return DESCRIPTION
// }

// func (p *Plugin) Headers() []string {
// 	return []string{"SCREENSHOT"}
// }

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	// if acc.Live() {

	// 	var buf []byte
	// 	url := fmt.Sprintf("http://%s", acc.Domain().String())
	// 	// capture entire browser viewport, returning png with quality=90
	// 	if err := chromedp.Run(p.ctx, fullScreenshot(url, 90, &buf)); err != nil {
	// 		url := fmt.Sprintf("https://%s", acc.Domain().String())
	// 		if err := chromedp.Run(p.ctx, fullScreenshot(url, 90, &buf)); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}

	// 	if err := acc.Save("index.png", buf); err != nil {
	// 		log.Error(err)
	// 	} else {
	// 		acc.SetMeta("SCREENSHOT", "index.png")
	// 	}
	// }
	return domain, err
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

func (p *Plugin) Close() {
	p.cancel()
}

// // Register the plugin
// func init() {
// 	var CODE = "img"
// 	collectors.Add(CODE, func() internal.Collector {
// 		return &Plugin{}
// 	})
// }

// Register the plugin
func init() {
	var CODE = "img"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       8,
				Code:      CODE,
				Title:     "Screenshot",
				Summary:   "Take screeshot of domains",
				DependsOn: []string{},
			},
		}
	})
}
