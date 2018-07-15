package main

import (
	"log"
	"github.com/DDRBoxman/streamdeck-go"
)

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
