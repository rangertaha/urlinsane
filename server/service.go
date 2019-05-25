// Copyright © 2018 CyberSecTech Inc
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
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/cybersectech-org/urlinsane"
	"github.com/cybersectech-org/urlinsane/languages"
	"github.com/spf13/cobra"
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
	properties = &Properties{
		"domain": Property{
			Type:        "input",
			Optional:    false,
			Description: "The domain",
		},
		"funcs": Property{
			Type:        "multi-select",
			Optional:    true,
			Description: "Extra functions for data or filtering (default [idna])",
			Values:      getFuncOptions(),
		},
		"typos": Property{
			Type:        "multi-select",
			Optional:    true,
			Description: "The domain",
			Values:      getTypoOptions(),
		},
		"keyboards": Property{
			Type:        "multi-select",
			Optional:    true,
			Description: "Keyboards/layouts ID to use (default [en1])",
			Values:      getKeyboardOptions(),
		},
	}

}

func getTypoOptions() (p []PropertyValue) {
	for _, t := range urlinsane.TRetrieve("all") {
		p = append(p, PropertyValue{t.Code, t.Name, t.Description})
	}
	return
}

func getFuncOptions() (p []PropertyValue) {
	for _, t := range urlinsane.FRetrieve("all") {
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
func NewResponse(results []urlinsane.TypoResult) (resp Response) {
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

// NewServer ...
func NewServer(cmd *cobra.Command, args []string) {
	// Echo instance
	e := echo.New()
	e.HideBanner = true

	address, err := cmd.Flags().GetString("addr.host")
	port, err := cmd.Flags().GetString("addr.port")
	stream, err := cmd.Flags().GetBool("stream")
	concurrency, err = cmd.Flags().GetInt("concurrency")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handlers
	if stream {
		// https://echo.labstack.com/cookbook/streaming-response
		e.POST("/", postStreamHandler)
	} else {
		e.POST("/", postHandler)
	}
	e.GET("/options", func(c echo.Context) error {
		return c.JSON(http.StatusOK, properties)
	})

	// Start server
	e.Logger.Fatal(e.Start(address + ":" + port))
}

// postHandler ....
func postHandler(c echo.Context) (err error) {
	// // Get parameters from json payload
	config := new(urlinsane.BasicConfig)
	config.Concurrency = concurrency
	if err = c.Bind(config); err != nil {
		c.Logger().Error(err)
		return
	}

	// Initialize urlinsane object
	urli := urlinsane.New(config.Config())

	// Execute returning results
	reponse := NewResponse(urli.Execute())

	// Return JSON results
	return c.JSON(http.StatusOK, reponse)
}

// postStreamHandler ...
func postStreamHandler(c echo.Context) (err error) {

	// Get parameters from the context
	config := new(urlinsane.BasicConfig)
	config.Concurrency = concurrency
	if err = c.Bind(config); err != nil {
		c.Logger().Error(err)
		return
	}

	// Initialize urlinsane object
	urli := urlinsane.New(config.Config())

	// Stream response
	results := urli.Stream()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	for r := range results {
		if err := json.NewEncoder(c.Response()).Encode(r); err != nil {
			return err
		}
		c.Response().Flush()
	}
	return nil
}
