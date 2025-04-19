package ambidata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Config 構造体で値を設定しなかった場合に使用されるデフォルト値。
var (
	DefaultScheme = "https"       // [Config.Scheme] のデフォルト値
	DefaultHost   = "ambidata.io" // [Config.Host] のデフォルト値
)

// ErrRequestEntityTooLarge はリクエストボディが大きすぎる場合に返されるエラーです。
// 主に、 [Sender.SendBulk] メソッドに渡したデータが多すぎる場合に発生します。
var ErrRequestEntityTooLarge = errors.New("request entity too large")

// Config は HTTP 通信の設定を保持する構造体です。
//
// ゼロ値の Config 構造体は、デフォルトの設定を使用して Ambient に接続する有効な構成となります。
type Config struct {
	// Scheme はリクエストのスキームを指定します。
	// 空文字列の場合は、 [DefaultScheme] が使用されます。
	//
	// TLS を使用できない環境では、"http" を指定してください。
	Scheme string

	// Host はリクエストのホスト名を指定します。
	// 空文字列の場合は、 [DefaultHost] が使用されます。
	Host string

	// Client は HTTP リクエストを送信するためのクライアントを指定します。
	// nil の場合は、 [http.DefaultClient] が使用されます。
	Client *http.Client
}

// APIError は API リクエストに関連するエラーを表す構造体です。
//
// クエリパラメータのうち、"userKey"、"readKey"、"writeKey" はフィルタリングされ、
// Query フィールドから除外されます。
type APIError struct {
	Method string     // 例: "GET"
	Path   string     // 例: "/api/v2/channels/"
	Query  url.Values // 例: url.Values{"devKey": []string{"02:00:00:00:00:01"}}
	Err    error
}

func newAPIError(req *httpRequest, err error) error {
	return &APIError{
		Method: req.Method,
		Path:   req.Path,
		Query:  filterQuery(req.Query),
		Err:    err,
	}
}

func (err *APIError) Error() string {
	if err == nil || err.Method == "" || err.Path == "" || err.Err == nil {
		return fmt.Sprintf("%#v", err)
	}
	q := err.Query

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString("ambidata: ")
	b.WriteString(err.Method)
	b.WriteByte(' ')
	b.WriteString(err.Path)
	if len(q) > 0 {
		b.WriteByte('?')
		b.WriteString(err.Query.Encode())
	}
	b.WriteString(": ")
	b.WriteString(err.Err.Error())
	return b.String()
}

func (err *APIError) Unwrap() error {
	return err.Err
}

// StatusCodeError は HTTP ステータスコードに関連するエラーを表す構造体です。
// 200 OK 以外のステータスコードが返された場合に使用されます。
//
// 本パッケージから返される StatusCodeError は、全て [APIError] によってラップされています。
type StatusCodeError struct {
	StatusCode int
}

func (err *StatusCodeError) Error() string {
	if err == nil || err.StatusCode == 0 {
		return fmt.Sprintf("%#v", err)
	}

	code := err.StatusCode
	text := http.StatusText(code)
	if text == "" {
		text = "Unknown Status Code"
	}
	return strconv.Itoa(code) + " " + text
}

func httpGet(ctx context.Context, cfg *Config, path string, query url.Values, v any) error {
	var err error

	req := &httpRequest{
		Config: cfg,
		Method: "GET",
		Path:   path,
		Query:  query,
	}

	resp, err := httpDo(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // ignore error

	d := json.NewDecoder(resp.Body)
	err = d.Decode(v)
	if err != nil {
		return newAPIError(req, err)
	}

	return nil
}

func httpPost(ctx context.Context, cfg *Config, path string, v any) error {
	return httpSend(ctx, cfg, "POST", path, v)
}

func httpPut(ctx context.Context, cfg *Config, path string, v any) error {
	return httpSend(ctx, cfg, "PUT", path, v)
}

func httpDelete(ctx context.Context, cfg *Config, path string, query url.Values) error {
	var err error

	req := &httpRequest{
		Config: cfg,
		Method: "DELETE",
		Path:   path,
		Query:  query,
	}

	resp, err := httpDo(ctx, req)
	if err != nil {
		return err
	}
	return httpReadError(req, resp.Body)
}

func httpSend(ctx context.Context, cfg *Config, method string, path string, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	req := &httpRequest{
		Config:      cfg,
		Method:      method,
		Path:        path,
		Body:        body,
		ContentType: "application/json",
	}

	resp, err := httpDo(ctx, req)
	if err != nil {
		return err
	}
	return httpReadError(req, resp.Body)
}

type httpRequest struct {
	Config      *Config
	Method      string
	Path        string
	Query       url.Values
	Body        []byte
	ContentType string
}

func httpDo(ctx context.Context, req *httpRequest) (*http.Response, error) {
	cfg := valueOrDefault(req.Config, &Config{})
	scheme := valueOrDefault(cfg.Scheme, DefaultScheme)
	host := valueOrDefault(cfg.Host, DefaultHost)
	c := valueOrDefault(cfg.Client, http.DefaultClient)

	u := &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     req.Path,
		RawQuery: req.Query.Encode(),
	}

	var body io.ReadCloser
	var getBody func() (io.ReadCloser, error)
	if len(req.Body) > 0 {
		getBody = func() (io.ReadCloser, error) {
			body := io.NopCloser(bytes.NewReader(req.Body))
			return body, nil
		}
		body, _ = getBody()
	}

	header := make(http.Header)
	header.Set("User-Agent", "") // disable sending User-Agent
	if req.ContentType != "" {
		header.Set("Content-Type", req.ContentType)
	}

	hreq := &http.Request{
		Method:        req.Method,
		URL:           u,
		Header:        header,
		Body:          body,
		ContentLength: int64(len(req.Body)),
		Host:          host,
	}
	hreq = hreq.WithContext(ctx)

	resp, err := c.Do(hreq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()

		var err error
		err = &StatusCodeError{StatusCode: resp.StatusCode}
		err = newAPIError(req, err)
		return nil, err
	}
	return resp, nil
}

func httpReadError(req *httpRequest, body io.ReadCloser) error {
	defer body.Close() // ignore error

	var err error
	const str = "request entity too large"
	buf := make([]byte, len(str))

	_, err = io.ReadFull(body, buf)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}
		return newAPIError(req, err)
	}

	_, err = body.Read(buf[len(str):])
	if err != io.EOF {
		return nil
	}

	if string(buf) == str {
		return newAPIError(req, ErrRequestEntityTooLarge)
	}
	return nil
}

func filterQuery(query url.Values) url.Values {
	f := url.Values{}
	for k, v := range query {
		switch k {
		case "userKey":
		case "readKey":
		case "writeKey":
		default:
			f[k] = v
		}
	}
	return f
}

func valueOrDefault[T comparable](v, def T) T {
	var zero T
	if v == zero {
		v = def
	}
	return v
}
