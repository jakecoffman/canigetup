package pi

import (
	"github.com/stianeikeland/go-rpio/v4"
	"log"
)

type Pi struct {
	pin rpio.Pin
}

func New() *Pi {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	pin := rpio.Pin(15)
	pin.Output()
	return &Pi{
		pin: pin,
	}
}

func (p *Pi) Close() {
	err := rpio.Close()
	if err != nil {
		log.Println(err)
	}
}

func (p *Pi) High() {

}

func (p *Pi) Low() {

}
