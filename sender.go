package ambidata

import (
	"context"
	"net/url"
	"time"
)

// Sender は Ambient にデータを送信するクライアントです。
type Sender struct {
	Ch       string // チャネルID
	WriteKey string // ライトキー

	// Config は HTTP 通信の設定を保持します。
	// nil の場合は、デフォルトの設定が使用されます。
	Config *Config
}

// NewSender は新しい [Sender] を作成します。
func NewSender(ch string, writeKey string) *Sender {
	return &Sender{Ch: ch, WriteKey: writeKey}
}

// NewSenderFromChannelAccess は [ChannelAccess] を基に新しい [Sender] を作成します。
func NewSenderFromChannelAccess(ca *ChannelAccess) *Sender {
	return NewSender(ca.Ch, ca.WriteKey)
}

// NewSenderFromChannelAccessLv1 は [ChannelAccessLv1] を基に新しい [Sender] を作成します。
func NewSenderFromChannelAccessLv1(ca1 *ChannelAccessLv1) *Sender {
	return NewSender(ca1.Ch, ca1.WriteKey)
}

// Send は単一のデータポイントをチャネルに送信します。
//
// [Data.Created] にゼロ値を指定した場合、サーバー側で現在時刻が設定されます。
//
// 推測に基づく情報: [Data.Cmnt] の最大長は 64 バイトのようです。
// 最大長を超えた部分はサーバーによって切り捨てられます。
//
// Send は [Data.Hide] フィールドを送信しません。
// [Data.Hide] フィールドを設定するには、データポイントを送信した後に、
// [Sender.SetHide] メソッドを使用してください。
func (s *Sender) Send(ctx context.Context, data Data) error {
	j := jsonSendDataRequest{
		jsonSendData: toJSONSendData(data),
		WriteKey:     s.WriteKey,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPost(ctx, s.Config, path, j)
}

// SendBulk は複数のデータポイントを一括でチャネルに送信します。
//
// [Data.Created] にゼロ値を指定した場合、サーバー側で現在時刻が設定されます。
//
// 推測に基づく情報: [Data.Cmnt] の最大長は 64 バイトのようです。
// 最大長を超えた部分はサーバーによって切り捨てられます。
//
// SendBulk は [Data.Hide] フィールドを送信しません。
// [Data.Hide] フィールドを設定するには、データポイントを送信した後に、
// [Sender.SetHide] メソッドを使用してください。
//
// 送信するデータポイントが多すぎる場合、リクエストボディのサイズ制限を超過して、
// [ErrRequestEntityTooLarge] エラーが発生することがあります。
// 推測に基づく情報: 執筆時点では、データポイントが 258 個以内であれば、
// サイズ制限に達しないようです。ただし、この個数は今後サーバー/クライアント双方の
// 更新によって変化する可能性があります。そのため、一度に送信する個数は、
// サイズ制限から余裕を持たせて、200 個以内にすることをお勧めします。
func (s *Sender) SendBulk(ctx context.Context, arr []Data) error {
	// Doc comment についてのメモ:
	// データポイントの最大数は ambitest.TestSenderSendBulkTooLarge でテストしています。

	j := jsonSendDataListRequest{
		WriteKey: s.WriteKey,
		Data:     toJSONSendDataList(arr),
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/dataarray"
	return httpPost(ctx, s.Config, path, j)
}

// SetCmnt は指定された時刻のデータポイントにコメントを設定します。
//
// 推測に基づく情報: 指定された時刻に該当するデータポイントが複数存在する場合、
// そのうちのどれか1つにコメントが設定されます。
func (s *Sender) SetCmnt(ctx context.Context, created time.Time, cmnt string) error {
	j := jsonSendCmnt{
		WriteKey: s.WriteKey,
		Created:  created,
		Cmnt:     cmnt,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPut(ctx, s.Config, path, j)
}

// SetHide は指定された時刻のデータポイントの表示/非表示状態を設定します。
//
// 推測に基づく情報: 指定された時刻に該当するデータポイントが複数存在する場合、
// そのうちのどれか1つに表示/非表示状態が設定されます。
func (s *Sender) SetHide(ctx context.Context, created time.Time, hide bool) error {
	j := jsonSendHide{
		WriteKey: s.WriteKey,
		Created:  created,
		Hide:     hide,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPut(ctx, s.Config, path, j)
}
