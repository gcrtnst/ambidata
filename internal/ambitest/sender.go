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

func TestSenderSendBulk(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	a1 := []ambidata.Data{
		{
			Created: time.Date(2015, 1, 3, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(301.0),
			D2:      ambidata.Just(302.0),
			D3:      ambidata.Just(303.0),
			D4:      ambidata.Just(304.0),
			D5:      ambidata.Just(305.0),
			D6:      ambidata.Just(306.0),
			D7:      ambidata.Just(307.0),
			D8:      ambidata.Just(308.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 309.0, Lng: 310.0}),
			Cmnt:    "cmnt 3",
		},
		{
			Created: time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(201.0),
			D2:      ambidata.Just(202.0),
			D3:      ambidata.Just(203.0),
			D4:      ambidata.Just(204.0),
			D5:      ambidata.Just(205.0),
			D6:      ambidata.Just(206.0),
			D7:      ambidata.Just(207.0),
			D8:      ambidata.Just(208.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 209.0, Lng: 210.0}),
			Cmnt:    "cmnt 2",
		},
		{
			Created: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(101.0),
			D2:      ambidata.Just(102.0),
			D3:      ambidata.Just(103.0),
			D4:      ambidata.Just(104.0),
			D5:      ambidata.Just(105.0),
			D6:      ambidata.Just(106.0),
			D7:      ambidata.Just(107.0),
			D8:      ambidata.Just(108.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 109.0, Lng: 110.0}),
			Cmnt:    "cmnt 1",
		},
	}

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.SendBulk(ctx, a1)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	a2, errFetch := f.FetchRange(ctx, len(a1), 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", a1, a2)
}
