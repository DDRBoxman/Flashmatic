package display

import (
	"fmt"
	"os"
)

type Icon struct {
	IconPath string `json:"icon_path"`
	KeyID    int    `json:"key_id"`
}

type Display struct {
	IconDir string
	Icons   []Icon
}

func MakeDisplay(iconDir string) (*Display, error) {
	if _, err := os.Stat(iconDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("icon path %s does not exist", iconDir)
	}

	display := &Display{
		IconDir: iconDir,
	}

	return display, nil
}

func (display *Display) AddIcon(keyID int, iconPath string) {
	display.Icons = append(display.Icons, Icon{
		IconPath: iconPath,
		KeyID:    keyID,
	})
}
