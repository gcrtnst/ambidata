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
// [Data.Created] フィールドにゼロ値を指定した場合、サーバー側で現在時刻が設定されます。
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
// [Data.Created] フィールドにゼロ値を指定した場合、サーバー側で現在時刻が設定されます。
//
// SendBulk は [Data.Hide] フィールドを送信しません。
// [Data.Hide] フィールドを設定するには、データポイントを送信した後に、
// [Sender.SetHide] メソッドを使用してください。
//
// 送信するデータポイントが多すぎる場合、リクエストボディのサイズ制限を超過して、
// [ErrRequestEntityTooLarge] エラーが発生することがあります。
func (s *Sender) SendBulk(ctx context.Context, arr []Data) error {
	j := jsonSendDataListRequest{
		WriteKey: s.WriteKey,
		Data:     toJSONSendDataList(arr),
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/dataarray"
	return httpPost(ctx, s.Config, path, j)
}

// SetCmnt は指定された時刻のデータポイントにコメントを設定します。
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
func (s *Sender) SetHide(ctx context.Context, created time.Time, hide bool) error {
	j := jsonSendHide{
		WriteKey: s.WriteKey,
		Created:  created,
		Hide:     hide,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPut(ctx, s.Config, path, j)
}
