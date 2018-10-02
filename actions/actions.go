package actions

import (
	"log"
	"github.com/tarm/serial"
	"github.com/chbmuc/lirc"
	"github.com/DDRBoxman/mfi"
	"github.com/DDRBoxman/Flashmatic/devices"
	"time"
)

var deviceActions = map[int]*devices.Device{}

var commandActions = map[int]*devices.Command{}

func addCommandAction(command devices.Command, key int) {
	commandActions[key] = &command
}

func addDeviceAction(device devices.Device, key int) {
	deviceActions[key] = &device
}

var keychan = make(chan int, 100)

func SetupDevices() (chan int) {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 19200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	ir, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}

	Client, err := mfi.MakeMFIClient("10.42.42.12", "ubnt", "ubnt")
	if err != nil {
		log.Panic(err)
	}

	err = Client.Auth("ubnt", "ubnt")
	if err != nil {
		log.Panic(err)
	}

	// OSSC Scart input
	osscScartCommand := &devices.LircCommand{
		Command: "ossc 1",
		Ir:      ir,
	}

	hdmi3Command := &devices.LircCommand{
		Command: "BN59-01041A HDMI3",
		Ir:      ir,
	}

	hdmi2Command := &devices.LircCommand{
		Command: "BN59-01041A HDMI2",
		Ir:      ir,
	}

	off := devices.Device{
		Name: "Off",
	}

	tv := devices.Device{
		Name: "TV",
		Power: &devices.IrRemotePower{
			OnCommand:  "BN59-01041A PowerOn",
			OffCommand: "BN59-01041A PowerOff",
			Ir:         ir,
		},
	}

	ossc := devices.Device{
		Name: "OSSC",
		Power: &devices.MfiPort{
			Port:   4,
			Client: Client,
		},
		Commands: []devices.Command{
			hdmi3Command,
		},
	}

	hydra := devices.Device{
		Name: "Hydra",
		Power: &devices.MfiPort{
			Port:   6,
			Client: Client,
		},
		Commands: []devices.Command{
			osscScartCommand,
		},
	}

	nes := devices.Device{
		Name: "NES",
		Power: &devices.MfiPort{
			Port:   1,
			Client: Client,
		},
		RequiredDevices: []*devices.Device{&tv, &ossc, &hydra},
	}

	segaCD := devices.Device{
		Name: "Sega CD",
		Power: &devices.MfiPort{
			Port:   3,
			Client: Client,
		},
	}

	genesis := devices.Device{
		Name: "Genesis",
		Power: &devices.MfiPort{
			Port:   2,
			Client: Client,
		},
		RequiredDevices: []*devices.Device{&tv, &ossc, &hydra, &segaCD},
	}

	n64 := devices.Device{
		Name: "N64",
		Power: &devices.MfiPort{
			Port:   7,
			Client: Client,
		},
		RequiredDevices: []*devices.Device{&tv, &ossc, &hydra},
	}

	ps1 := devices.Device{
		Name: "PS1",
		Power: &devices.MfiPort{
			Port:   8,
			Client: Client,
		},
		RequiredDevices: []*devices.Device{&tv, &ossc, &hydra},
	}

	hdmiSwitch := devices.Device{
		Commands: []devices.Command{hdmi2Command},
	}

	gamecube := devices.Device{
		Name:            "GameCube",
		RequiredDevices: []*devices.Device{&tv, &hdmiSwitch},
		Power: &devices.MfiPort{
			Port:   5,
			Client: Client,
		},
		Commands: []devices.Command{&devices.AviorCommand{
			Port:    s,
			Command: "sw i02\r",
		}},
	}

	snes := devices.Device{
		Name:            "SNES",
		RequiredDevices: []*devices.Device{&tv, &hdmiSwitch},
		Commands: []devices.Command{&devices.AviorCommand{
			Port:    s,
			Command: "sw i01\r",
		}},
	}

	xbone := devices.Device{
		Name:            "XBox One",
		RequiredDevices: []*devices.Device{&tv, &hdmiSwitch},
		Commands: []devices.Command{&devices.AviorCommand{
			Port:    s,
			Command: "sw i03\r",
		}},
	}

	ps4 := devices.Device{
		Name:            "PS4",
		RequiredDevices: []*devices.Device{&tv, &hdmiSwitch},
		Commands: []devices.Command{&devices.AviorCommand{
			Port:    s,
			Command: "sw i05\r",
		}},
	}

	ps3 := devices.Device{
		Name: "PS3",
		RequiredDevices: []*devices.Device{&tv, &hdmiSwitch},
		Commands: []devices.Command{&devices.AviorCommand{
			Port:    s,
			Command: "sw i06\r",
		}},
	}

	addDeviceAction(off, 10)

	addDeviceAction(nes, 4)
	addDeviceAction(genesis, 3)
	addDeviceAction(n64, 2)
	addDeviceAction(ps1, 1)
	addDeviceAction(gamecube, 0)
	addDeviceAction(snes, 9)
	addDeviceAction(xbone, 6)
	addDeviceAction(ps4, 8)
	addDeviceAction(ps3,7)

	volup := devices.LircCommand{
		Command: "BN59-01041A V_UP",
		Ir:      ir,
	}

	voldown := devices.LircCommand{
		Command: "BN59-01041A V_DOWN",
		Ir:      ir,
	}

	addCommandAction(&volup, 13)
	addCommandAction(&voldown, 14)

	go processKeyEvents()

	return keychan
}

func processKeyEvents() {
	for {
		select {
		case event := <-keychan:
			if device, ok := deviceActions[event]; ok {
				go updateToDeviceDesiredState(device)
			}

			if cmd, ok := commandActions[event]; ok {
				(*cmd).Send()
			}
		}
	}
}

var activeDevices = []*devices.Device{}

func updateToDeviceDesiredState(d *devices.Device) {
	desiredActive := []*devices.Device{}
	for _, device := range d.RequiredDevices {
		desiredActive = append(desiredActive, device)
	}

	desiredActive = append(desiredActive, d)

	desiredInactive := []*devices.Device{}
	for _, s1 := range activeDevices {
		found := false
		for _, s2 := range desiredActive {
			if s1 == s2 {
				found = true
				break
			}
		}

		if !found {
			desiredInactive = append(desiredInactive, s1)
		}
	}

	activeDevices = desiredActive

	for _, device := range desiredActive {
		if device.Power != nil {
			log.Println("Turning on: ", device.Name)
			device.Power.On()
		}
	}

	for _, device := range desiredInactive {
		if device.Power != nil {
			log.Println("Turning off: ", device.Name)
			device.Power.Off()
		}
	}

	time.Sleep(5 * time.Second)

	// Send commands for device
	sendCommands(d)
}

func sendCommands(d *devices.Device) {
	for _, command := range d.Commands {
		command.Send()
	}

	for _, device := range d.RequiredDevices {
		sendCommands(device)
	}
}
