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
