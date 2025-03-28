package ambidata

import (
	"context"
	"net/url"
	"time"
)

type Sender struct {
	Ch       string
	WriteKey string
	Config   *Config
}

func NewSender(ch string, writeKey string) *Sender {
	return &Sender{Ch: ch, WriteKey: writeKey}
}

func NewSenderFromChannelAccess(ca *ChannelAccess) *Sender {
	return NewSender(ca.Ch, ca.WriteKey)
}

func NewSenderFromChannelAccessLv1(ca1 *ChannelAccessLv1) *Sender {
	return NewSender(ca1.Ch, ca1.WriteKey)
}

func (s *Sender) Send(ctx context.Context, data Data) error {
	j := jsonSendDataRequest{
		jsonSendData: toJSONSendData(data),
		WriteKey:     s.WriteKey,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPost(ctx, s.Config, path, j)
}

func (s *Sender) SendBulk(ctx context.Context, arr []Data) error {
	j := jsonSendDataListRequest{
		WriteKey: s.WriteKey,
		Data:     toJSONSendDataList(arr),
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/dataarray"
	return httpPost(ctx, s.Config, path, j)
}

func (s *Sender) SetCmnt(ctx context.Context, created time.Time, cmnt string) error {
	j := jsonSendCmnt{
		WriteKey: s.WriteKey,
		Created:  created,
		Cmnt:     cmnt,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPut(ctx, s.Config, path, j)
}

func (s *Sender) SetHide(ctx context.Context, created time.Time, hide bool) error {
	j := jsonSendHide{
		WriteKey: s.WriteKey,
		Created:  created,
		Hide:     hide,
	}

	path := "/api/v2/channels/" + url.PathEscape(s.Ch) + "/data"
	return httpPut(ctx, s.Config, path, j)
}
