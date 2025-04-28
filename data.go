package ambidata

import (
	"image/color"
	"time"
)

// ChannelAccess はチャネルへのアクセス情報を保持する構造体です。
type ChannelAccess struct {
	ChannelInfo
	ReadKey  string // リードキー
	WriteKey string // ライトキー
}

// ToLv1 は [ChannelAccess] から [ChannelAccessLv1] への変換を行います。
func (ch *ChannelAccess) ToLv1() ChannelAccessLv1 {
	return ChannelAccessLv1{
		Ch:       ch.Ch,
		WriteKey: ch.WriteKey,
	}
}

// ChannelAccessLv1 はチャネルへのアクセス情報のうち、ID とライトキーのみを保持する構造体です。
type ChannelAccessLv1 struct {
	Ch       string
	WriteKey string
}

// ChannelInfo はチャネルの詳細情報を保持する構造体です。
//
// 推測に基づく情報: チャネル情報の形式は公式ライブラリやリファレンスでは明示されていないため、
// この構造体のフィールドの多くは本パッケージ開発者の推測に基づいています。
type ChannelInfo struct {
	Ch         string          // チャネルID
	User       string          // ユーザーID
	Created    time.Time       // チャネル作成日時
	Modified   time.Time       // チャネル更新日時
	LastPost   time.Time       // データの最終送信日時(一度も送信されていない場合はゼロ値)
	Charts     int             // 不明
	DataPerDay int             // データの一日あたりの平均数
	DCh        bool            // 不明
	ChName     string          // チャネル名
	ChDesc     string          // チャネルの説明
	D1         FieldInfo       // データ1
	D2         FieldInfo       // データ2
	D3         FieldInfo       // データ3
	D4         FieldInfo       // データ4
	D5         FieldInfo       // データ5
	D6         FieldInfo       // データ6
	D7         FieldInfo       // データ7
	D8         FieldInfo       // データ8
	Loc        Maybe[Location] // 位置情報
	PhotoID    string          // 写真の Embed code
	DevKeys    []string        // デバイスキー
	Bd         string          // ボードID
	LastData   LastData        // 最後に送信されたデータ
}

// FieldInfo はデータフィールドの情報を保持する構造体です。
type FieldInfo struct {
	Name  string     // データ名
	Color FieldColor // 色ID
}

// FieldColor はデータフィールドの色IDを表す型です。
// [FieldColor.ToRGBA] 関数で  [color.RGBA] 値に変換できます。
type FieldColor string

// データフィールドに使用できる色IDの定義。
// これらの色は Ambient のグラフや UI で使用されます。
const (
	FieldColorBlue    FieldColor = "1"  // #3B5998
	FieldColorRed     FieldColor = "2"  // #DC3912
	FieldColorOrange  FieldColor = "3"  // #FF9900
	FieldColorPurple  FieldColor = "4"  // #990099
	FieldColorGreen   FieldColor = "5"  // #109618
	FieldColorSkyBlue FieldColor = "6"  // #0099C6
	FieldColorPink    FieldColor = "7"  // #DD4477
	FieldColorBrown   FieldColor = "8"  // #996633
	FieldColorOlive   FieldColor = "9"  // #66AA00
	FieldColorCyan    FieldColor = "10" // #00FFFF
	FieldColorYellow  FieldColor = "11" // #FFFF00
	FieldColorBlack   FieldColor = "12" // #000000
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

// ToRGBA は指定された [FieldColor] を [color.RGBA] 型の値に変換します。
// 有効な色IDの場合、対応する [color.RGBA] 値を返し、ok は true に設定されます。
// 無効な色IDの場合、ゼロ値の [color.RGBA] を返し、ok は false に設定されます。
func (c FieldColor) ToRGBA() (rgba color.RGBA, ok bool) {
	rgba, ok = colorMap[c]
	return
}

// LastData はチャネルに最後に送信されたデータを表す構造体です。
//
// チャネルに一度もデータが送信されていない場合、ID は空文字列になります。
type LastData struct {
	Data
	ID string
}

// Data はチャネルに保存されるデータポイントを表す構造体です。
type Data struct {
	Created time.Time       // データの生成時刻
	D1      Maybe[float64]  // データ1
	D2      Maybe[float64]  // データ2
	D3      Maybe[float64]  // データ3
	D4      Maybe[float64]  // データ4
	D5      Maybe[float64]  // データ5
	D6      Maybe[float64]  // データ6
	D7      Maybe[float64]  // データ7
	D8      Maybe[float64]  // データ8
	Loc     Maybe[Location] // 位置情報
	Cmnt    string          // コメント
	Hide    bool            // 非表示フラグ
}

// Location は位置情報を表す構造体です。
type Location struct {
	Lat float64 // 緯度
	Lng float64 // 経度
}

// Maybe はオプショナル値を表すジェネリック型です。
type Maybe[T any] struct {
	V  T    // 値
	OK bool // 値が存在する場合は true
}

// Just は値 v を含む [Maybe] 型を作成します。
// 作成された [Maybe] 型は OK フィールドが true に設定されます。
func Just[T any](v T) Maybe[T] {
	return Maybe[T]{V: v, OK: true}
}
