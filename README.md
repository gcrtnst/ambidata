# ambidata

[![Go Reference](https://pkg.go.dev/badge/github.com/gcrtnst/ambidata.svg)](https://pkg.go.dev/github.com/gcrtnst/ambidata)

ambidata は、IoTデータ可視化サービス [Ambient](https://ambidata.io/) の Go 言語向け非公式クライアントライブラリです。

## 特徴

- データの送信/取得/削除、チャネル情報の取得など、公式ライブラリと同等の機能を提供します。
- データポイントやチャネル情報は構造体として表現され、型安全に扱うことができます。
- Doc comments が日本語で詳細に記載されており、[pkg.go.dev](https://pkg.go.dev/github.com/gcrtnst/ambidata) やエディタ上でドキュメントを参照しながらコーディングできます。
- [The Unlicense](LICENSE) の下で公開されており、パブリックドメインとして自由に使用、改変、配布が可能です。

## インストール

```bash
go get github.com/gcrtnst/ambidata
```

## 基本的な使い方

ここでは、代表的な機能であるデータの送受信の例を示します。その他の機能については [pkg.go.dev](https://pkg.go.dev/github.com/gcrtnst/ambidata) を参照してください。

これらのサンプルコードを実行するには、以下の環境変数を事前に設定する必要があります。
- `AMBIDATA_CH`: チャネルID
- `AMBIDATA_WRITEKEY`: ライトキー
- `AMBIDATA_READKEY`: リードキー

### データの送信

```go
package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gcrtnst/ambidata"
)

func main() {
	// 設定値を取得
	ch := os.Getenv("AMBIDATA_CH")             // チャネルID
	writeKey := os.Getenv("AMBIDATA_WRITEKEY") // ライトキー

	// Sender を作成
	s := ambidata.NewSender(ch, writeKey)

	// 30秒ごとにデータを送信
	ctx := context.Background()
	for t := range time.Tick(30 * time.Second) {
		// 測定値を取得 (実際のセンサー値を取得する処理に置き換えてください)
		temp := dummyTemperature()
		humi := dummyHumidity()

		// 送信するデータを作成
		// データフィールドは Maybe[T] というオプショナル型となっており、
		// 値が無いフィールドは送信されません。
		data := ambidata.Data{
			Created: t,                   // データの作成日時
			D1:      ambidata.Just(temp), // 温度 (Just 関数は Maybe[T] 型の値を作成します)
			D2:      ambidata.Just(humi), // 湿度 (Just 関数は Maybe[T] 型の値を作成します)
		}

		// データを送信
		err := s.Send(ctx, data)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func dummyTemperature() float64 { return rand.NormFloat64()*10 + 25 }
func dummyHumidity() float64    { return rand.Float64() * 100 }
```

### データの取得

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gcrtnst/ambidata"
)

func main() {
	// 設定値を取得
	ch := os.Getenv("AMBIDATA_CH")           // チャネルID
	readKey := os.Getenv("AMBIDATA_READKEY") // リードキー

	// Fetcher を作成
	f := ambidata.NewFetcher(ch, readKey)

	// データを取得
	ctx := context.Background()
	arr, err := f.FetchRange(ctx, 10, 0) // 最新から0件スキップして10件取得
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, data := range arr {
		// データの生成時刻を出力
		fmt.Printf("created=%s", data.Created.Format(time.RFC3339Nano))

		// D1 フィールドが存在する場合は出力
		if data.D1.OK {
			fmt.Printf(", d1=%f", data.D1.V)
		}

		// D2 フィールドが存在する場合は出力
		if data.D2.OK {
			fmt.Printf(", d2=%f", data.D2.V)
		}

		fmt.Println()
	}
}
```

## 免責事項

- これは非公式ライブラリです。公式のサポートや保証はありません。サーバー側の仕様変更などにより、予告なく動作しなくなる可能性があります。
- [Ambient利用規約](https://ambidata.io/about/terms/) を遵守してください。
- データの送信間隔は、[諸元/制限](https://ambidata.io/refs/spec/) に基づきユーザー側で制限してください。本ライブラリでは送信間隔を制限しておりません。

## ライセンス

このライブラリは The Unlicense の下で公開されています。詳細は [LICENSE](LICENSE) を参照してください。
