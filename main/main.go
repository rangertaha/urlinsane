package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.npmjs.com/package/api-server")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	// fmt.Println(string(body))
}

// import (
// 	"fmt"
// 	"net"
// 	"time"
// )

// func main() {
// 	host := "facebook.com"
// 	// port := os.Args[2]

// 	// Connect to the target host and port
// 	conn, err := net.DialTimeout("tcp", host+":80", 5*time.Second)
// 	if err != nil {
// 		fmt.Println("Error:", err.Error())
// 		return
// 	}
// 	defer conn.Close()

// 	// Send the request to the server
// 	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\n\r\n", host)

// 	// Read the response from the server
// 	buffer := make([]byte, 1024)
// 	n, _ := conn.Read(buffer)

// 	// Print the response
// 	response := string(buffer[:n])
// 	fmt.Println(response)
// }

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/chromedp/cdproto/cdp"
// 	"github.com/chromedp/chromedp"
// )

// type Product struct {
// 	Url, Image, Name, Price string
// }

// func main() {
// 	var products []Product

// 	// initialize a chrome instance
// 	ctx, cancel := chromedp.NewContext(
// 		context.Background(),
// 		chromedp.WithLogf(log.Printf),
// 	)
// 	defer cancel()

// 	// navigate to the target web page and select the HTML elements of interest
// 	var nodes []*cdp.Node
// 	chromedp.Run(ctx,
// 		chromedp.Navigate("https://www.scrapingcourse.com/ecommerce"),
// 		chromedp.Nodes(".product", &nodes, chromedp.ByQueryAll),
// 	)

// 	// scraping data from each node
// 	var url, image, name, price string
// 	for _, node := range nodes {
// 		chromedp.Run(ctx,
// 			chromedp.AttributeValue("a", "href", &url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
// 			chromedp.AttributeValue("img", "src", &image, nil, chromedp.ByQuery, chromedp.FromNode(node)),
// 			chromedp.Text(".product-name", &name, chromedp.ByQuery, chromedp.FromNode(node)),
// 			chromedp.Text(".price", &price, chromedp.ByQuery, chromedp.FromNode(node)),
// 		)

// 		product := Product{}

// 		product.Url = url
// 		product.Image = image
// 		product.Name = name
// 		product.Price = price

// 		products = append(products, product)
// 	}

// 	fmt.Print(products)

// 	//... export logic
// }

// import "github.com/gocolly/colly"

// // initialize a data structure to keep the scraped data
// type Product struct {
// 	Url, Image, Name, Price string
// }

// func main() {
// 	var products []Product

// 	// instantiate a new collector object
// 	c := colly.NewCollector(
// 		colly.AllowedDomains("www.scrapingcourse.com"),
// 	)

// 	// OnHTML callback
// 	c.OnHTML("li.product", func(e *colly.HTMLElement) {

// 		// initialize a new Product instance
// 		product := Product{}

// 		// scrape the target data
// 		product.Url = e.ChildAttr("a", "href")
// 		product.Image = e.ChildAttr("img", "src")
// 		product.Name = e.ChildText(".product-name")
// 		product.Price = e.ChildText(".price")

// 		// add the product instance with scraped data to the list of products
// 		products = append(products, product)

// 		// OnHTML callback
// 		c.OnHTML("li.product", func(e *colly.HTMLElement) {

// 			// initialize a new Product instance
// 			product := Product{}

// 			// scrape the target data
// 			product.Url = e.ChildAttr("a", "href")
// 			product.Image = e.ChildAttr("img", "src")
// 			product.Name = e.ChildText(".product-name")
// 			product.Price = e.ChildText(".price")

// 			// add the product instance with scraped data to the list of products
// 			products = append(products, product)

// 		})

// 	})

// }

// import (
// 	"fmt"
// 	"log"
// 	"net/netip"

// 	"github.com/oschwald/maxminddb-golang/v2"
// )
// func main() {
// 	db, err := maxminddb.Open("../GeoLite2-City.mmdb")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	addr := netip.MustParseAddr("81.2.69.142")

// 	var record any
// 	err = db.Lookup(addr).Decode(&record)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	fmt.Printf("%v", record)
// }
