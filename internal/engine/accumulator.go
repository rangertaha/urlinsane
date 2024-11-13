package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	log "github.com/sirupsen/logrus"
)

type accumulator struct {
	out    chan<- internal.Domain
	domain internal.Domain
}

func NewAccumulator(out chan<- internal.Domain, domain internal.Domain) internal.Accumulator {
	return &accumulator{
		out:    out,
		domain: domain,
	}
}

func (ac *accumulator) Add(domain internal.Domain) {
	log.WithFields(log.Fields{"domain": domain.String()}).Debug("Accumulator:Add()")
	ac.out <- domain
}

func (c *accumulator) Mkdir(dirpath, name string) (dir string, err error) {
	dir = filepath.Join(dirpath, name)
	if err = os.MkdirAll(dir, 0750); err != nil {
		return
	}
	return
}

func (c *accumulator) Mkfile(dir, name string, content []byte) (file string, err error) {
	file = filepath.Join(dir, name)
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = os.WriteFile(file, content, 0644)
		if err != nil {
			return
		}
	}
	return
}

func (c *accumulator) Domain() internal.Domain {
	return c.domain
}

func (c *accumulator) GetJson(key string) json.RawMessage {
	key = strings.ToLower(key)
	return c.domain.GetData(key)
}

func (c *accumulator) SetJson(key string, value json.RawMessage) {
	key = strings.ToLower(key)
	c.domain.SetData(key, value)
}

func (c *accumulator) Unmarshal(key string, v interface{}) (err error) {
	if data := c.GetJson(key); data != nil {

		Source := (*json.RawMessage)(&data)
		return json.Unmarshal(*Source, &v)

		// Source := (*json.RawMessage)(&data)
		// return json.Unmarshal(*Source, v)
	}

	return fmt.Errorf("nothing to retrive")
}

func (c *accumulator) GetMeta(key string) (data string) {
	return c.domain.GetMeta(key)
}

func (c *accumulator) SetMeta(key string, value string) {
	c.domain.SetMeta(key, value)
}

func (c *accumulator) Save(key string, data []byte) (err error) {
	return err
}

func (c *accumulator) Next() (err error) {
	c.Add(c.domain)
	return err
}

func (c *accumulator) Live(v ...bool) bool {
	c.domain.Live(v...)
	return c.domain.Live()
}
