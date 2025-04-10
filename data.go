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
	Color FieldColor
}

type FieldColor string

const (
	FieldColorBlue    FieldColor = "1"
	FieldColorRed     FieldColor = "2"
	FieldColorOrange  FieldColor = "3"
	FieldColorPurple  FieldColor = "4"
	FieldColorGreen   FieldColor = "5"
	FieldColorSkyBlue FieldColor = "6"
	FieldColorPink    FieldColor = "7"
	FieldColorBrown   FieldColor = "8"
	FieldColorOlive   FieldColor = "9"
	FieldColorCyan    FieldColor = "10"
	FieldColorYellow  FieldColor = "11"
	FieldColorBlack   FieldColor = "12"
)

var colorMap = map[FieldColor]color.RGBA{
	FieldColorBlue:    {0x3B, 0x59, 0x98, 0xFF},
	FieldColorRed:     {0xDC, 0x39, 0x12, 0xFF},
	FieldColorOrange:  {0xFF, 0x99, 0x00, 0xFF},
	FieldColorPurple:  {0x99, 0x00, 0x99, 0xFF},
	FieldColorGreen:   {0x10, 0x96, 0x18, 0xFF},
	FieldColorSkyBlue: {0x00, 0x99, 0xC6, 0xFF},
	FieldColorPink:    {0xDD, 0x44, 0x77, 0xFF},
	FieldColorBrown:   {0x99, 0x66, 0x33, 0xFF},
	FieldColorOlive:   {0x66, 0xAA, 0x00, 0xFF},
	FieldColorCyan:    {0x00, 0xFF, 0xFF, 0xFF},
	FieldColorYellow:  {0xFF, 0xFF, 0x00, 0xFF},
	FieldColorBlack:   {0x00, 0x00, 0x00, 0xFF},
}

func FieldColorToRGBA(c FieldColor) (color.RGBA, bool) {
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
