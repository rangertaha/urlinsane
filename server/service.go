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
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"

	"github.com/cybersectech-org/urlinsane"
	"github.com/cybersectech-org/urlinsane/languages"
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

// NewTCPServer ...
func NewTCPServer(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")
	l, nerr := net.Listen("tcp", address+":"+port)
	if nerr != nil {
		fmt.Println("ERROR", nerr)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		// Read
		if err != nil {
			fmt.Println("ERROR", err)
			continue
		}

		go func(conn net.Conn) {
			r := bufio.NewReader(conn)
			for {
				input, err := r.ReadBytes(byte('\n'))
				switch err {
				case nil:
					break
				case io.EOF:
				default:
					fmt.Println("ERROR", err)
				}

				config := new(urlinsane.BasicConfig)
				config.Concurrency = concurrency
				if err := json.Unmarshal(input, &config); err != nil {
					fmt.Println("ERROR", err)
				}
				fmt.Println(string(input), config)

				urli := urlinsane.New(config.Config())

				// Stream response
				results := urli.Stream()
				for r := range results {
					// Write
					fmt.Println(r)
					data, err := json.Marshal(r)
					if err != nil {
						fmt.Println("ERROR", err)
					}
					fmt.Println(string(data))
					conn.Write(data)
				}
			}

		}(conn)
	}

}

// NewServer ...
func NewWebServer(cmd *cobra.Command, args []string) {
	// Echo instance
	e := echo.New()
	e.HideBanner = true

	address, err := cmd.Flags().GetString("host")
	port, err := cmd.Flags().GetString("port")
	proto, err := cmd.Flags().GetString("type")
	concurrency, err = cmd.Flags().GetInt("concurrency")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handlers
	fmt.Println(proto)
	fmt.Println("Start handlers")
	if proto == "http" {
		e.POST("/", httpHandler)

	} else if proto == "ws" {
		e.GET("/", websocketHandler)
	}
	e.GET("/options", func(c echo.Context) error {
		return c.JSON(http.StatusOK, properties)
	})

	// Start server
	e.Logger.Fatal(e.Start(address + ":" + port))
}

// postHandler ....
func httpHandler(c echo.Context) (err error) {
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

func websocketHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Read
			config := new(urlinsane.BasicConfig)
			config.Concurrency = concurrency
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			if err := json.Unmarshal([]byte(msg), &config); err != nil {
				c.Logger().Error(err)
			}

			// Initialize urlinsane object
			urli := urlinsane.New(config.Config())

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
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
