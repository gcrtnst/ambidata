package ambidata

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestFetcherGetChannelNormal(t *testing.T) {
	const inReadKey = "74545caba2bfd44f"
	const inCh = "83601"
	const inBody = `{"ch":"83601","user":"27143","created":"2006-01-03T15:04:05.999Z","modified":"2006-01-04T15:04:05.999Z","lastpost":"2006-01-05T15:04:05.999Z","charts":0,"dataperday":288,"d_ch":true,"bd":"98929","devkeys":["08:A9:0C:9E:E0:C3"],"chDesc":"chDesc","chName":"chName","d1":{"name":"d1","color":"1"},"d2":{"name":"d2","color":"2"},"d3":{"name":"d3","color":"3"},"d4":{"name":"d4","color":"4"},"d5":{"name":"d5","color":"5"},"d6":{"name":"d6","color":"6"},"d7":{"name":"d7","color":"7"},"d8":{"name":"d8","color":"8"},"loc":[9,10],"photoid":"https://drive.google.com/file/d/1MK59q8lV8tDZCOvvjjudWerKHYUejBCt/view?usp=sharing","lastdata":{"d1":1,"d2":2,"d3":3,"d4":4,"d5":5,"d6":6,"d7":7,"d8":8,"loc":[10,9],"cmnt":"cmnt","created":"2006-01-05T15:04:05.999Z","_id":"67d82f1f81e5845e0e8e9b8d"}}`

	wantRet := ChannelInfo{
		Ch:         inCh,
		User:       "27143",
		Created:    time.Date(2006, 1, 3, 15, 04, 05, 999000000, time.UTC),
		Modified:   time.Date(2006, 1, 4, 15, 04, 05, 999000000, time.UTC),
		LastPost:   time.Date(2006, 1, 5, 15, 04, 05, 999000000, time.UTC),
		Charts:     0,
		DataPerDay: 288,
		DCh:        true,
		Bd:         "98929",
		DevKeys:    []string{"08:A9:0C:9E:E0:C3"},
		ChDesc:     "chDesc",
		ChName:     "chName",
		D1:         FieldInfo{Name: "d1", Color: "1"},
		D2:         FieldInfo{Name: "d2", Color: "2"},
		D3:         FieldInfo{Name: "d3", Color: "3"},
		D4:         FieldInfo{Name: "d4", Color: "4"},
		D5:         FieldInfo{Name: "d5", Color: "5"},
		D6:         FieldInfo{Name: "d6", Color: "6"},
		D7:         FieldInfo{Name: "d7", Color: "7"},
		D8:         FieldInfo{Name: "d8", Color: "8"},
		Loc:        Just(Location{Lat: 10, Lng: 9}),
		PhotoID:    "https://drive.google.com/file/d/1MK59q8lV8tDZCOvvjjudWerKHYUejBCt/view?usp=sharing",
		LastData: LastData{
			ID: "67d82f1f81e5845e0e8e9b8d",
			Data: Data{
				D1:      Just(1.0),
				D2:      Just(2.0),
				D3:      Just(3.0),
				D4:      Just(4.0),
				D5:      Just(5.0),
				D6:      Just(6.0),
				D7:      Just(7.0),
				D8:      Just(8.0),
				Loc:     Just(Location{Lat: 9, Lng: 10}),
				Cmnt:    "cmnt",
				Created: time.Date(2006, 1, 5, 15, 04, 05, 999000000, time.UTC),
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("GET /api/v2/channels/"+inCh+"/", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	f := &Fetcher{
		Ch:      inCh,
		ReadKey: inReadKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotRet, gotErr := f.GetChannel(ctx)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if diff := cmp.Diff(wantRet, gotRet); diff != "" {
		t.Errorf("ret: mismatch (-want, +got)\n%s", diff)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}

	if gotQuery, err := url.ParseQuery(gotReq.URL.RawQuery); err != nil {
		t.Errorf("request: query: %v", err)
	} else if gotReadKey := gotQuery.Get("readKey"); gotReadKey != inReadKey {
		t.Errorf("request: readKey: expected %#v, got %#v", inReadKey, gotReadKey)
	}
}

func TestFetcherGetChannelErrCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	f := &Fetcher{
		Ch:      "83601",
		ReadKey: "74545caba2bfd44f",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := f.GetChannel(ctx)
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestFetcherGetChannelErrStatus(t *testing.T) {
	const inReadKey = "74545caba2bfd44f"
	const inCh = "83601"
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/83601/"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	f := &Fetcher{
		Ch:      inCh,
		ReadKey: inReadKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := f.GetChannel(ctx)
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

func TestFetcherGetChannelErrJSON(t *testing.T) {
	const inReadKey = "74545caba2bfd44f"
	const inCh = "83601"
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/83601/"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	f := &Fetcher{
		Ch:      inCh,
		ReadKey: inReadKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := f.GetChannel(ctx)
	if gotAPIErr := (&APIError{}); !errors.As(gotErr, &gotAPIErr) {
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
	}
	if !errors.Is(gotErr, io.EOF) {
		t.Errorf("err: expected %#v, got %#v", io.EOF.Error(), gotErr.Error())
	}
}
