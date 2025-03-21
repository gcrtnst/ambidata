package ambidata

import (
	"context"
	"net/url"
)

type Fetcher struct {
	ReadKey string
	Ch      string
	Config  *Config
}

func (f *Fetcher) GetChannel(ctx context.Context) (ChannelInfo, error) {
	path := "/api/v2/channels/" + url.PathEscape(f.Ch) + "/"
	var j jsonRecvChannelInfo
	err := f.httpGetJSON(ctx, path, nil, &j)
	if err != nil {
		return ChannelInfo{}, err
	}

	ret := j.ToChannelInfo()
	return ret, nil
}

func (f *Fetcher) httpGetJSON(ctx context.Context, path string, query url.Values, v any) error {
	const key = "readKey"
	val := f.ReadKey
	if query == nil {
		query = url.Values{key: []string{val}}
	} else {
		query.Set(key, val)
	}

	return httpGetJSON(ctx, f.Config, path, query, v)
}
