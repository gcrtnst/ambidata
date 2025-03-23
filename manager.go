package ambidata

import (
	"context"
	"net/url"
)

type Manager struct {
	UserKey string
	Config  *Config
}

func NewManager(userKey string) *Manager {
	return &Manager{UserKey: userKey}
}

func (m *Manager) GetChannelList(ctx context.Context) ([]ChannelAccess, error) {
	var j jsonRecvChannelAccessList
	err := m.httpGetJSON(ctx, "/api/v2/channels/", nil, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToChannelAccessList()
	return ret, nil
}

func (m *Manager) GetDeviceChannel(ctx context.Context, devKey string) (ChannelAccess, error) {
	var j jsonRecvChannelAccess
	query := url.Values{"devKey": []string{devKey}}
	err := m.httpGetJSON(ctx, "/api/v2/channels/", query, &j)
	if err != nil {
		return ChannelAccess{}, err
	}

	ret := j.ToChannelAccess()
	return ret, nil
}

func (m *Manager) GetDeviceChannelLv1(ctx context.Context, devKey string) (ChannelAccessLv1, error) {
	var j jsonRecvChannelAccessLv1
	query := url.Values{"devKey": []string{devKey}, "level": []string{"1"}}
	err := m.httpGetJSON(ctx, "/api/v2/channels/", query, &j)
	if err != nil {
		return ChannelAccessLv1{}, err
	}

	ret := j.ToChannelAccessLv1()
	return ret, nil
}

func (m *Manager) DeleteData(ctx context.Context, ch string) error {
	path := "/api/v2/channels/" + url.PathEscape(ch) + "/data"
	return m.httpDelete(ctx, path, nil)
}

func (m *Manager) httpGetJSON(ctx context.Context, path string, query url.Values, v any) error {
	query = m.ensureUserKey(query)
	return httpGetJSON(ctx, m.Config, path, query, v)
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
