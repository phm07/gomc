package textcomponent

import (
	"encoding/json"
	"gomc/src/nbt"
	"strings"
)

const (
	ColorBlack       = "black"
	ColorDarkBlue    = "dark_blue"
	ColorDarkGreen   = "dark_green"
	ColorDarkAqua    = "dark_aqua"
	ColorDarkRed     = "dark_red"
	ColorPurple      = "dark_purple"
	ColorGold        = "gold"
	ColorGray        = "gray"
	ColorDarkGray    = "dark_gray"
	ColorBlue        = "blue"
	ColorGreen       = "green"
	ColorAqua        = "aqua"
	ColorRed         = "red"
	ColorLightPurple = "light_purple"
	ColorYellow      = "yellow"
	ColorWhite       = "white"
)

type Component struct {
	Text          string      `json:"text"`
	Color         string      `json:"color,omitempty"`
	Bold          bool        `json:"bold,omitempty"`
	Italic        bool        `json:"italic,omitempty"`
	Underlined    bool        `json:"underlined,omitempty"`
	Strikethrough bool        `json:"strikethrough,omitempty"`
	Obfuscated    bool        `json:"obfuscated,omitempty"`
	Insertion     string      `json:"insertion,omitempty"`
	Extra         []Component `json:"extra,omitempty"`
}

func New(text string) *Component {
	return &Component{Text: text}
}

func (c *Component) SetColor(color string) *Component {
	c.Color = color
	return c
}

func (c *Component) SetBold(bold bool) *Component {
	c.Bold = bold
	return c
}

func (c *Component) SetItalic(italic bool) *Component {
	c.Italic = italic
	return c
}

func (c *Component) SetUnderlined(underlined bool) *Component {
	c.Underlined = underlined
	return c
}

func (c *Component) SetStrikethrough(strikethrough bool) *Component {
	c.Strikethrough = strikethrough
	return c
}

func (c *Component) SetObfuscated(obfuscated bool) *Component {
	c.Obfuscated = obfuscated
	return c
}

func (c *Component) SetInsertion(insertion string) *Component {
	c.Insertion = insertion
	return c
}

func (c *Component) AddExtra(extra *Component) *Component {
	c.Extra = append(c.Extra, *extra)
	return c
}

func (c *Component) Plain() string {
	var sb strings.Builder
	sb.WriteString(c.Text)
	for _, e := range c.Extra {
		sb.WriteString(e.Plain())
	}
	return sb.String()
}

func (c *Component) MarshalJSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Component) MarshalNBT() nbt.Tag {
	return nbt.Marshal(*c)
}
