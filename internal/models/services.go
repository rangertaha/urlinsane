package models

type IP struct {
	Address  string   `json:"ip,omitempty"`
	Ports    []Port   `json:"ports,omitempty"`
	Location Location `json:"geo,omitempty"`
}

type Service struct {
	Name   string `json:"name,omitempty"`
	Banner string `json:"banner,omitempty"`
}

type Port struct {
	Number  int     `json:"number,omitempty"`
	State   string  `json:"state,omitempty"`
	Service Service `json:"service,omitempty"`
}
