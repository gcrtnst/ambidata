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
