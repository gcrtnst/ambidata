package ambidata

import (
	"context"
	"net/url"
)

// Manager はambidataのチャンネル管理機能を提供するクライアントです。
// ユーザーキーを使用して、チャンネルの一覧取得やデータの削除などの管理操作を行います。
type Manager struct {
	UserKey string
	Config  *Config
}

// NewManager は新しいManagerインスタンスを作成します。
// ユーザーキーを指定して、チャンネル管理用のクライアントを初期化します。
func NewManager(userKey string) *Manager {
	return &Manager{UserKey: userKey}
}

// GetChannelList はユーザーが所有するすべてのチャンネルのリストを取得します。
// 各チャンネルのアクセス情報（読み取りキーと書き込みキーを含む）を返します。
func (m *Manager) GetChannelList(ctx context.Context) ([]ChannelAccess, error) {
	var j jsonRecvChannelAccessList
	err := m.httpGet(ctx, "/api/v2/channels/", nil, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToChannelAccessList()
	return ret, nil
}

// GetDeviceChannel はデバイスキーに関連付けられたチャンネルのアクセス情報を取得します。
// チャンネル情報と読み取りキー、書き込みキーを含むChannelAccessオブジェクトを返します。
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

// GetDeviceChannelLv1 はデバイスキーに関連付けられたチャンネルのレベル1アクセス情報を取得します。
// チャンネルIDと書き込みキーのみを含む簡易的なアクセス情報を返します。
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

// DeleteData は指定されたチャンネルのすべてのデータを削除します。
// この操作は元に戻せないため、注意して使用してください。
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
