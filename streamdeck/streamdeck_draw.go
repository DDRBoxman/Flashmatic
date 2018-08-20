package streamdeck

import (
	"github.com/DDRBoxman/streamdeck-go"
	"time"
	"os"
	"image/gif"
	"image"
	"image/draw"
	"github.com/disintegration/imaging"
)

var animatedKeys = []*animatedKey{}

type animatedKey struct {
	gif          *gif.GIF
	currentFrame int
	key          int
	lastDraw     time.Time
}

func clearKeys(deck streamdeck.StreamDeck) {
	for i := 0; i<15; i++ {
		dest := image.NewRGBA(image.Rect(0, 0, streamdeck.ICON_SIZE, streamdeck.ICON_SIZE))
		deck.WriteImageToKey(dest, i)
	}
}

func animateKeys(deck streamdeck.StreamDeck) {
	for _, key := range animatedKeys {
		if len(key.gif.Image) <= 1 {
			continue
		}
		keyDuration := time.Duration(key.gif.Delay[key.currentFrame]*10)*time.Millisecond
		if time.Now().Sub(key.lastDraw) > keyDuration {
			key.lastDraw = time.Now()
			key.currentFrame++
			if key.currentFrame >= len(key.gif.Image) {
				key.currentFrame = 0
			}
			writePalletedFrameToKey(deck, key.key, key.gif.Image[key.currentFrame])
		}
	}
}

func setIconAnimated(key int, path string) error {
	iconFile, _ := os.Open(path)
	defer iconFile.Close()

	gifData, err := gif.DecodeAll(iconFile)
	if err != nil {
		return err
	}

	animatedKeys = append(animatedKeys, &animatedKey{
		key:          key,
		gif:          gifData,
		currentFrame: 0,
		lastDraw: time.Now(),
	})

	return nil
}

func writePalletedFrameToKey(deck streamdeck.StreamDeck, key int, frame *image.Paletted) {
	dest := image.NewRGBA(image.Rect(0, 0, streamdeck.ICON_SIZE, streamdeck.ICON_SIZE))
	draw.Draw(dest, frame.Bounds(), frame, image.ZP, draw.Over)

	deck.WriteImageToKey(dest, key)
}

func writeIconToKey(deck streamdeck.StreamDeck, key int, path string) {
	iconFile, _ := os.Open(path)
	defer iconFile.Close()
	icon, _, _ := image.Decode(iconFile)

	dest := image.NewRGBA(image.Rect(0, 0, streamdeck.ICON_SIZE, streamdeck.ICON_SIZE))
	scaled := imaging.Resize(icon, streamdeck.ICON_SIZE, 0, imaging.NearestNeighbor)
	draw.Draw(dest, scaled.Bounds(), scaled, image.ZP, draw.Over)

	deck.WriteImageToKey(dest, key)
}
