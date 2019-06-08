// The MIT License (MIT)
//
// Copyright © 2018 Tal Hachi
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
	"fmt"
	"os"

	"github.com/cybersectech-org/urlinsane/pkg/typo"
	"github.com/cybersectech-org/urlinsane/pkg/typo/languages"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Property ...
type Property struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Optional    bool            `json:"optional"`
	Values      []PropertyValue `json:"values,omitempty"`
}

// PropertyValue ...
type PropertyValue struct {
	Value       string `json:"value"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Response ...
type Response struct {
	Headers []string                 `json:"headers"`
	Rows    []map[string]interface{} `json:"rows"`
}

// Properties ...
type Properties map[string]Property

var concurrency int
var properties *Properties

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)

	// properties = &Properties{
	// 	"domain": Property{
	// 		Type:        "input",
	// 		Optional:    false,
	// 		Description: "The domain",
	// 	},
	// 	"funcs": Property{
	// 		Type:        "multi-select",
	// 		Optional:    true,
	// 		Description: "Extra functions for data or filtering (default [idna])",
	// 		Values:      getFuncOptions(),
	// 	},
	// 	"typos": Property{
	// 		Type:        "multi-select",
	// 		Optional:    true,
	// 		Description: "The domain",
	// 		Values:      getTypoOptions(),
	// 	},
	// 	"keyboards": Property{
	// 		Type:        "multi-select",
	// 		Optional:    true,
	// 		Description: "Keyboards/layouts ID to use (default [en1])",
	// 		Values:      getKeyboardOptions(),
	// 	},
	// }

}

func errorHandler(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func getTypoOptions() (p []PropertyValue) {
	for _, t := range typo.TRetrieve("all") {
		p = append(p, PropertyValue{t.Code, t.Name, t.Description})
	}
	return
}

func getFuncOptions() (p []PropertyValue) {
	for _, t := range typo.FRetrieve("all") {
		p = append(p, PropertyValue{t.Code, t.Name, t.Description})
	}
	return
}

func getKeyboardOptions() (p []PropertyValue) {
	for _, t := range languages.KEYBOARDS.Keyboards("all") {
		p = append(p, PropertyValue{t.Code, t.Name, t.Description})
	}
	return
}

// NewResponse ...
func NewResponse(results []typo.TypoResult) (resp Response) {
	for _, record := range results {
		m := make(map[string]interface{})

		for key, value := range record.Data {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)
			m[strKey] = strValue
		}

		m["Live"] = record.Live
		m["Variant"] = record.Variant.String()
		m["Typo"] = record.Typo.Name
		resp.Rows = append(resp.Rows, m)
	}
	if len(resp.Rows) > 0 {
		for k := range resp.Rows[0] {
			resp.Headers = append(resp.Headers, k)
		}
	}

	return resp
}

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
				urli := typo.New(config.Config())

				// Stream response
				results := urli.Stream()
				for r := range results {
					// Write
					data, _ := json.Marshal(r)
					fmt.Println(string(data))
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
