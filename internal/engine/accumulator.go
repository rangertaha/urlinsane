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
	cfg    internal.Config
	kv     internal.Database
	log    log.Entry
	dir    string
}

func NewAccumulator(out chan<- internal.Domain, domain internal.Domain, conf internal.Config) internal.Accumulator {
	logger := log.WithFields(log.Fields{"domain": domain.String(), "func": "accumulator"})
	var dir string
	if conf.AssetDir() != "" {
		dir = filepath.Join(conf.AssetDir(), domain.String())
	}

	return &accumulator{
		out:    out,
		domain: domain,
		cfg:    conf,
		log:    *logger,
		kv:     conf.Database(),
		dir:    dir,
	}
}

func (ac *accumulator) Add(domain internal.Domain) {
	ac.log.Debug("Passing domain to next")
	ac.out <- domain
}

// func (c *accumulator) Mkdir(dirpath, name string) (dir string, err error) {
// 	dir = filepath.Join(dirpath, name)
// 	if err = os.MkdirAll(dir, 0750); err != nil {
// 		return
// 	}
// 	return
// }

// func (c *accumulator) Mkfile(dir, name string, content []byte) (file string, err error) {
// 	file = filepath.Join(dir, name)
// 	_, err = os.Stat(file)
// 	if os.IsNotExist(err) {
// 		err = os.WriteFile(file, content, 0644)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

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

func (c *accumulator) Marshal(key string, v interface{}) (err error) {
	var records []byte
	records, err = json.Marshal(v)

	if err != nil {
		log.Error(err)
		return err
	}
	key = fmt.Sprintf("%s:%s", c.domain.String(), key)
	// log.Debug(key, "write:", string(records))
	c.SetJson(key, records)
	// log.Debug(key, "read:", string(c.GetJson(key)))

	return
}

func (c *accumulator) Unmarshal(key string, v interface{}) (err error) {
	if data := c.GetJson(key); data != nil {
		return json.Unmarshal(data, v)
	}
	c.log.Error("nothing to unmarshal")

	return fmt.Errorf("nothing to unmarshal")
}

func (c *accumulator) GetMeta(key string) (data string) {
	return c.domain.GetMeta(key)
}

func (c *accumulator) SetMeta(key string, value string) {
	c.domain.SetMeta(key, value)
}

func (c *accumulator) Metadata() map[string]string {
	return c.domain.Meta()
}

func (c *accumulator) Save(name string, data []byte) (err error) {
	if _, err = os.Stat(c.dir); os.IsNotExist(err) {
		log.Debug(err, "creating..")
		if err = os.MkdirAll(c.dir, 0750); err != nil {
			log.Error(err)
			return
		}
	}

	file := filepath.Join(c.dir, name)
	// if _, err = os.Stat(file); os.IsNotExist(err) {
	log.Debugf("creating %s", file)
	if err = os.WriteFile(file, data, 0644); err != nil {
		log.Error(err)
		return
		// }
	}

	return err
}

func (c *accumulator) Next() (err error) {
	c.Add(c.domain)
	return err
}

func (c *accumulator) Live(v ...bool) bool {
	return c.domain.Live(v...)
}

func (c *accumulator) Cached(v ...bool) bool {
	return c.domain.Cached(v...)
}
