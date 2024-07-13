package influxdb

import (
	"github.com/cybint/athena/db"
	"github.com/cybint/athena/interfaces"
	"github.com/cybint/athena/plugins"
	"github.com/cybint/athena/plugins/outputs"
	influxdb2 "github.com/influxdata/influxdb-client-go"
)

type InfluxDB struct {
	plugins.Plugin
	token  string
	Url    string `hcl:"url"`
	Org    string `hcl:"org"`
	Batch  uint   `hcl:"batch,optional"`
	bucket string

	client influxdb2.Client
	bot    interfaces.Bot
}

func (p *InfluxDB) Init(bot interfaces.Bot) error {
	p.bot = bot
	if bot.Status() == db.TESTING {
		p.bucket = "testing"
	}
	if bot.Status() == db.TRAINING {
		p.bucket = "training"
	}

	p.client = influxdb2.NewClientWithOptions(p.Url, p.token,
		influxdb2.DefaultOptions().SetBatchSize(p.Batch))

	return nil
}

func (p *InfluxDB) Run(diss interfaces.Disseminator) error {
	writeAPI := p.client.WriteAPI(p.Org, p.bucket)
	for metric := range diss.Metrics() {

		// create point
		p := influxdb2.NewPoint(
			metric.Name(),
			metric.Tags(),
			metric.Fields(),
			metric.Time())

		// write asynchronously
		writeAPI.WritePoint(p)
	}

	// Force all unwritten data to be sent
	writeAPI.Flush()

	return nil
}

func (p *InfluxDB) Train(diss interfaces.Disseminator) error {
	return p.Run(diss)
}

func (p *InfluxDB) Test(diss interfaces.Disseminator) error {
	return p.Run(diss)
}

func (p *InfluxDB) Stop() {
	// p.client.Close()
}

// Register the plugin
func init() {
	outputs.Add("influxdb", func() interfaces.Output {
		return &InfluxDB{
			// Set the default values
			token:  "WtFkD1OilUFKlTZJ2WX1FgiqB2w8bGOVB4lr-Io_Jzu9ZW0YXXBZkgKVelbs6pvzFFmmRvmqpOYDnXACpimNwQ==",
			Url:    "http://127.0.0.1:8086",
			Org:    "athena",
			bucket: "trading",
			Batch:  1000,
		}
	})
}
