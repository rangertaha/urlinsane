package influxdb

import (
	"fmt"

	"github.com/cybint/athena/interfaces"
	"github.com/cybint/athena/plugins/outputs"
)

type None struct {
	bot interfaces.Bot
}

func (p *None) Init(bot interfaces.Bot) error {
	p.bot = bot

	return nil
}

func (p *None) Run(diss interfaces.Disseminator) error {
	fmt.Println("output NONE")
	for metric := range diss.Metrics() {
		fmt.Println("OUTPUT", metric)
		fmt.Println("PROGRESS", p.bot.Timeline().Progress())
	}
	return nil
}

func (p *None) Train(diss interfaces.Disseminator) error {
	return p.Run(diss)
}

func (p *None) Test(diss interfaces.Disseminator) error {
	return p.Run(diss)
}

func (p *None) Stop() {

}

// Register the plugin
func init() {
	outputs.Add("none", func() interfaces.Output {
		return &None{}
	})
}
