package ambidata

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Fetcher はambidataからデータを取得するためのクライアントです。
// チャンネルIDと読み取りキーを使用して、特定のチャンネルからデータを取得します。
type Fetcher struct {
	Ch      string
	ReadKey string
	Config  *Config
}

// NewFetcher は新しいFetcherインスタンスを作成します。
// チャンネルIDと読み取りキーを指定して、データ取得用のクライアントを初期化します。
func NewFetcher(ch string, readKey string) *Fetcher {
	return &Fetcher{Ch: ch, ReadKey: readKey}
}

// NewFetcherFromChannelAccess はChannelAccessオブジェクトから新しいFetcherインスタンスを作成します。
// ChannelAccessに含まれるチャンネルIDと読み取りキーを使用してFetcherを初期化します。
func NewFetcherFromChannelAccess(ca *ChannelAccess) *Fetcher {
	return NewFetcher(ca.Ch, ca.ReadKey)
}

// GetChannel はチャンネルの詳細情報を取得します。
// 指定されたコンテキストを使用してAPIリクエストを実行し、チャンネル情報を返します。
func (f *Fetcher) GetChannel(ctx context.Context) (ChannelInfo, error) {
	path := "/api/v2/channels/" + url.PathEscape(f.Ch) + "/"
	var j jsonRecvChannelInfo
	err := f.httpGet(ctx, path, nil, &j)
	if err != nil {
		return ChannelInfo{}, err
	}

	ret := j.ToChannelInfo()
	return ret, nil
}

// FetchRange は指定された範囲のデータを取得します。
// n個のデータを取得し、skip個のデータをスキップします。
// nとskipは非負の値である必要があります。
func (f *Fetcher) FetchRange(ctx context.Context, n int, skip int) ([]Data, error) {
	if n < 0 || skip < 0 {
		err := fmt.Errorf("ambidata: (*Fetcher).FetchRange: n and skip must be non-negative (n=%d, skip=%d)", n, skip)
		return nil, err
	}
	if n <= 0 {
		return []Data{}, nil
	}

	query := url.Values{"n": []string{strconv.Itoa(n)}}
	if skip > 0 {
		query.Set("skip", strconv.Itoa(skip))
	}

	path := "/api/v2/channels/" + url.PathEscape(f.Ch) + "/data"
	var j jsonRecvDataList
	err := f.httpGet(ctx, path, query, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToDataList()
	return ret, nil
}

// FetchPeriod は指定された期間のデータを取得します。
// 開始時刻から終了時刻までの間に作成されたデータを返します。
// 開始時刻が終了時刻より後の場合は空のスライスを返します。
func (f *Fetcher) FetchPeriod(ctx context.Context, start time.Time, end time.Time) ([]Data, error) {
	if !start.Before(end) {
		return []Data{}, nil
	}

	path := "/api/v2/channels/" + url.PathEscape(f.Ch) + "/data"
	query := url.Values{
		"start": []string{start.Format(time.RFC3339Nano)},
		"end":   []string{end.Format(time.RFC3339Nano)},
	}

	var j jsonRecvDataList
	err := f.httpGet(ctx, path, query, &j)
	if err != nil {
		return nil, err
	}

	ret := j.ToDataList()
	return ret, nil
}

func (f *Fetcher) httpGet(ctx context.Context, path string, query url.Values, v any) error {
	const key = "readKey"
	val := f.ReadKey
	if query == nil {
		query = url.Values{key: []string{val}}
	} else {
		query.Set(key, val)
	}

	return httpGet(ctx, f.Config, path, query, v)
}
