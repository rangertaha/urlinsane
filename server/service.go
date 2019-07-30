// The MIT License (MIT)
//
// Copyright © 2018 Rangertaha <rangertaha@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package server

import (
	"encoding/json"

	"github.com/cybint/urlinsane/pkg/typo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
	// "github.com/davecgh/go-spew/spew"
)

// NewWebSocketServer ...
func NewWebSocketServer(host, port string, concurrency int) {
	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				// Read
				config := new(typo.BasicConfig)
				config.Concurrency = concurrency
				msg := ""
				err := websocket.Message.Receive(ws, &msg)
				if err != nil {
					c.Logger().Error(err)
				}
				if err := json.Unmarshal([]byte(msg), &config); err != nil {
					c.Logger().Error(err)
				}

				// Initialize typo object
				typosquatting := typo.New(config.Config())

				// Stream response
				results := typosquatting.Stream()
				for r := range results {
					// Write
					//  spew.Dump(r)
					data, err := json.Marshal(r)
					if err != nil {
						c.Logger().Error(err)
					}

					err = websocket.Message.Send(ws, string(data))
					if err != nil {
						c.Logger().Error(err)
					}
				}
				msgDone := `{"total": {"progress": 100}, "status": "done"}`
				err = websocket.Message.Send(ws, msgDone)
				if err != nil {
					c.Logger().Error(err)
				}
				break
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Start server
	e.Logger.Fatal(e.Start(host + ":" + port))
}
