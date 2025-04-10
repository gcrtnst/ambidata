package main

import (
	"context"
	"time"

	"github.com/gcrtnst/ambidata"
)

func TestFetcherGetChannel(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	stt := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().Add(time.Hour)
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

	const wantDataPerDayMin = 1
	const wantDCh = true

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	// Send data to ambidata to ensure that fields like lastpost, lastdata and
	// dataperday will have non-zero values when retrived.
	t.PostWait()
	errSend := s.Send(ctx, data)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	c, err := f.GetChannel(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	assertEqual(t, "ch: ", t.Config.Ch, c.Ch)
	assertAtoi(t, "user: ", c.User)
	assertTimeIsBetween(t, "created: ", stt, end, c.Created)
	assertTimeIsBetween(t, "modified: ", stt, end, c.Modified)
	assertTimeIsBetween(t, "lastpost: ", stt, end, c.LastPost)
	assertEqual(t, "charts: ", 0, c.Charts)
	assertGreaterOrEqual(t, "dataperday: ", wantDataPerDayMin, c.DataPerDay)
	assertEqual(t, "d_ch: ", wantDCh, c.DCh)
	assertNotEqual(t, "chName: ", "", c.ChName)
	assertNotEqual(t, "chDesc: ", "", c.ChDesc)
	assertNotEqual(t, "d1.name: ", "", c.D1.Name)
	assertFieldColor(t, "d1.color: ", c.D1.Color)
	assertNotEqual(t, "d2.name: ", "", c.D2.Name)
	assertFieldColor(t, "d2.color: ", c.D2.Color)
	assertNotEqual(t, "d3.name: ", "", c.D3.Name)
	assertFieldColor(t, "d3.color: ", c.D3.Color)
	assertNotEqual(t, "d4.name: ", "", c.D4.Name)
	assertFieldColor(t, "d4.color: ", c.D4.Color)
	assertNotEqual(t, "d5.name: ", "", c.D5.Name)
	assertFieldColor(t, "d5.color: ", c.D5.Color)
	assertNotEqual(t, "d6.name: ", "", c.D6.Name)
	assertFieldColor(t, "d6.color: ", c.D6.Color)
	assertNotEqual(t, "d7.name: ", "", c.D7.Name)
	assertFieldColor(t, "d7.color: ", c.D7.Color)
	assertNotEqual(t, "d8.name: ", "", c.D8.Name)
	assertFieldColor(t, "d8.color: ", c.D8.Color)
	assertNonZeroLocation(t, "loc: ", c.Loc)
	assertNotEqual(t, "photoid: ", "", c.PhotoID)
	assertContains(t, "devkeys: ", t.Config.DevKey, c.DevKeys)
	assertAtoi(t, "bd: ", c.Bd)
	assertNotEqual(t, "lastdata._id: ", "", c.LastData.ID)
	assertCmp(t, "lastdata: ", data, c.LastData.Data)
}

func TestFetcherFetchRange(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	sent := []ambidata.Data{
		{
			Created: time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(501.0),
			D2:      ambidata.Just(502.0),
			D3:      ambidata.Just(503.0),
			D4:      ambidata.Just(504.0),
			D5:      ambidata.Just(505.0),
			D6:      ambidata.Just(506.0),
			D7:      ambidata.Just(507.0),
			D8:      ambidata.Just(508.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 509.0, Lng: 510.0}),
			Cmnt:    "cmnt 5",
		},
		{
			Created: time.Date(2015, 1, 4, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(401.0),
			D2:      ambidata.Just(402.0),
			D3:      ambidata.Just(403.0),
			D4:      ambidata.Just(404.0),
			D5:      ambidata.Just(405.0),
			D6:      ambidata.Just(406.0),
			D7:      ambidata.Just(407.0),
			D8:      ambidata.Just(408.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 409.0, Lng: 410.0}),
			Cmnt:    "cmnt 4",
		},
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
	want := sent[1:4]

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.SendBulk(ctx, sent)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	got, errFetch := f.FetchRange(ctx, 3, 1)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", want, got)
}

func TestFetcherFetchPeriod(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	stt := time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)
	end := time.Date(2015, 1, 4, 0, 0, 0, 1000000, time.UTC)
	sent := []ambidata.Data{
		{
			Created: time.Date(2015, 1, 5, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(501.0),
			D2:      ambidata.Just(502.0),
			D3:      ambidata.Just(503.0),
			D4:      ambidata.Just(504.0),
			D5:      ambidata.Just(505.0),
			D6:      ambidata.Just(506.0),
			D7:      ambidata.Just(507.0),
			D8:      ambidata.Just(508.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 509.0, Lng: 510.0}),
			Cmnt:    "cmnt 5",
		},
		{
			Created: time.Date(2015, 1, 4, 0, 0, 0, 0, time.UTC),
			D1:      ambidata.Just(401.0),
			D2:      ambidata.Just(402.0),
			D3:      ambidata.Just(403.0),
			D4:      ambidata.Just(404.0),
			D5:      ambidata.Just(405.0),
			D6:      ambidata.Just(406.0),
			D7:      ambidata.Just(407.0),
			D8:      ambidata.Just(408.0),
			Loc:     ambidata.Just(ambidata.Location{Lat: 409.0, Lng: 410.0}),
			Cmnt:    "cmnt 4",
		},
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
	want := sent[1:4]

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.SendBulk(ctx, sent)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	got, errFetch := f.FetchPeriod(ctx, stt, end)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", want, got)
}
