package actor

import "fmt"

type Printer struct {
	*Actor[string]
}

func NewPrinter() *Printer {
	p := &Printer{}
	p.Actor = New(func(s string) {
		fmt.Println(s)
	})
	return p
}

func (p *Printer) Print(s string) {
	p.Send(s)
}
