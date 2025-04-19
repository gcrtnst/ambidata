package ambidata

import (
	"context"
	"net/url"
)

// Manager は Ambient のチャネル管理機能を提供するクライアントです。
type Manager struct {
	UserKey string // ユーザーキー

	// Config は HTTP 通信の設定を保持します。
	// nil の場合は、デフォルトの設定が使用されます。
	Config *Config
}

// NewManager は新しい [Manager] を作成します。
func NewManager(userKey string) *Manager {
	return &Manager{UserKey: userKey}
}

// GetChannelList はユーザーが所有するすべてのチャネルのリストを取得します。
func (m *Manager) GetChannelList(ctx context.Context) ([]ChannelAccess, error) {
	var j jsonRecvChannelAccessList
	err := m.httpGet(ctx, "/api/v2/channels/", nil, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToChannelAccessList()
	return ret, nil
}

// GetDeviceChannel は指定されたデバイスキーに関連付けられたチャネルの情報を取得します。
func (m *Manager) GetDeviceChannel(ctx context.Context, devKey string) (ChannelAccess, error) {
	var j jsonRecvChannelAccess
	query := url.Values{"devKey": []string{devKey}}
	err := m.httpGet(ctx, "/api/v2/channels/", query, &j)
	if err != nil {
		return ChannelAccess{}, err
	}

	ret := j.ToChannelAccess()
	return ret, nil
}

// GetDeviceChannelLv1 は指定されたデバイスキーに関連付けられたチャネルの情報のうち、ID とライトキーのみを取得します。
//
// Lv1 という名前は、API に対するクエリパラメータ level=1 に由来します。
func (m *Manager) GetDeviceChannelLv1(ctx context.Context, devKey string) (ChannelAccessLv1, error) {
	var j jsonRecvChannelAccessLv1
	query := url.Values{"devKey": []string{devKey}, "level": []string{"1"}}
	err := m.httpGet(ctx, "/api/v2/channels/", query, &j)
	if err != nil {
		return ChannelAccessLv1{}, err
	}

	ret := j.ToChannelAccessLv1()
	return ret, nil
}

// DeleteData は指定されたチャネルのすべてのデータを削除します。
// 部分的に削除する機能はありません。
// Ambient サイトのチャネルページの「データー削除」ボタンと同じ動作をします。
// チャネルそのものやチャネルの設定情報は削除されずに残ります。
// 削除したデーターは復元できないので、注意してください。
func (m *Manager) DeleteData(ctx context.Context, ch string) error {
	path := "/api/v2/channels/" + url.PathEscape(ch) + "/data"
	return m.httpDelete(ctx, path, nil)
}

func (m *Manager) httpGet(ctx context.Context, path string, query url.Values, v any) error {
	query = m.ensureUserKey(query)
	return httpGet(ctx, m.Config, path, query, v)
}

func (m *Manager) httpDelete(ctx context.Context, path string, query url.Values) error {
	query = m.ensureUserKey(query)
	return httpDelete(ctx, m.Config, path, query)
}

func (m *Manager) ensureUserKey(query url.Values) url.Values {
	const key = "userKey"
	val := m.UserKey
	if query == nil {
		query = url.Values{key: []string{val}}
	} else {
		query.Set(key, val)
	}
	return query
}
