package ambidata

import (
	"context"
	"net/url"
	"time"
)

// Sender はambidataにデータを送信するためのクライアントです。
// チャネルIDと書き込みキーを使用して、特定のチャネルにデータを送信します。
type Sender struct {
	Ch       string
	WriteKey string
	Config   *Config
}

// NewSender は新しいSenderインスタンスを作成します。
// チャネルIDと書き込みキーを指定して、データ送信用のクライアントを初期化します。
func NewSender(ch string, writeKey string) *Sender {
	return &Sender{Ch: ch, WriteKey: writeKey}
}

// NewSenderFromChannelAccess はChannelAccessオブジェクトから新しいSenderインスタンスを作成します。
// ChannelAccessに含まれるチャネルIDと書き込みキーを使用してSenderを初期化します。
func NewSenderFromChannelAccess(ca *ChannelAccess) *Sender {
	return NewSender(ca.Ch, ca.WriteKey)
}

// NewSenderFromChannelAccessLv1 はChannelAccessLv1オブジェクトから新しいSenderインスタンスを作成します。
// ChannelAccessLv1に含まれるチャネルIDと書き込みキーを使用してSenderを初期化します。
func NewSenderFromChannelAccessLv1(ca1 *ChannelAccessLv1) *Sender {
	return NewSender(ca1.Ch, ca1.WriteKey)
}

// Send は単一のデータポイントをチャネルに送信します。
// 指定されたデータオブジェクトをJSON形式に変換してAPIに送信します。
func (s *Sender) Send(ctx context.Context, data Data) error {
	j := jsonSendDataRequest{
		jsonSendData: toJSONSendData(data),
		WriteKey:     s.WriteKey,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPost(ctx, s.Config, path, j)
}

// SendBulk は複数のデータポイントを一括でチャネルに送信します。
// データの配列をJSON形式に変換して一度のAPIリクエストで送信します。
// 大量のデータを効率的に送信する場合に使用します。
func (s *Sender) SendBulk(ctx context.Context, arr []Data) error {
	j := jsonSendDataListRequest{
		WriteKey: s.WriteKey,
		Data:     toJSONSendDataList(arr),
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/dataarray"
	return httpPost(ctx, s.Config, path, j)
}

// SetCmnt は指定された時刻のデータポイントにコメントを設定します。
// 既存のデータポイントのコメントフィールドを更新します。
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
// hideがtrueの場合、データポイントはグラフやリストに表示されなくなります。
func (s *Sender) SetHide(ctx context.Context, created time.Time, hide bool) error {
	j := jsonSendHide{
		WriteKey: s.WriteKey,
		Created:  created,
		Hide:     hide,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPut(ctx, s.Config, path, j)
}
