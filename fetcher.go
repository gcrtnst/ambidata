package ambidata

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Fetcher は Ambient からデータを取得するクライアントです。
type Fetcher struct {
	Ch      string // チャネルID
	ReadKey string // リードキー

	// Config は HTTP 通信の設定を保持します。
	// nil の場合は、デフォルトの設定が使用されます。
	Config *Config
}

// NewFetcher は新しい [Fetcher] を作成します。
func NewFetcher(ch string, readKey string) *Fetcher {
	return &Fetcher{Ch: ch, ReadKey: readKey}
}

// NewFetcherFromChannelAccess は [ChannelAccess] を基に新しい [Fetcher] を作成します。
func NewFetcherFromChannelAccess(ca *ChannelAccess) *Fetcher {
	return NewFetcher(ca.Ch, ca.ReadKey)
}

// GetChannel はチャネルの詳細情報を取得します。
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
// 最新から skip 件のデータを読み飛ばし、その先 n 件のデータを取得します。
// n と skip は非負の値である必要があります。
// データは新しいものから古いものの順に並びます。
//
// 推測に基づく情報: 取得できるデータ数は最大3000件です。
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
// データは新しいものから古いものの順に並びます。
//
// 推測に基づく情報: Ambient サーバーにおける時刻の精度はミリ秒単位のようです。
// より高精度な時刻を指定した場合、ミリ秒単位になるように切り捨てられます。
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
