package main

import (
	"log"
	_ "image/gif"
	_ "image/png"
	"time"
	"image/gif"
	"github.com/DDRBoxman/streamdeck-go"
	"github.com/DDRBoxman/mfi"
	"github.com/chbmuc/lirc"
)

var animatedKeys = []*animatedKey{}

type animatedKey struct {
	gif          *gif.GIF
	currentFrame int
	key          int
	lastDraw     time.Time
}

type eventType string

const ButtonUpEvent eventType = "buttonUp"
const ButtonDownEvent eventType = "buttonDown"

const backlightTimeout = 15 * time.Second

var streamDeck streamdeck.StreamDeck

func main() {
	decks := streamdeck.FindDecks()
	if len(decks) < 1 {
		log.Panic("No streamdeck found.")
	}

	streamDeck = decks[0]

	err := streamDeck.Open()
	if err != nil {
		log.Panic(err)
	}

	clearKeys(streamDeck)

	setupDevices()

	streamDeck.SetBrightness(100)

	go readKeysLoop(streamDeck)

	backlightTimer := time.NewTimer(backlightTimeout)

	backlightOff := false

	for {
		select {
		case event := <-eventChan:
			switch event.eventType {
			case ButtonUpEvent:
				if backlightOff {
					backlightTimer.Reset(backlightTimeout)
					streamDeck.SetBrightness(100)
					backlightOff = false
				} else {
					if device, ok := deviceIcons[event.id]; ok {
						go updateToDeviceDesiredState(device)
					}
				}
			}
		case <-backlightTimer.C:
			streamDeck.SetBrightness(0)
			backlightOff = true
		default:
			animateKeys(streamDeck)
			time.Sleep(150)
		}
	}
}

func setupDevices() {
	ir, err := lirc.Init("/var/run/lirc/lircd")
	if err != nil {
		panic(err)
	}

	client, err := mfi.MakeMFIClient("10.42.42.12", "ubnt", "ubnt")
	if err != nil {
		log.Panic(err)
	}

	err = client.Auth("ubnt", "ubnt")
	if err != nil {
		log.Panic(err)
	}

	// OSSC Scart input
	osscScartCommand := &lircCommand{
		command: "ossc 1",
		ir: ir,
	}

	hdmi3Command := &lircCommand{
		command: "BN59-01041A HDMI3",
		ir: ir,
	}

	off := device{
		Name: "Off",
	}

	tv := device{
		Name: "TV",
		Power: &irRemotePower{
			onCommand: "BN59-01041A PowerOn",
			offCommand: "BN59-01041A PowerOff",
			ir: ir,
		},
	}

	ossc := device{
		Name: "OSSC",
		Power: &mfiPort{
			port:   4,
			client: client,
		},
		Commands: []command {
			hdmi3Command,
		},
	}

	hydra := device{
		Name: "Hydra",
		Power: &mfiPort{
			port:   6,
			client: client,
		},
		Commands: []command {
			osscScartCommand,
		},
	}

	nes := device{
		Name: "NES",
		Power: &mfiPort{
			port:   1,
			client: client,
		},
		RequiredDevices: []*device{&tv, &ossc, &hydra},
	}

	segaCD := device{
		Name: "Sega CD",
		Power: &mfiPort{
			port:   3,
			client: client,
		},
	}

	genesis := device{
		Name: "Genesis",
		Power: &mfiPort{
			port:   2,
			client: client,
		},
		RequiredDevices: []*device{&tv, &ossc, &hydra, &segaCD},
	}

	n64 := device{
		Name: "N64",
		Power: &mfiPort{
			port:   7,
			client: client,
		},
		RequiredDevices: []*device{&tv, &ossc, &hydra},
	}

	ps1 := device{
		Name: "PS1",
		Power: &mfiPort{
			port:   8,
			client: client,
		},
		RequiredDevices: []*device{&tv,  &ossc, &hydra},
	}

	gamecube := device{
		Name: "GameCube",
		Power: &mfiPort{
			port:   8,
			client: client,
		},
	}

	makeDeviceAnimatedIcon(off, "./assets/icons/kirby_sleeping.gif", 10)

	makeDeviceIcon(streamDeck, nes, "./assets/icons/consoles/Nes.gif", 4)
	makeDeviceIcon(streamDeck, genesis, "./assets/icons/consoles/MegaCd01.gif", 3)
	makeDeviceIcon(streamDeck, n64, "./assets/icons/consoles/N64.gif", 2)
	makeDeviceIcon(streamDeck, ps1, "./assets/icons/consoles/Psone.gif", 1)
	makeDeviceIcon(streamDeck, gamecube, "./assets/icons/Ngc_Violet03.png", 0)

  /*
	writeIconToKey(deck, 0, "./assets/icons/consoles/Ps2_07.gif")
	writeIconToKey(deck, 9, "./assets/icons/consoles/DC03.gif")
	writeIconToKey(deck, 8, "./assets/icons/Ngc_Violet03.png")
	writeIconToKey(deck, 7, "./assets/icons/consoles/Snes2.gif")*/
}

var deviceIcons = map[int]*device{}

func makeDeviceIcon(deck streamdeck.StreamDeck, device device, iconPath string, key int) {
	writeIconToKey(deck, key, iconPath)
	deviceIcons[key] = &device
}

func makeDeviceAnimatedIcon(device device, iconPath string, key int) {
	err := setIconAnimated(key, iconPath)
	if err != nil {
		log.Panic(err)
	}
	deviceIcons[key] = &device
}

var activeDevices = []*device{}

func updateToDeviceDesiredState(d *device) {
	desiredActive := []*device{d}
	for _, device := range d.RequiredDevices {
		desiredActive = append(desiredActive, device)
	}

	desiredInactive := []*device{}
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

func sendCommands(d *device) {
	for _, command := range d.Commands {
		command.Send()
	}

	for _, device := range d.RequiredDevices {
		sendCommands(device)
	}
}