package main

import (
	"sync"
	"github.com/DDRBoxman/Flashmatic/streamdeck"
		"github.com/DDRBoxman/Flashmatic/web"
	"github.com/DDRBoxman/Flashmatic/display"
	"log"
	"github.com/DDRBoxman/Flashmatic/actions"
)

func main() {
	keychan := actions.SetupDevices()

	display, err := display.MakeDisplay("./assets/icons")
	if err != nil {
		log.Fatal(err)
	}
	display.AddIcon(10, "kirby_sleeping.gif")
	display.AddIcon(4, "consoles/Nes.gif")
	display.AddIcon(3, "consoles/MegaCd01.gif")
	display.AddIcon(2, "consoles/N64.gif")
	display.AddIcon(1, "consoles/Psone.gif")
	display.AddIcon(0, "Ngc_Violet03.png")
	display.AddIcon(9, "consoles/Snes2.gif")
	display.AddIcon(6, "xbone.png")
	display.AddIcon(8, "ps4.png")
	display.AddIcon(7, "ps3.png")
	display.AddIcon(13, "vol_up.png")
	display.AddIcon(14, "vol_down.png")

	go web.StartServer(display, keychan)

	go streamdeck.StartStreamdeck(display, keychan)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
