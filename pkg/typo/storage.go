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
	"fmt"

	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

// Storager stores and queries result records
type Storager interface {
	Query(Domain) []Domain
	Save(Domain)
}

// ElasticStorage ...
type ElasticStorage struct {
	client *elastic.Client
	index  string
}

// NewElasticsearch ...
func NewElasticsearch() (es ElasticStorage, err error) {
	host := viper.GetString("elastic.host")
	port := viper.GetString("elastic.port")
	username := viper.GetString("elastic.username")
	password := viper.GetString("elastic.password")
	index := viper.GetString("elastic.index")
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%s", host, port)),
		elastic.SetBasicAuth(username, password))
	es.client = client
	es.index = index
	return es, err
}

// Query ...
func (es *ElasticStorage) Query(r Domain) (res []Domain) {
	searchResult, err := es.client.Search().Index(es.index).Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(searchResult)
	return
}

// Save ...
func (es *ElasticStorage) Save(r Domain) {
	ctx := context.Background()
	index := viper.GetString("elastic.index")
	_, err := es.client.Index().Index(index).Type("domain").Id(r.String()).BodyJson(r).Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
