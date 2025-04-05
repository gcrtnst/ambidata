package ambidata

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSenderSendNormal(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = ""
	const wantCT = "application/json"

	inData := Data{
		Created: time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC),
		D1:      Just(101.0),
		D2:      Just(102.0),
		D3:      Just(103.0),
		D4:      Just(104.0),
		D5:      Just(105.0),
		D6:      Just(106.0),
		D7:      Just(107.0),
		D8:      Just(108.0),
		Loc:     Just(Location{Lat: 109, Lng: 110}),
		Cmnt:    "111",
		Hide:    true,
	}

	wantJSON := map[string]any{
		"writeKey": inWriteKey,
		"created":  "1970-01-01T01:00:00Z",
		"d1":       101.0,
		"d2":       102.0,
		"d3":       103.0,
		"d4":       104.0,
		"d5":       105.0,
		"d6":       106.0,
		"d7":       107.0,
		"d8":       108.0,
		"lat":      109.0,
		"lng":      110.0,
		"cmnt":     "111",
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	var gotReqBody []byte
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("POST /api/v2/channels/"+inCh+"/data", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		gotReqBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.Send(ctx, inData)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}
	if gotCT := gotReq.Header.Get("Content-Type"); gotCT != wantCT {
		t.Errorf("request: Content-Type: expected %#v, got %#v", wantCT, gotCT)
	}

	var gotJSON map[string]any
	if err := json.Unmarshal(gotReqBody, &gotJSON); err != nil {
		t.Errorf("response: body: %v", err)
	} else if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("response: body: mismatch (-want, +got)\n%s", diff)
	}
}

func TestSenderSendEmpty(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = ""
	const wantCT = "application/json"

	inData := Data{}
	wantJSON := map[string]any{"writeKey": inWriteKey}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	var gotReqBody []byte
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("POST /api/v2/channels/"+inCh+"/data", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		gotReqBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.Send(ctx, inData)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}
	if gotCT := gotReq.Header.Get("Content-Type"); gotCT != wantCT {
		t.Errorf("request: Content-Type: expected %#v, got %#v", wantCT, gotCT)
	}

	var gotJSON map[string]any
	if err := json.Unmarshal(gotReqBody, &gotJSON); err != nil {
		t.Errorf("response: body: %v", err)
	} else if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("response: body: mismatch (-want, +got)\n%s", diff)
	}
}

func TestSenderSendErrNaN(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = ""
	inData := Data{D1: Just(math.NaN())}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq bool
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("POST /api/v2/channels/"+inCh+"/data", func(w http.ResponseWriter, r *http.Request) {
		gotReq = true
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.Send(ctx, inData)
	if gotJSONErr := (&json.UnsupportedValueError{}); !errors.As(gotErr, &gotJSONErr) {
		t.Errorf("err: expected %T, got %T", gotJSONErr, gotErr)
	}
	if gotReq {
		t.Errorf("request: unexpected HTTP request received")
	}
}

func TestSenderSendErrCanceled(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inData := Data{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.Send(ctx, inData)
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestSenderSendErrStatus(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inData := Data{}
	const wantMethod = "POST"
	const wantPath = "/api/v2/channels/83601/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.Send(ctx, inData)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotStatusErr, ok := gotAPIErr.Err.(*StatusCodeError); !ok {
			t.Errorf("err.Err: expected (*ambidata.StatusCodeError), got %T", gotAPIErr.Err)
		} else if gotStatusErr.StatusCode != http.StatusNotFound {
			t.Errorf("err.StatusCode: expected %d, got %d", http.StatusNotFound, gotStatusErr.StatusCode)
		}
	}
}

func TestSenderSendErrRequestEntityTooLarge(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = "request entity too large"
	inData := Data{}
	const wantMethod = "POST"
	const wantPath = "/api/v2/channels/83601/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.Send(ctx, inData)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotAPIErr.Err != ErrRequestEntityTooLarge {
			t.Errorf("err:Err: expected ErrRequestEntityTooLarge, got %q", gotAPIErr.Err.Error())
		}
	}
}

func TestSenderSendBulkNormal(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = ""
	const wantCT = "application/json"

	inArr := []Data{
		{},
		{
			Created: time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC),
			D1:      Just(101.0),
			D2:      Just(102.0),
			D3:      Just(103.0),
			D4:      Just(104.0),
			D5:      Just(105.0),
			D6:      Just(106.0),
			D7:      Just(107.0),
			D8:      Just(108.0),
			Loc:     Just(Location{Lat: 109, Lng: 110}),
			Cmnt:    "111",
			Hide:    true,
		},
	}

	wantJSON := map[string]any{
		"writeKey": inWriteKey,
		"data": []any{
			map[string]any{},
			map[string]any{
				"created": "1970-01-01T01:00:00Z",
				"d1":      101.0,
				"d2":      102.0,
				"d3":      103.0,
				"d4":      104.0,
				"d5":      105.0,
				"d6":      106.0,
				"d7":      107.0,
				"d8":      108.0,
				"lat":     109.0,
				"lng":     110.0,
				"cmnt":    "111",
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	var gotReqBody []byte
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("POST /api/v2/channels/"+inCh+"/dataarray", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		gotReqBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SendBulk(ctx, inArr)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}
	if gotCT := gotReq.Header.Get("Content-Type"); gotCT != wantCT {
		t.Errorf("request: Content-Type: expected %#v, got %#v", wantCT, gotCT)
	}

	var gotJSON map[string]any
	if err := json.Unmarshal(gotReqBody, &gotJSON); err != nil {
		t.Errorf("response: body: %v", err)
	} else if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("response: body: mismatch (-want, +got)\n%s", diff)
	}
}

func TestSenderSendBulkErrNaN(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = ""
	inArr := []Data{{D1: Just(math.NaN())}}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq bool
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("POST /api/v2/channels/"+inCh+"/dataarray", func(w http.ResponseWriter, r *http.Request) {
		gotReq = true
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SendBulk(ctx, inArr)
	if gotJSONErr := (&json.UnsupportedValueError{}); !errors.As(gotErr, &gotJSONErr) {
		t.Errorf("err: expected %T, got %T", gotJSONErr, gotErr)
	}
	if gotReq {
		t.Errorf("request: unexpected HTTP request received")
	}
}

func TestSenderSendBulkErrCanceled(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inArr := []Data{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SendBulk(ctx, inArr)
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestSenderSendBulkErrStatus(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inArr := []Data{}
	const wantMethod = "POST"
	const wantPath = "/api/v2/channels/83601/dataarray"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SendBulk(ctx, inArr)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotStatusErr, ok := gotAPIErr.Err.(*StatusCodeError); !ok {
			t.Errorf("err.Err: expected (*ambidata.StatusCodeError), got %T", gotAPIErr.Err)
		} else if gotStatusErr.StatusCode != http.StatusNotFound {
			t.Errorf("err.StatusCode: expected %d, got %d", http.StatusNotFound, gotStatusErr.StatusCode)
		}
	}
}

func TestSenderSendBulkErrRequestEntityTooLarge(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	const inCode = 200
	const inBody = "request entity too large"
	inArr := []Data{}
	const wantMethod = "POST"
	const wantPath = "/api/v2/channels/83601/dataarray"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SendBulk(ctx, inArr)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotAPIErr.Err != ErrRequestEntityTooLarge {
			t.Errorf("err:Err: expected ErrRequestEntityTooLarge, got %q", gotAPIErr.Err.Error())
		}
	}
}

func TestSenderSetCmntNormal(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inCmnt = "comment"
	const inCode = 200
	const inBody = "OK"
	const wantCT = "application/json"

	wantJSON := map[string]any{
		"writeKey": inWriteKey,
		"created":  "1970-01-01T01:00:00Z",
		"cmnt":     "comment",
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	var gotReqBody []byte
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("PUT /api/v2/channels/"+inCh+"/data", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		gotReqBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetCmnt(ctx, inCreated, inCmnt)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}
	if gotCT := gotReq.Header.Get("Content-Type"); gotCT != wantCT {
		t.Errorf("request: Content-Type: expected %#v, got %#v", wantCT, gotCT)
	}

	var gotJSON map[string]any
	if err := json.Unmarshal(gotReqBody, &gotJSON); err != nil {
		t.Errorf("response: body: %v", err)
	} else if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("response: body: mismatch (-want, +got)\n%s", diff)
	}
}

func TestSenderSetCmntErrCanceled(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inCmnt = "comment"
	const inCode = 200
	const inBody = "OK"

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetCmnt(ctx, inCreated, inCmnt)
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestSenderSetCmntErrStatus(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inCmnt = "comment"
	const wantMethod = "PUT"
	const wantPath = "/api/v2/channels/83601/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetCmnt(ctx, inCreated, inCmnt)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotStatusErr, ok := gotAPIErr.Err.(*StatusCodeError); !ok {
			t.Errorf("err.Err: expected (*ambidata.StatusCodeError), got %T", gotAPIErr.Err)
		} else if gotStatusErr.StatusCode != http.StatusNotFound {
			t.Errorf("err.StatusCode: expected %d, got %d", http.StatusNotFound, gotStatusErr.StatusCode)
		}
	}
}

func TestSenderSetCmntErrRequestEntityTooLarge(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inCmnt = "comment"
	const inCode = 200
	const inBody = "request entity too large"
	const wantMethod = "PUT"
	const wantPath = "/api/v2/channels/83601/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetCmnt(ctx, inCreated, inCmnt)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotAPIErr.Err != ErrRequestEntityTooLarge {
			t.Errorf("err:Err: expected ErrRequestEntityTooLarge, got %q", gotAPIErr.Err.Error())
		}
	}
}

func TestSenderSetHideNormal(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inHide = true
	const inCode = 200
	const inBody = "OK"
	const wantCT = "application/json"

	wantJSON := map[string]any{
		"writeKey": inWriteKey,
		"created":  "1970-01-01T01:00:00Z",
		"hide":     true,
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	var gotReqBody []byte
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("PUT /api/v2/channels/"+inCh+"/data", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		gotReqBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetHide(ctx, inCreated, inHide)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}
	if gotCT := gotReq.Header.Get("Content-Type"); gotCT != wantCT {
		t.Errorf("request: Content-Type: expected %#v, got %#v", wantCT, gotCT)
	}

	var gotJSON map[string]any
	if err := json.Unmarshal(gotReqBody, &gotJSON); err != nil {
		t.Errorf("response: body: %v", err)
	} else if diff := cmp.Diff(wantJSON, gotJSON); diff != "" {
		t.Errorf("response: body: mismatch (-want, +got)\n%s", diff)
	}
}

func TestSenderSetHideErrCanceled(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inHide = true
	const inCode = 200
	const inBody = "OK"

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetHide(ctx, inCreated, inHide)
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestSenderSetHideErrStatus(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inHide = true
	const wantMethod = "PUT"
	const wantPath = "/api/v2/channels/83601/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetHide(ctx, inCreated, inHide)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotStatusErr, ok := gotAPIErr.Err.(*StatusCodeError); !ok {
			t.Errorf("err.Err: expected (*ambidata.StatusCodeError), got %T", gotAPIErr.Err)
		} else if gotStatusErr.StatusCode != http.StatusNotFound {
			t.Errorf("err.StatusCode: expected %d, got %d", http.StatusNotFound, gotStatusErr.StatusCode)
		}
	}
}

func TestSenderSetHideErrRequestEntityTooLarge(t *testing.T) {
	const inCh = "83601"
	const inWriteKey = "52e2cd7ddbfe2fed"
	inCreated := time.Date(1970, 1, 1, 1, 0, 0, 0, time.UTC)
	const inHide = true
	const inCode = 200
	const inBody = "request entity too large"
	const wantMethod = "PUT"
	const wantPath = "/api/v2/channels/83601/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(inCode)
		w.Write([]byte(inBody))
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	s := &Sender{
		Ch:       inCh,
		WriteKey: inWriteKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := s.SetHide(ctx, inCreated, inHide)
	if gotAPIErr, ok := gotErr.(*APIError); !ok {
		t.Errorf("err: expected (*ambidata.APIError), got %T", gotErr)
	} else {
		if gotAPIErr.Method != wantMethod {
			t.Errorf("err.Method: expected %#v, got %#v", wantMethod, gotAPIErr.Method)
		}
		if gotAPIErr.Path != wantPath {
			t.Errorf("err.Path: expected %#v, got %#v", wantPath, gotAPIErr.Path)
		}
		if diff := cmp.Diff(wantQuery, gotAPIErr.Query); diff != "" {
			t.Errorf("err.Query: mismatch (-want, +got)\n%s", diff)
		}
		if gotAPIErr.Err != ErrRequestEntityTooLarge {
			t.Errorf("err:Err: expected ErrRequestEntityTooLarge, got %q", gotAPIErr.Err.Error())
		}
	}
}
