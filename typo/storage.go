// // Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
// //
// // This program is free software: you can redistribute it and/or modify
// // it under the terms of the GNU General Public License as published by
// // the Free Software Foundation, either version 3 of the License, or
// // (at your option) any later version.
// //
// // This program is distributed in the hope that it will be useful,
// // but WITHOUT ANY WARRANTY; without even the implied warranty of
// // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// // GNU General Public License for more details.
// //
// // You should have received a copy of the GNU General Public License
// // along with this program.  If not, see <http://www.gnu.org/licenses/>.
package typo

// import (
// 	"context"
// 	"fmt"

// 	"github.com/olivere/elastic"
// 	"github.com/spf13/viper"
// )

// // Storager stores and queries result records
// type Storager interface {
// 	Query(Domain) []Domain
// 	Save(Domain)
// }

// // ElasticStorage ...
// type ElasticStorage struct {
// 	client *elastic.Client
// 	index  string
// }

// // NewElasticsearch ...
// func NewElasticsearch() (es ElasticStorage, err error) {
// 	host := viper.GetString("elastic.host")
// 	port := viper.GetString("elastic.port")
// 	username := viper.GetString("elastic.username")
// 	password := viper.GetString("elastic.password")
// 	index := viper.GetString("elastic.index")
// 	client, err := elastic.NewClient(
// 		elastic.SetURL(fmt.Sprintf("http://%s:%s", host, port)),
// 		elastic.SetBasicAuth(username, password))
// 	es.client = client
// 	es.index = index
// 	return es, err
// }

// // Query ...
// func (es *ElasticStorage) Query(r Domain) (res []Domain) {
// 	searchResult, err := es.client.Search().Index(es.index).Query(elastic.NewMatchAllQuery()).Do(context.Background())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(searchResult)
// 	return
// }

// // Save ...
// func (es *ElasticStorage) Save(r Domain) {
// 	ctx := context.Background()
// 	index := viper.GetString("elastic.index")
// 	_, err := es.client.Index().Index(index).Type("domain").Id(r.String()).BodyJson(r).Do(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
