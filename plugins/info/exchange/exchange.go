package exchange

import (
	"fmt"

	"github.com/cybint/athena/interfaces"
	"github.com/cybint/athena/plugins/outputs"
)

type Exchange struct {
	exc interfaces.Exchange
}

func (p *Exchange) Init(bot interfaces.Bot) error {
	p.exc = bot.Exc()
	return nil
}

func (p *Exchange) Run(diss interfaces.Disseminator) error {
	fmt.Println("Output....")
	// change to diss.Orders()
	// for order := range diss.Metrics() {
	// 	p.exc.Submit(order)
	// }
	return nil
}

func (p *Exchange) Train(diss interfaces.Disseminator) error {
	// p.exc.Mode(db.TRAINING)
	return p.Run(diss)
}

func (p *Exchange) Test(diss interfaces.Disseminator) error {
	// p.exc.Mode(db.TESTING)
	return p.Run(diss)
}

func (p *Exchange) Stop() {}

// Register the plugin
func init() {
	outputs.Add("exchange", func() interfaces.Output {
		return &Exchange{
			// Set the default values

		}
	})
}
