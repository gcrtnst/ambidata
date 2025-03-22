package ambidata

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type Fetcher struct {
	Ch      string
	ReadKey string
	Config  *Config
}

func NewFetcher(ch string, readKey string) *Fetcher {
	return &Fetcher{Ch: ch, ReadKey: readKey}
}

func NewFetcherFromChannelAccess(ca *ChannelAccess) *Fetcher {
	return NewFetcher(ca.Ch, ca.ReadKey)
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

func (f *Fetcher) FetchRange(ctx context.Context, n int, skip int) ([]Data, error) {
	if n < 0 || skip < 0 {
		err := fmt.Errorf("ambidata: (*Fetcher).FetchRange: n and skip must be non-negative (n=%d, skip=%d)", n, skip)
		return nil, err
	}
	if n-skip <= 0 {
		return []Data{}, nil
	}

	query := url.Values{"n": []string{strconv.Itoa(n)}}
	if skip > 0 {
		query.Set("skip", strconv.Itoa(skip))
	}

	path := "/api/v2/channels/" + url.PathEscape(f.Ch) + "/data"
	var j jsonRecvDataList
	err := f.httpGetJSON(ctx, path, query, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToDataList()
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
