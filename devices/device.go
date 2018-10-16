package devices

import (
	"github.com/DDRBoxman/mfi"
	"github.com/chbmuc/lirc"
	"github.com/tarm/serial"
	"log"
)

type Device struct {
	Name            string
	Power           Power
	RequiredDevices []*Device
	Commands        []Command
}

type Command interface {
	Send()
}

type Avinput interface {
}

type Power interface {
	On()
	Off()
}

type IrRemotePower struct {
	OnCommand  string
	OffCommand string
	Ir         *lirc.Router
}

func (p *IrRemotePower) On() {
	p.Ir.Send(p.OnCommand)
}

func (p *IrRemotePower) Off() {
	p.Ir.Send(p.OffCommand)
}

type MfiPort struct {
	Port   int
	Client *mfi.MFIClient
}

func (p *MfiPort) On() {
	p.Client.SetOutputEnabled(p.Port, true)
}

func (p *MfiPort) Off() {
	p.Client.SetOutputEnabled(p.Port, false)
}

type LircCommand struct {
	Command string
	Ir      *lirc.Router
}

func (l *LircCommand) Send() {
	l.Ir.Send(l.Command)
}

// Logitech HDMI Switch
// Uses ATEN VS0801H
// swmode default - disables detect switching useful for consoles like wii u
type AviorCommand struct {
	Command string
	Port    *serial.Port
}

func (a *AviorCommand) Send() {
	n, err := a.Port.Write([]byte(a.Command))
	if err != nil {
		log.Print("Failed to write to serial port ", err)
	}

	log.Printf("Wrote %d bytes to serial port \n", n)

	buf := make([]byte, 128)
	n, err = a.Port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])
}
