package status

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/spf13/viper"
	"image/png"
	"os"
)

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Sample struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Sample `json:"sample"`
}

type Description struct {
	Text string `json:"text"`
}

type Status struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favicon,omitempty"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
	PreviewsChat       bool        `json:"previewsChat"`
}

var (
	motd       string
	maxPlayers int
	favicon    string
)

func Init() error {
	motd = viper.GetString("motd")
	maxPlayers = viper.GetInt("max_players")
	return loadFavicon()
}

func loadFavicon() error {
	content, err := os.ReadFile("server-icon.png")
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		return err
	}
	if img.Bounds().Dx() != 64 || img.Bounds().Dy() != 64 {
		return errors.New("favicon must be 64x64")
	}

	var buf bytes.Buffer
	_, err = base64.NewEncoder(base64.StdEncoding, &buf).Write(content)
	if err != nil {
		return err
	}
	favicon = "data:image/png;base64," + buf.String()
	return nil
}

func GetStatus() *Status {
	return &Status{
		Version: Version{
			Name:     "1.20.4",
			Protocol: 765,
		},
		Players: Players{
			Max:    maxPlayers,
			Online: 5,
		},
		Description: Description{
			Text: motd,
		},
		Favicon:            favicon,
		EnforcesSecureChat: true,
		PreviewsChat:       true,
	}
}
