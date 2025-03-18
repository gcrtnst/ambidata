package ambidata

import "time"

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
