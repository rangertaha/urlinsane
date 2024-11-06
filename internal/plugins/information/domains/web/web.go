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
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
)

const (
	ORDER       = 2
	CODE        = "web"
	NAME        = "Download Webpage"
	DESCRIPTION = "Retrieving the web page contents"
)

type None struct {
	types []string
}

func (p *None) Id() string {
	return CODE
}

func (p *None) Order() int {
	return ORDER
}

func (p *None) Description() string {
	return DESCRIPTION
}

func (p *None) Headers() []string {
	return []string{"IMAGE"}
}

func (p *None) Exec(in internal.Typo) (out internal.Typo) {
	if v := in.Variant(); v.Live() {
		file := p.Screenshot(v.Name())
		v.Add("IMAGE", file)
	}

	return in
}

func (p *None) Screenshot(domain string) (filename string) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	var buf []byte
	url := fmt.Sprintf("http://%s", domain)
	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
		url := fmt.Sprintf("https://%s", domain)
		if err := chromedp.Run(ctx, fullScreenshot(url, 90, &buf)); err != nil {
			// log.Fatal(err)
			return ""
		}
	}
	filename = fmt.Sprintf("main/files/%s.png", domain)
	if err := os.WriteFile(filename, buf, 0o644); err != nil {
		// log.Fatal(err)
		return ""
	}
	return filename
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

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &None{}
	})
}
