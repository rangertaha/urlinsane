package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigFile struct {
	Name       string   `yaml:"name"`
	Age        int      `yaml:"age"`
	DnsServers []string `yaml:"dnsServers"`
}

type AppDir struct {
	config  ConfigFile
	homeDir string
	appDir  string
	appCfg  string
}

var defaultConfig = ConfigFile{
	Name: "John Doe",
	Age:  30,
	DnsServers: []string{
		"127.0.0.53",
		"100.64.100.1",
		"8.8.8.8",         // Google
		"8.8.4.4",         // Google
		"1.1.1.1",         // Cloudflare
		"1.0.0.1",         // Cloudflare
		"1.1.1.2",         // Cloudflare  malware blocking
		"1.0.0.2",         // Cloudflare  malware blocking
		"1.1.1.3",         // Cloudflare  adult blocking
		"1.0.0.3",         // Cloudflare  adult blocking
		"9.9.9.9",         // Quad9
		"149.112.112.112", // Quad9
		"185.228.168.9",   // Cleanbrowsing
		"185.228.169.9",   // Cleanbrowsing
		"8.26.56.26",      // Comodo Secure DNS
		"8.20.247.20",     // Comodo Secure DNS
	},
}

func NewAppDir() (*AppDir, error) {
	dir := &AppDir{}
	if err := dir.Init(); err != nil {
		return nil, err
	}

	return dir, nil
}

func (c *AppDir) Init() (err error) {
	if c.homeDir == "" {
		if c.homeDir, err = os.UserHomeDir(); err != nil {
			if c.homeDir, err = os.Getwd(); err != nil {
				c.homeDir = ""
			}
		}
	}
	if c.appDir == "" {
		c.appDir = filepath.Join(c.homeDir, ".urlinsane")
	}
	if c.appCfg == "" {
		c.appCfg = filepath.Join(c.appDir, "config.yml")
	}

	return c.getOrCreate()
}

func (c *AppDir) getOrCreate() (err error) {
	// Create app directory in user's home directory or local directory
	err = os.MkdirAll(c.appDir, 0750)
	if err != nil {
		return err
	}

	// Create or load config file
	_, err = os.Stat(c.appCfg)
	if os.IsNotExist(err) {
		c.SaveConfig(defaultConfig)
	} else {
		c.config = c.LoadConfig()
	}
	return nil
}

func (c *AppDir) SaveConfig(config ConfigFile) {
	file, err := os.Create(c.appCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		log.Fatal(err)
	}
}

func (c *AppDir) LoadConfig() ConfigFile {
	file, err := os.Open(c.appCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// var data Data
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&c.config); err != nil {
		log.Fatal(err)
	}

	return c.config
}
