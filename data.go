package ambidata

import (
	"image/color"
	"time"
)

type ChannelAccess struct {
	ChannelInfo
	ReadKey  string
	WriteKey string
}

func (ch *ChannelAccess) ToLv1() ChannelAccessLv1 {
	return ChannelAccessLv1{
		Ch:       ch.Ch,
		WriteKey: ch.WriteKey,
	}
}

type ChannelAccessLv1 struct {
	Ch       string
	WriteKey string
}

type ChannelInfo struct {
	Ch         string
	User       string
	Created    time.Time
	Modified   time.Time
	LastPost   time.Time
	Charts     int
	DataPerDay int
	DCh        bool
	ChName     string
	ChDesc     string
	D1         FieldInfo
	D2         FieldInfo
	D3         FieldInfo
	D4         FieldInfo
	D5         FieldInfo
	D6         FieldInfo
	D7         FieldInfo
	D8         FieldInfo
	Loc        Maybe[Location]
	PhotoID    string
	DevKeys    []string
	Bd         string
	LastData   LastData
}

type FieldInfo struct {
	Name  string
	Color Color
}

type Color string

const (
	ColorBlue    Color = "1"
	ColorRed     Color = "2"
	ColorOrange  Color = "3"
	ColorPurple  Color = "4"
	ColorGreen   Color = "5"
	ColorSkyBlue Color = "6"
	ColorPink    Color = "7"
	ColorBrown   Color = "8"
	ColorOlive   Color = "9"
	ColorCyan    Color = "10"
	ColorYellow  Color = "11"
	ColorBlack   Color = "12"
)

var colorMap = map[Color]color.RGBA{
	ColorBlue:    {0x3B, 0x59, 0x98, 0xFF},
	ColorRed:     {0xdc, 0x39, 0x12, 0xFF},
	ColorOrange:  {0xFF, 0x99, 0x00, 0xFF},
	ColorPurple:  {0x99, 0x00, 0x99, 0xFF},
	ColorGreen:   {0x10, 0x96, 0x18, 0xFF},
	ColorSkyBlue: {0x00, 0x99, 0xC6, 0xFF},
	ColorPink:    {0xDD, 0x44, 0x77, 0xFF},
	ColorBrown:   {0x99, 0x66, 0x33, 0xFF},
	ColorOlive:   {0x66, 0xAA, 0x00, 0xFF},
	ColorCyan:    {0x00, 0xFF, 0xFF, 0xFF},
	ColorYellow:  {0xFF, 0xFF, 0x00, 0xFF},
	ColorBlack:   {0x00, 0x00, 0x00, 0xFF},
}

func ColorToRGBA(c Color) (color.RGBA, bool) {
	rgba, ok := colorMap[c]
	return rgba, ok
}

type LastData struct {
	Data
	ID string
}

type Data struct {
	Created time.Time
	D1      Maybe[float64]
	D2      Maybe[float64]
	D3      Maybe[float64]
	D4      Maybe[float64]
	D5      Maybe[float64]
	D6      Maybe[float64]
	D7      Maybe[float64]
	D8      Maybe[float64]
	Loc     Maybe[Location]
	Cmnt    string
	Hide    bool
}

type Location struct {
	Lat, Lng float64
}

type Maybe[T any] struct {
	V  T
	OK bool
}

func Just[T any](v T) Maybe[T] {
	return Maybe[T]{V: v, OK: true}
}
