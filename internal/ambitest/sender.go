package main

import (
	"context"
	"time"

	"github.com/gcrtnst/ambidata"
)

func TestSenderSend(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	data := ambidata.Data{
		Created: time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC),
		D1:      ambidata.Just(101.0),
		D2:      ambidata.Just(102.0),
		D3:      ambidata.Just(103.0),
		D4:      ambidata.Just(104.0),
		D5:      ambidata.Just(105.0),
		D6:      ambidata.Just(106.0),
		D7:      ambidata.Just(107.0),
		D8:      ambidata.Just(108.0),
		Loc:     ambidata.Just(ambidata.Location{Lat: 109.0, Lng: 110.0}),
		Cmnt:    "cmnt",
	}

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.Send(ctx, data)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	a, errFetch := f.FetchRange(ctx, 1, 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", []ambidata.Data{data}, a)
}
