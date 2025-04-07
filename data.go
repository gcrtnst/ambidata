package ambidata

import (
	"image/color"
	"time"
)

// ChannelAccess はチャネルへのアクセス情報を保持する構造体です。
// チャネル情報に加えて、読み取りキーと書き込みキーを含みます。
type ChannelAccess struct {
	ChannelInfo
	ReadKey  string
	WriteKey string
}

// ToLv1 はChannelAccessからChannelAccessLv1への変換を行います。
// レベル1のチャネルアクセス情報（チャネルIDと書き込みキーのみ）を返します。
func (ch *ChannelAccess) ToLv1() ChannelAccessLv1 {
	return ChannelAccessLv1{
		Ch:       ch.Ch,
		WriteKey: ch.WriteKey,
	}
}

// ChannelAccessLv1 は簡易的なチャネルアクセス情報を保持する構造体です。
// チャネルIDと書き込みキーのみを含む軽量な構造体です。
type ChannelAccessLv1 struct {
	Ch       string
	WriteKey string
}

// ChannelInfo はチャネルの詳細情報を保持する構造体です。
// チャネルの基本情報、フィールド情報、最新データなどを含みます。
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

// FieldInfo はデータフィールドの情報を保持する構造体です。
// フィールド名と色情報を含みます。
type FieldInfo struct {
	Name  string
	Color FieldColor
}

// FieldColor はフィールドの色を表す型です。
// 文字列として保存され、[FieldColorToRGBA] 関数で RGBA 値に変換できます。
type FieldColor string

// 以下の定数はフィールドに使用できる色を定義しています。
// これらの色は Ambient のグラフやUIで使用されます。
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

// FieldColorToRGBA は指定された [Color] を RGBA 値に変換します。
// 変換が成功した場合は RGBA 値と true を返し、失敗した場合はゼロ値と false を返します。
func FieldColorToRGBA(c FieldColor) (color.RGBA, bool) {
	rgba, ok := colorMap[c]
	return rgba, ok
}

// LastData はチャネルの最新データを表す構造体です。
// 通常のデータに加えて、データのIDを含みます。
type LastData struct {
	Data
	ID string
}

// Data はチャネルに保存されるデータポイントを表す構造体です。
// タイムスタンプ、8つの数値フィールド、位置情報、コメント、表示/非表示状態を含みます。
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

// Location は位置情報を表す構造体です。
// 緯度と経度の座標を含みます。
type Location struct {
	Lat, Lng float64
}

// Maybe はオプショナル値を表すジェネリック型です。
// 値が存在するかどうかを示すOKフィールドと、実際の値を保持するVフィールドを持ちます。
// Goにおけるオプショナル型の実装として機能します。
type Maybe[T any] struct {
	V  T
	OK bool
}

// Just は値vを含むMaybe型を作成します。
// 作成されたMaybe型はOKフィールドがtrueに設定されます。
func Just[T any](v T) Maybe[T] {
	return Maybe[T]{V: v, OK: true}
}
