package streamdeck

import (
	"github.com/DDRBoxman/streamdeck-go"
	"log"
	"time"
)

type eventType string

const ButtonUpEvent eventType = "buttonUp"
const ButtonDownEvent eventType = "buttonDown"

const backlightTimeout = 15 * time.Second

type event struct {
	eventType eventType
	id        int
}

var eventChan = make(chan event, 100)

var keyData = make([]byte, 16)

func readKeysLoop(deck streamdeck.StreamDeck) {
	data := make([]byte, 255)
	for {
		size, err := deck.Device.Read(data)
		if err != nil {
			log.Println(err)
			continue
		}

		if size != 17 {
			continue
		}

		for i, state := range data {
			// Skip first byte
			if i == 0 || i > 15 {
				continue
			}

			if keyData[i] != state {
				if state == 1 {
					eventChan <- event{
						eventType: ButtonDownEvent,
						id:        i - 1,
					}
				} else {
					eventChan <- event{
						eventType: ButtonUpEvent,
						id:        i - 1,
					}
				}
			}
		}

		copy(keyData, data)
	}
}
