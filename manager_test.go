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

func TestManagerGetChannelListNormal(t *testing.T) {
	const inBody = `[{"ch":"83602","readKey":"234e4ce9ec33dacf","writeKey":"575d743ebb4d2c2d","user":"27143","created":"2006-01-02T15:04:05.999Z","modified":"2006-01-02T15:04:05.999Z","lastpost":"1970-01-01T00:00:00.000Z","charts":0,"dataperday":0,"d_ch":true},{"ch":"83601","readKey":"74545caba2bfd44f","writeKey":"52e2cd7ddbfe2fed","user":"27143","created":"2006-01-03T15:04:05.999Z","modified":"2006-01-04T15:04:05.999Z","lastpost":"2006-01-05T15:04:05.999Z","charts":0,"dataperday":288,"d_ch":true,"bd":"98929","devkeys":["08:A9:0C:9E:E0:C3"],"chDesc":"chDesc","chName":"chName","d1":{"name":"d1","color":"1"},"d2":{"name":"d2","color":"2"},"d3":{"name":"d3","color":"3"},"d4":{"name":"d4","color":"4"},"d5":{"name":"d5","color":"5"},"d6":{"name":"d6","color":"6"},"d7":{"name":"d7","color":"7"},"d8":{"name":"d8","color":"8"},"loc":[9,10],"photoid":"https://drive.google.com/file/d/1MK59q8lV8tDZCOvvjjudWerKHYUejBCt/view?usp=sharing","lastdata":{"d1":1,"d2":2,"d3":3,"d4":4,"d5":5,"d6":6,"d7":7,"d8":8,"loc":[10,9],"cmnt":"cmnt","created":"2006-01-05T15:04:05.999Z","_id":"67d82f1f81e5845e0e8e9b8d"}}]`
	const inUserKey = "4ef42dcecf7e7ceba2"

	wantRet := []ChannelAccess{
		{
			ReadKey:  "234e4ce9ec33dacf",
			WriteKey: "575d743ebb4d2c2d",
			ChannelInfo: ChannelInfo{
				Ch:         "83602",
				User:       "27143",
				Created:    time.Date(2006, 1, 2, 15, 04, 05, 999000000, time.UTC),
				Modified:   time.Date(2006, 1, 2, 15, 04, 05, 999000000, time.UTC),
				LastPost:   time.Time{},
				Charts:     0,
				DataPerDay: 0,
				DCh:        true,
			},
		},
		{
			ReadKey:  "74545caba2bfd44f",
			WriteKey: "52e2cd7ddbfe2fed",
			ChannelInfo: ChannelInfo{
				Ch:         "83601",
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
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("GET /api/v2/channels/{$}", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: inUserKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotRet, gotErr := m.GetChannelList(ctx)
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
	} else if gotUserKey := gotQuery.Get("userKey"); gotUserKey != inUserKey {
		t.Errorf("request: userKey: expected %#v, got %#v", inUserKey, gotUserKey)
	}
}

func TestManagerGetChannelListErrCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetChannelList(ctx)
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestManagerGetChannelListErrStatus(t *testing.T) {
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetChannelList(ctx)
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

func TestManagerGetChannelListErrJSON(t *testing.T) {
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/"
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

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetChannelList(ctx)
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

func TestManagerGetDeviceChannelNormal(t *testing.T) {
	const inBody = `{"ch":"83601","readKey":"74545caba2bfd44f","writeKey":"52e2cd7ddbfe2fed","user":"27143","created":"2006-01-03T15:04:05.999Z","modified":"2006-01-04T15:04:05.999Z","lastpost":"2006-01-05T15:04:05.999Z","charts":0,"dataperday":288,"d_ch":true,"bd":"98929","devkeys":["08:A9:0C:9E:E0:C3"],"chDesc":"chDesc","chName":"chName","d1":{"name":"d1","color":"1"},"d2":{"name":"d2","color":"2"},"d3":{"name":"d3","color":"3"},"d4":{"name":"d4","color":"4"},"d5":{"name":"d5","color":"5"},"d6":{"name":"d6","color":"6"},"d7":{"name":"d7","color":"7"},"d8":{"name":"d8","color":"8"},"loc":[9,10],"photoid":"https://drive.google.com/file/d/1MK59q8lV8tDZCOvvjjudWerKHYUejBCt/view?usp=sharing","lastdata":{"d1":1,"d2":2,"d3":3,"d4":4,"d5":5,"d6":6,"d7":7,"d8":8,"loc":[10,9],"cmnt":"cmnt","created":"2006-01-05T15:04:05.999Z","_id":"67d82f1f81e5845e0e8e9b8d"}}`
	const inUserKey = "4ef42dcecf7e7ceba2"
	const inDevKey = "08:A9:0C:9E:E0:C3"

	wantRet := ChannelAccess{
		ReadKey:  "74545caba2bfd44f",
		WriteKey: "52e2cd7ddbfe2fed",
		ChannelInfo: ChannelInfo{
			Ch:         "83601",
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
		},
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("GET /api/v2/channels/{$}", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: inUserKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotRet, gotErr := m.GetDeviceChannel(ctx, inDevKey)
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
	} else {
		if gotUserKey := gotQuery.Get("userKey"); gotUserKey != inUserKey {
			t.Errorf("request: userKey: expected %#v, got %#v", inUserKey, gotUserKey)
		}
		if gotDevKey := gotQuery.Get("devKey"); gotDevKey != inDevKey {
			t.Errorf("request: devKey: expected %#v, got %#v", inDevKey, gotDevKey)
		}
	}
}

func TestManagerGetDeviceChannelErrCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetDeviceChannel(ctx, "08:A9:0C:9E:E0:C3")
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestManagerGetDeviceChannelErrStatus(t *testing.T) {
	const inDevKey = "08:A9:0C:9E:E0:C3"
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/"
	wantQuery := url.Values{"devKey": []string{inDevKey}}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetDeviceChannel(ctx, inDevKey)
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

func TestManagerGetDeviceChannelErrJSON(t *testing.T) {
	const inDevKey = "08:A9:0C:9E:E0:C3"
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/"
	wantQuery := url.Values{"devKey": []string{inDevKey}}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetDeviceChannel(ctx, inDevKey)
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

func TestManagerGetDeviceChannelLv1Normal(t *testing.T) {
	const inBody = `{"ch":"83601","writeKey":"52e2cd7ddbfe2fed"}`
	const inUserKey = "4ef42dcecf7e7ceba2"
	const inDevKey = "08:A9:0C:9E:E0:C3"
	const wantLevel = "1"

	wantRet := ChannelAccessLv1{
		WriteKey: "52e2cd7ddbfe2fed",
		Ch:       "83601",
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("GET /api/v2/channels/{$}", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		w.Write([]byte(inBody))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: inUserKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotRet, gotErr := m.GetDeviceChannelLv1(ctx, inDevKey)
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
	} else {
		if gotUserKey := gotQuery.Get("userKey"); gotUserKey != inUserKey {
			t.Errorf("request: userKey: expected %#v, got %#v", inUserKey, gotUserKey)
		}
		if gotDevKey := gotQuery.Get("devKey"); gotDevKey != inDevKey {
			t.Errorf("request: devKey: expected %#v, got %#v", inDevKey, gotDevKey)
		}
		if gotLevel := gotQuery.Get("level"); gotLevel != wantLevel {
			t.Errorf("request: level: expected %#v, got %#v", wantLevel, gotLevel)
		}
	}
}

func TestManagerGetDeviceChannelLv1ErrCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetDeviceChannelLv1(ctx, "08:A9:0C:9E:E0:C3")
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestManagerGetDeviceChannelLv1ErrStatus(t *testing.T) {
	const inDevKey = "08:A9:0C:9E:E0:C3"
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/"
	wantQuery := url.Values{"devKey": []string{inDevKey}, "level": []string{"1"}}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetDeviceChannelLv1(ctx, inDevKey)
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

func TestManagerGetDeviceChannelLv1ErrJSON(t *testing.T) {
	const inDevKey = "08:A9:0C:9E:E0:C3"
	const wantMethod = "GET"
	const wantPath = "/api/v2/channels/"
	wantQuery := url.Values{"devKey": []string{inDevKey}, "level": []string{"1"}}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	_, gotErr := m.GetDeviceChannelLv1(ctx, inDevKey)
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

func TestManagerDeleteDataNormal(t *testing.T) {
	const inUserKey = "4ef42dcecf7e7ceba2"
	const inCh = "83602"

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	var gotReq *http.Request
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("DELETE /api/v2/channels/"+inCh+"/data", func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: inUserKey,
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := m.DeleteData(ctx, inCh)
	if gotErr != nil {
		t.Fatalf("err: %v", gotErr)
	}

	if gotUA := gotReq.Header.Values("User-Agent"); len(gotUA) > 0 {
		t.Errorf("request: User-Agent: expected not to send, got %#v", gotUA)
	}

	if gotQuery, err := url.ParseQuery(gotReq.URL.RawQuery); err != nil {
		t.Errorf("request: query: %v", err)
	} else if gotUserKey := gotQuery.Get("userKey"); gotUserKey != inUserKey {
		t.Errorf("request: userKey: expected %#v, got %#v", inUserKey, gotUserKey)
	}
}

func TestManagerDeleteDataErrCanceled(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	cancel()

	var handler http.HandlerFunc
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	srv := httptest.NewServer(handler)
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := m.DeleteData(ctx, "83602")
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("err: expected %#v, got %#v", context.Canceled.Error(), gotErr.Error())
	}
}

func TestManagerDeleteDataErrStatus(t *testing.T) {
	const inCh = "83602"
	const wantMethod = "DELETE"
	const wantPath = "/api/v2/channels/" + inCh + "/data"
	wantQuery := url.Values{}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	srvURL, _ := url.Parse(srv.URL)

	m := &Manager{
		UserKey: "4ef42dcecf7e7ceba2",
		Config: &Config{
			Scheme: srvURL.Scheme,
			Host:   srvURL.Host,
			Client: srv.Client(),
		},
	}

	gotErr := m.DeleteData(ctx, inCh)
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
