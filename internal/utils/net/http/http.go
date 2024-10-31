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

package http

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	httpLib "net/http"
)

type (
	Userinfo struct {
		Username    string `json:"username,omitempty"`
		Password    string `json:"password,omitempty"`
		PasswordSet bool   `json:"password_set,omitempty"`
	}
	// ConnectionState records basic TLS details about the connection.
	ConnectionState struct {
		Version                     uint16                `json:"version,omitempty"`                      // TLS version used by the connection (e.g. VersionTLS12); added in Go 1.3
		Complete                    bool                  `json:"complete,omitempty"`                     // TLS handshake is complete
		DidResume                   bool                  `json:"did_resume,omitempty"`                   // connection resumes a previous TLS connection; added in Go 1.1
		CipherSuite                 uint16                `json:"cipher_suite,omitempty"`                 // cipher suite in use (TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, ...)
		NegotiatedProtocol          string                `json:"negotiated_protocol,omitempty"`          // negotiated next protocol (not guaranteed to be from Config.NextProtos)
		NegotiatedProtocolIsMutual  bool                  `json:"negotiated_protocolIsMutual,omitempty"`  // negotiated protocol was advertised by server (client side only)
		ServerName                  string                `json:"server_name,omitempty"`                  // server name requested by client, if any (server side only)
		PeerCertificates            []*x509.Certificate   `json:"-"`                                      // certificate chain presented by remote peer
		VerifiedChains              [][]*x509.Certificate `json:"-"`                                      // verified chains built from PeerCertificates
		SignedCertificateTimestamps [][]byte              `json:"signed_certificateTimestamps,omitempty"` // SCTs from the peer, if any; added in Go 1.5
		OCSPResponse                []byte                `json:"ocsp_response,omitempty"`                // stapled OCSP response from peer, if any; added in Go 1.5
	}
	URL struct {
		Scheme     string    `json:"scheme,omitempty"`
		Opaque     string    `json:"opaque,omitempty"`      // encoded opaque data
		User       *Userinfo `json:"user,omitempty"`        // username and password information
		Host       string    `json:"host,omitempty"`        // host or host:port
		Path       string    `json:"path,omitempty"`        // path (relative paths may omit leading slash)
		RawPath    string    `json:"raw_path,omitempty"`    // encoded path hint (see EscapedPath method)
		ForceQuery bool      `json:"force_query,omitempty"` // append a query ('?') even if RawQuery is empty
		RawQuery   string    `json:"raw_query,omitempty"`   // encoded query values, without '?'
		Fragment   string    `json:"fragment,omitempty"`    // fragment for references, without '#'
	}

	// Response represents the response from an HTTP request.
	Response struct {
		Status           string              `json:"status,omitempty"`      // e.g. "200 OK"
		StatusCode       int                 `json:"status_code,omitempty"` // e.g. 200
		Proto            string              `json:"proto,omitempty"`       // e.g. "HTTP/1.0"
		ProtoMajor       int                 `json:"proto_major,omitempty"`
		ProtoMinor       int                 `json:"proto_minor,omitempty"`
		Header           map[string][]string `json:"header,omitempty"`
		Body             string              `json:"body,omitempty"`
		Length           int64               `json:"length,omitempty"`
		SSDeep           string              `json:"ssdeep,omitempty"`
		TransferEncoding []string            `json:"encoding,omitempty"`
		Uncompressed     bool                `json:"uncompressed,omitempty"`
		Trailer          map[string][]string `json:"trailer,omitempty"`
		Request          Request             `json:"request,omitempty"`
		TLS              ConnectionState     `json:"tls,omitempty"`
	}
	Request struct {
		Method           string               `json:"method,omitempty"` // Method specifies the HTTP method (GET, POST, PUT, etc.).
		URL              *URL                 `json:"url,omitempty"`
		Proto            string               `json:"proto,omitempty"`
		ProtoMajor       int                  `json:"proto_major,omitempty"`
		ProtoMinor       int                  `json:"proto_minor,omitempty"`
		Header           map[string][]string  `json:"header,omitempty"`
		Body             string               `json:"body,omitempty"`
		ContentLength    int64                `json:"length,omitempty"`
		TransferEncoding []string             `json:"encoding,omitempty"`
		Close            bool                 `json:"close,omitempty"`
		Host             string               `json:"host,omitempty"`
		Form             map[string][]string  `json:"form,omitempty"`
		PostForm         map[string][]string  `json:"post_form,omitempty"`
		MultipartForm    *Form                `json:"multipart_form,omitempty"`
		Trailer          map[string][]string  `json:"trailer,omitempty"`
		RemoteAddr       string               `json:"remote_addr,omitempty"`
		RequestURI       string               `json:"request_uri,omitempty"`
		TLS              *tls.ConnectionState `json:"tls,omitempty"`
	}
	Form struct {
		Value map[string][]string      `json:"value,omitempty"`
		File  map[string][]*FileHeader `json:"file,omitempty"`
	}
	FileHeader struct {
		Filename string              `json:"filename,omitempty"`
		Header   map[string][]string `json:"header,omitempty"`
		Size     int64               `json:"size,omitempty"`
	}

	Http struct {
		Request  Request  `json:"request,omitempty"`
		Response Response `json:"response,omitempty"`
	}
)

// NewResponse ...
func NewResponse(resp *httpLib.Response) (nres Response) {
	nres.Status = resp.Status
	nres.StatusCode = resp.StatusCode
	nres.Proto = resp.Proto
	nres.ProtoMajor = resp.ProtoMajor
	nres.ProtoMinor = resp.ProtoMinor
	nres.Header = resp.Header
	if obj, err := ioutil.ReadAll(resp.Body); err == nil {
		nres.Body = string(obj)
	}

	// nres.SSDeep = resp.SSDeep
	nres.TransferEncoding = resp.TransferEncoding
	nres.Uncompressed = resp.Uncompressed
	nres.Trailer = resp.Trailer

	// nres.TLS = resp.TLS
	// nres.TLS.Version = resp.TLS.Version

	// Request
	nres.Request.Proto = resp.Request.Proto
	nres.Request.ProtoMajor = resp.Request.ProtoMajor
	nres.Request.ProtoMinor = resp.Request.ProtoMinor
	nres.Request.Header = resp.Request.Header
	// if obj, err := ioutil.ReadAll(resp.Request.Body); err == nil {
	// 	nres.Request.Body = string(obj)
	// }
	nres.Request.ContentLength = resp.Request.ContentLength
	nres.Request.TransferEncoding = resp.Request.TransferEncoding
	nres.Request.Close = resp.Request.Close
	nres.Request.Host = resp.Request.Host
	nres.Request.Trailer = resp.Request.Trailer
	nres.Request.RemoteAddr = resp.Request.RemoteAddr
	nres.Request.RequestURI = resp.Request.RequestURI

	// Request.TLS
	nres.Request.TLS = resp.Request.TLS

	nres.Request.Form = resp.Request.Form
	nres.Request.PostForm = resp.Request.PostForm
	// nres.Request.MultipartForm.Value = resp.Request.MultipartForm.Value
	// nres.Request.MultipartForm.File = resp.Request.MultipartForm.File
	// nres.Request.TLS = resp.Request.URL.TLS
	return
}
