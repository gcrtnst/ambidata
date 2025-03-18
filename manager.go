package ambidata

import (
	"context"
	"net/url"
)

type Manager struct {
	UserKey string
	Config  *Config
}

func (m *Manager) GetChannelList(ctx context.Context) ([]ChannelAccess, error) {
	var j jsonRecvChannelAccessList
	err := m.httpGetJSON(ctx, PathGetChannelList, nil, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToChannelAccessList()
	return ret, nil
}

func (m *Manager) httpGetJSON(ctx context.Context, path string, query url.Values, v any) error {
	const key = "userKey"
	val := m.UserKey
	if query == nil {
		query = url.Values{key: []string{val}}
	} else {
		query.Set(key, val)
	}

	return httpGetJSON(ctx, m.Config, path, query, v)
}
