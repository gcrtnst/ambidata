package ambidata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	DefaultScheme = "https"
	DefaultHost   = "ambidata.io"
)

var defaultHeader = http.Header{
	"User-Agent": nil, // disable sending User-Agent
}

type Config struct {
	Scheme string
	Host   string
	Client *http.Client
}

type APIError struct {
	Method string
	Path   string
	Query  url.Values
	Err    error
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

func httpGetJSON(ctx context.Context, cfg *Config, path string, query url.Values, v any) error {
	const method = "GET"
	var err error

	req := &httpRequest{
		Config: cfg,
		Method: method,
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
		return &APIError{
			Method: method,
			Path:   path,
			Query:  filterQuery(query),
			Err:    err,
		}
	}

	return nil
}

func httpDelete(ctx context.Context, cfg *Config, path string, query url.Values) (*http.Response, error) {
	return httpDo(ctx, &httpRequest{
		Config: cfg,
		Method: "DELETE",
		Path:   path,
		Query:  query,
	})
}

type httpRequest struct {
	Config *Config
	Method string
	Path   string
	Query  url.Values
	Body   []byte
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

	hreq := &http.Request{
		Method:        req.Method,
		URL:           u,
		Header:        defaultHeader,
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
		return nil, &APIError{
			Method: hreq.Method,
			Path:   hreq.URL.Path,
			Query:  filterQuery(req.Query),
			Err:    &StatusCodeError{StatusCode: resp.StatusCode},
		}
	}
	return resp, nil
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
