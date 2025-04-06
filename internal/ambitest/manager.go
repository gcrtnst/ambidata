package main

import (
	"context"
	"time"

	"github.com/gcrtnst/ambidata"
)

func TestManagerGetChannelList(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
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

	l, err := m.GetChannelList(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	i := -1
	for j := range l {
		if l[j].Ch == t.Config.Ch {
			i = j
			break
		}
	}
	if i < 0 {
		t.Error("target channel not found")
		return
	}

	assertAtoi(t, "user: ", l[i].User)
	assertEqual(t, "readKey: ", t.Config.ReadKey, l[i].ReadKey)
	assertEqual(t, "writeKey: ", t.Config.WriteKey, l[i].WriteKey)
	assertTimeIsBetween(t, "created: ", stt, end, l[i].Created)
	assertTimeIsBetween(t, "modified: ", stt, end, l[i].Modified)
	assertTimeIsBetween(t, "lastpost: ", stt, end, l[i].LastPost)
	assertEqual(t, "charts: ", 0, l[i].Charts)
	assertGreaterOrEqual(t, "dataperday: ", wantDataPerDayMin, l[i].DataPerDay)
	assertEqual(t, "d_ch: ", wantDCh, l[i].DCh)
	assertNotEqual(t, "chName: ", "", l[i].ChName)
	assertNotEqual(t, "chDesc: ", "", l[i].ChDesc)
	assertNotEqual(t, "d1.name: ", "", l[i].D1.Name)
	assertColor(t, "d1.color: ", l[i].D1.Color)
	assertNotEqual(t, "d2.name: ", "", l[i].D2.Name)
	assertColor(t, "d2.color: ", l[i].D2.Color)
	assertNotEqual(t, "d3.name: ", "", l[i].D3.Name)
	assertColor(t, "d3.color: ", l[i].D3.Color)
	assertNotEqual(t, "d4.name: ", "", l[i].D4.Name)
	assertColor(t, "d4.color: ", l[i].D4.Color)
	assertNotEqual(t, "d5.name: ", "", l[i].D5.Name)
	assertColor(t, "d5.color: ", l[i].D5.Color)
	assertNotEqual(t, "d6.name: ", "", l[i].D6.Name)
	assertColor(t, "d6.color: ", l[i].D6.Color)
	assertNotEqual(t, "d7.name: ", "", l[i].D7.Name)
	assertColor(t, "d7.color: ", l[i].D7.Color)
	assertNotEqual(t, "d8.name: ", "", l[i].D8.Name)
	assertColor(t, "d8.color: ", l[i].D8.Color)
	assertNonZeroLocation(t, "loc: ", l[i].Loc)
	assertNotEqual(t, "photoid: ", "", l[i].PhotoID)
	assertContains(t, "devkeys: ", t.Config.DevKey, l[i].DevKeys)
	assertAtoi(t, "bd: ", l[i].Bd)
	assertNotEqual(t, "lastdata._id: ", "", l[i].LastData.ID)
	assertCmp(t, "lastdata: ", data, l[i].LastData.Data)
}

func TestManagerGetDeviceChannel(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
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

	c, err := m.GetDeviceChannel(ctx, t.Config.DevKey)
	if err != nil {
		t.Error(err)
		return
	}

	assertEqual(t, "ch: ", t.Config.Ch, c.Ch)
	assertAtoi(t, "user: ", c.User)
	assertEqual(t, "readKey: ", t.Config.ReadKey, c.ReadKey)
	assertEqual(t, "writeKey: ", t.Config.WriteKey, c.WriteKey)
	assertTimeIsBetween(t, "created: ", stt, end, c.Created)
	assertTimeIsBetween(t, "modified: ", stt, end, c.Modified)
	assertTimeIsBetween(t, "lastpost: ", stt, end, c.LastPost)
	assertEqual(t, "charts: ", 0, c.Charts)
	assertGreaterOrEqual(t, "dataperday: ", wantDataPerDayMin, c.DataPerDay)
	assertEqual(t, "d_ch: ", wantDCh, c.DCh)
	assertNotEqual(t, "chName: ", "", c.ChName)
	assertNotEqual(t, "chDesc: ", "", c.ChDesc)
	assertNotEqual(t, "d1.name: ", "", c.D1.Name)
	assertColor(t, "d1.color: ", c.D1.Color)
	assertNotEqual(t, "d2.name: ", "", c.D2.Name)
	assertColor(t, "d2.color: ", c.D2.Color)
	assertNotEqual(t, "d3.name: ", "", c.D3.Name)
	assertColor(t, "d3.color: ", c.D3.Color)
	assertNotEqual(t, "d4.name: ", "", c.D4.Name)
	assertColor(t, "d4.color: ", c.D4.Color)
	assertNotEqual(t, "d5.name: ", "", c.D5.Name)
	assertColor(t, "d5.color: ", c.D5.Color)
	assertNotEqual(t, "d6.name: ", "", c.D6.Name)
	assertColor(t, "d6.color: ", c.D6.Color)
	assertNotEqual(t, "d7.name: ", "", c.D7.Name)
	assertColor(t, "d7.color: ", c.D7.Color)
	assertNotEqual(t, "d8.name: ", "", c.D8.Name)
	assertColor(t, "d8.color: ", c.D8.Color)
	assertNonZeroLocation(t, "loc: ", c.Loc)
	assertNotEqual(t, "photoid: ", "", c.PhotoID)
	assertContains(t, "devkeys: ", t.Config.DevKey, c.DevKeys)
	assertAtoi(t, "bd: ", c.Bd)
	assertNotEqual(t, "lastdata._id: ", "", c.LastData.ID)
	assertCmp(t, "lastdata: ", data, c.LastData.Data)
}

func TestGetDeviceChannelLv1(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)

	c, err := m.GetDeviceChannelLv1(ctx, t.Config.DevKey)
	if err != nil {
		t.Error(err)
		return
	}

	assertEqual(t, "ch: ", t.Config.Ch, c.Ch)
	assertEqual(t, "writeKey: ", t.Config.WriteKey, c.WriteKey)
}
