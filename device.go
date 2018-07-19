package main

import (
	"github.com/DDRBoxman/mfi"
	"github.com/chbmuc/lirc"
)

type device struct {
	Name            string
	Power           power
	RequiredDevices []*device
	Commands []command
}

type command interface {
	Send()
}

type avinput interface {
}

type power interface {
	On()
	Off()
}

type irRemotePower struct {
	onCommand string
	offCommand string
	ir *lirc.Router
}

func (p *irRemotePower) On() {
	p.ir.Send(p.onCommand)
}

func (p *irRemotePower) Off() {
	p.ir.Send(p.offCommand)
}

type mfiPort struct {
	port   int
	client *mfi.MFIClient
}

func (p *mfiPort) On() {
	p.client.SetOutputEnabled(p.port, true)
}

func (p *mfiPort) Off() {
	p.client.SetOutputEnabled(p.port, false)
}

type lircCommand struct {
	command string
	ir *lirc.Router
}

func (l *lircCommand) Send() {
	l.ir.Send(l.command)
}