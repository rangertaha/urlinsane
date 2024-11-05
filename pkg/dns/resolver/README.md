# Resolver MT [![GoDoc][doc-img]][doc] [![Coverage Status][cov-img]][cov] ![test][test-img]

This package provides a multithreaded and thread-safe DNS resolver client for Go. It is used to quickly resolve DNS queries in applications that need to perform requests concurrently, while also respecting a limit of queries per second for each resolver.

Create a new resolver client and set the resolvers to use, the limit of queries per second and the total number of goroutines to create. When done, close the connections by calling Client.Close().

```
package resolvermt_test

import (
	"fmt"

	"github.com/d3mondev/resolvermt"
)

func main() {
	resolvers := []string{
		"8.8.8.8",
		"9.9.9.10",
		"1.1.1.1",
	}

	domains := []string{
		"www.google.com",
		"facebook.com",
		"uber.com",
		"apple.com",
		"blogger.com",
		"instagram.com",
		"linkedin.com",
		"en.wikipedia.org",
		"cloudflare.com",
		"wordpress.org",
	}

	client := resolvermt.New(resolvers, 3, 10, 5)
	defer client.Close()

	results := client.Resolve(domains, resolvermt.TypeA)

	for _, record := range results {
		fmt.Printf("%s %s %s\n", record.Question, record.Type, record.Answer)
	}
}
```

Output:

```
www.google.com A 172.217.13.164
uber.com A 104.36.195.150
facebook.com A 31.13.67.35
apple.com A 17.253.144.10
instagram.com A 34.226.174.41
instagram.com A 54.86.160.135
instagram.com A 54.210.252.14
instagram.com A 52.87.69.127
instagram.com A 54.205.161.210
instagram.com A 54.174.36.241
instagram.com A 54.165.128.247
instagram.com A 3.208.131.236
blogger.com A 172.217.165.41
cloudflare.com A 104.16.132.229
cloudflare.com A 104.16.133.229
linkedin.com A 13.107.42.14
en.wikipedia.org CNAME dyna.wikimedia.org
en.wikipedia.org A 208.80.154.224
wordpress.org A 198.143.164.252
```

In its current implementation, the client will create at least one or more sockets per resolver to send and receive DNS queries. As each socket uses an open file descriptor, it is easy to hit the limit of open file descriptors in the OS. If this happens, the client will handle the failed connections but the performance may suffer with a large number of resolvers and threads. If you are trying to run a large number of requests concurrently (100+), increase the limit with `ulimit -n 102400`. A future version may change the way queries are made to prevent this.

[doc-img]: https://pkg.go.dev/badge/github.com/d3mondev/resolvermt
[doc]: https://pkg.go.dev/github.com/d3mondev/resolvermt
[cov-img]: https://codecov.io/gh/d3mondev/resolvermt/branch/master/graph/badge.svg?token=3D1I4J6YA2
[cov]: https://codecov.io/gh/d3mondev/resolvermt
[test-img]: https://github.com/d3mondev/resolvermt/workflows/test/badge.svg
