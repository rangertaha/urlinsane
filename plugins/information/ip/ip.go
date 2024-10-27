package ip

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/information"
)

const CODE = "ip"

type Ipaddr struct {
	id    int
	types []string
}

func (n *Ipaddr) Id() int {
	return n.id
}

func (n *Ipaddr) Code() string {
	return CODE
}

func (n *Ipaddr) Name() string {
	return "IP Address"
}

func (n *Ipaddr) Description() string {
	return "Domain IP addresses"
}

func (n *Ipaddr) Fields() []string {
	return []string{}
}

func (n *Ipaddr) Headers() []string {
	return []string{}
}

func (n *Ipaddr) Exec(in urlinsane.Typo) (out urlinsane.Typo) {
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() urlinsane.Information {
		return &Ipaddr{
			id:    0,
			types: []string{information.DOMAINS},
		}
	})
}
