// The MIT License (MIT)
//
// Copyright © 2018 rangertaha rangertaha@gmail.com
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

package typo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olivere/elastic"
)

// Storager stores and queries result records
type Storager interface {
	Query(Domain) []Domain
	Save(Domain)
}

// ElasticStorage ...
type ElasticStorage struct {
}

// Query ...
func (es *ElasticStorage) Query(r Domain) (res []Domain) {
	// Create a client and connect to http://127.0.0.1:9201
	ES, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9201"))
	if err != nil {
		// Handle error
	}

	searchResult, err := ES.Search().Index("urlinsane").Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(searchResult)
	return
}

// Save ...
func (es *ElasticStorage) Save(r Domain) {
	json, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(json))
	//fmt.Println(r)
	// ctx := context.Background()
	// ES, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9201"))
	// if err != nil {
	// 	// Handle error
	// }

	// searchResult, err := ES.Index().Index("urlinsane").Type("domain").Id("1").BodyJson(r).Do(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(searchResult)
}
