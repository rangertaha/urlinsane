package main

import (
	"fmt"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func main() {
	// Extract the domain from a string
	// using the default list
	fmt.Println(publicsuffix.Domain("example.com"))       // example.com
	fmt.Println(publicsuffix.Domain("www.example.com"))   // example.com
	fmt.Println(publicsuffix.Domain("example.co.uk"))     // example.co.uk
	fmt.Println(publicsuffix.Domain("www.example.co.uk")) // example.co.uk

	// Parse the domain from a string
	// using the default list
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println(publicsuffix.Parse("example.com"))       // &DomainName{"com", "example", ""}
	fmt.Println(publicsuffix.Parse("www.example.com"))   // &DomainName{"com", "example", "www"}
	fmt.Println(publicsuffix.Parse("example.co.uk"))     // &DomainName{"co.uk", "example", ""}
	fmt.Println(publicsuffix.Parse("www.example.co.uk")) // &DomainName{"co.uk", "example", "www"}
	fmt.Println(publicsuffix.Parse("www.example.co.uk")) // &DomainName{"co.uk", "example", "www"}
	fmt.Println(publicsuffix.Parse("example"))           // &DomainName{"co.uk", "example", "www"}
	d, _ := publicsuffix.Parse("example.com.uk.io")
	fmt.Println(d.String(), " SLD:", d.SLD,  " TLD:",d.TLD,  " TRD:",d.TRD)
}

// package main

// import (
//     "fmt"
//     "net"
//     "net/url"
// )

// func main() {

// // We’ll parse this example URL, which includes a scheme, authentication info, host, port, path, query params, and query fragment.

//     // s := "postgres://user:pass@host.com:5432/path?k=v#f"
// 	// s := "postgres://username@host.com:5432/path?k=v#f"
// 	s := "https://username@host.com/one/two/three/four/username"

// // Parse the URL and ensure there are no errors.

//     u, err := url.Parse(s)
//     if err != nil {
//         panic(err)
//     }

// // Accessing the scheme is straightforward.

//     fmt.Println(u.Scheme)

// // User contains all authentication info; call Username and Password on this for individual values.

//     fmt.Println(u.User)
//     fmt.Println(u.User.Username())
//     p, _ := u.User.Password()
//     fmt.Println(p)

// // The Host contains both the hostname and the port, if present. Use SplitHostPort to extract them.

//     fmt.Println("host", u.Host)
//     host, port, _ := net.SplitHostPort(u.Host)
//     fmt.Println(host)
//     fmt.Println(port)

// // Here we extract the path and the fragment after the #.

//     fmt.Println(u.Path)
//     fmt.Println(u.Fragment)

// // To get query params in a string of k=v format, use RawQuery. You can also parse query params into a map. The parsed query param maps are from strings to slices of strings, so index into [0] if you only want the first value.

//     // fmt.Println(u.RawQuery)
//     // m, _ := url.ParseQuery(u.RawQuery)
//     // fmt.Println(m)
//     // fmt.Println(m["k"][0])
// }
