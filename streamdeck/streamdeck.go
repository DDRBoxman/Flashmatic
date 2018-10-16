package streamdeck

import (
	"github.com/DDRBoxman/Flashmatic/display"
	"github.com/DDRBoxman/streamdeck-go"
	"log"
	"path/filepath"
	"time"
)

var streamDeck streamdeck.StreamDeck

func StartStreamdeck(display *display.Display, keychan chan int) {
	decks := streamdeck.FindDecks()
	if len(decks) < 1 {
		log.Println("No streamdeck found.")
	} else {
		streamDeck = decks[0]

		err := streamDeck.Open()
		if err != nil {
			log.Panic(err)
		}

		clearKeys(streamDeck)

		setupIcons(display)

		streamDeck.SetBrightness(100)

		go readKeysLoop(streamDeck)

		backlightTimer := time.NewTimer(backlightTimeout)

		backlightOff := false

		animateTimer := time.NewTicker(150 * time.Millisecond)

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
						keychan <- event.id
					}
				}
			case <-backlightTimer.C:
				streamDeck.SetBrightness(0)
				backlightOff = true
			case <-animateTimer.C:
				animateKeys(streamDeck)
			}
		}
	}
}

func setupIcons(display *display.Display) {
	for _, icon := range display.Icons {
		if filepath.Ext(icon.IconPath) == ".gif" {
			setIconAnimated(icon.KeyID, filepath.Join(display.IconDir, icon.IconPath))
		} else {
			writeIconToKey(streamDeck, icon.KeyID, filepath.Join(display.IconDir, icon.IconPath))
		}
	}
}
