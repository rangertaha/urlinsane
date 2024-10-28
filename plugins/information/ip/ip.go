package ip

// // ipLookupFunc
// func ipLookupFunc(tr Result) (results []Result) {
// 	results = append(results, checkIP(tr))
// 	return
// }


import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "ip"

type Ipaddr struct {
	types []string
}

func (n *Ipaddr) Id() string {
	return CODE
}

func (n *Ipaddr) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Ipaddr) Name() string {
	return "Ip Address"
}

func (n *Ipaddr) Description() string {
	return "Domain IP addresses"
}

// func (n *Ipaddr) Fields() []string {
// 	return []string{}
// }

func (n *Ipaddr) Headers() []string {
	return []string{"Online", "IPv4", "IPv6"}
}

func (n *Ipaddr) Exec(in urlinsane.Typo) (out urlinsane.Typo) {

	in.Variant().Add("Online", true)
	in.Variant().Add("IPv4", "100.0.0.0")
	in.Variant().Add("IPv6", "100.0.0.0")
	in.Variant().Add("JSON", "{}")
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() urlinsane.Information {
		return &Ipaddr{
			types: []string{urlinsane.DOMAIN},
		}
	})
}
