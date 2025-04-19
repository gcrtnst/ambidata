package main

import (
	"context"
	"errors"
	"math"
	"slices"
	"strings"
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

func TestSenderSendTimePrecision(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	// "推測に基づく情報" 用のテスト
	in := ambidata.Data{Created: time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.UTC)}
	want := ambidata.Data{Created: time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC)}

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.Send(ctx, in)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	got, errFetch := f.FetchRange(ctx, 1, 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", []ambidata.Data{want}, got)
}

func TestSenderSendCmntSize(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	// "推測に基づく情報" 用のテスト
	created := time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC)
	inCmnt := strings.Repeat(".", 128)
	wantCmnt := strings.Repeat(".", 64)
	in := ambidata.Data{Created: created, Cmnt: inCmnt}
	want := ambidata.Data{Created: created, Cmnt: wantCmnt}

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.Send(ctx, in)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	got, errFetch := f.FetchRange(ctx, 1, 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", []ambidata.Data{want}, got)
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

func TestSenderSendBulkTooLarge(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	data := ambidata.Data{
		Created: time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.FixedZone("UTC+7", 7*60*60)),
		D1:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D2:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D3:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D4:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D5:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D6:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D7:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		D8:      ambidata.Just(-math.Nextafter(1e-6, 0)),
		Loc:     ambidata.Just(ambidata.Location{Lat: -math.Nextafter(1e-6, 0), Lng: -math.Nextafter(1e-6, 0)}),
		Cmnt:    strings.Repeat("-", 64),
	}

	// "推測に基づく情報" 用のテスト
	const maxlen = 258
	arrOK := slices.Repeat([]ambidata.Data{data}, maxlen)
	arrNG := slices.Repeat([]ambidata.Data{data}, maxlen+1)

	errDelOK := m.DeleteData(ctx, t.Config.Ch)
	if errDelOK != nil {
		t.Error(errDelOK)
		return
	}

	t.PostWait()
	errSendOK := s.SendBulk(ctx, arrOK)
	t.PostDone()
	if errSendOK != nil {
		t.Errorf("errSendOK: expected nil, got %#v", errSendOK)
		return
	}

	errDelNG := m.DeleteData(ctx, t.Config.Ch)
	if errDelNG != nil {
		t.Error(errDelOK)
		return
	}

	t.PostWait()
	errSendNG := s.SendBulk(ctx, arrNG)
	t.PostDone()
	if !errors.Is(errSendNG, ambidata.ErrRequestEntityTooLarge) {
		t.Errorf("errSendNG: expected ErrRequestEntityTooLarge, got %#v", errSendNG)
		return
	}
}

func TestSenderSetCmnt(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	created := time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC)
	const cmntOld = "cmnt old"
	const cmntNew = "cmnt new"

	sent := ambidata.Data{
		Created: created,
		D1:      ambidata.Just(101.0),
		D2:      ambidata.Just(102.0),
		D3:      ambidata.Just(103.0),
		D4:      ambidata.Just(104.0),
		D5:      ambidata.Just(105.0),
		D6:      ambidata.Just(106.0),
		D7:      ambidata.Just(107.0),
		D8:      ambidata.Just(108.0),
		Loc:     ambidata.Just(ambidata.Location{Lat: 109.0, Lng: 110.0}),
		Cmnt:    cmntOld,
	}
	want := sent
	want.Cmnt = cmntNew

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.Send(ctx, sent)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	s.SetCmnt(ctx, created, cmntNew)

	a, errFetch := f.FetchRange(ctx, 1, 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", []ambidata.Data{want}, a)
}

func TestSenderSetHide(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	created := time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC)
	sent := ambidata.Data{
		Created: created,
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
	want := sent
	want.Hide = true

	errDel := m.DeleteData(ctx, t.Config.Ch)
	if errDel != nil {
		t.Error(errDel)
		return
	}

	t.PostWait()
	errSend := s.Send(ctx, sent)
	t.PostDone()
	if errSend != nil {
		t.Error(errSend)
		return
	}

	s.SetHide(ctx, created, true)

	a, errFetch := f.FetchRange(ctx, 1, 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", []ambidata.Data{want}, a)
}

func TestSenderSetHideNonexistent(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	// "推測に基づく情報" 用のテスト
	sentCreated := time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC)
	setCreated := time.Date(2006, 1, 2, 15, 4, 5, 998000000, time.UTC)

	sent := []ambidata.Data{
		{
			Created: sentCreated,
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
	want := slices.Clone(sent)

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

	s.SetHide(ctx, setCreated, true)

	got, errFetch := f.FetchRange(ctx, len(want), 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", want, got)
}

func TestSenderSetHideMultiple(t *T) {
	ctx := context.Background()
	m := ambidata.NewManager(t.Config.UserKey)
	f := ambidata.NewFetcher(t.Config.Ch, t.Config.ReadKey)
	s := ambidata.NewSender(t.Config.Ch, t.Config.WriteKey)

	// "推測に基づく情報" 用のテスト
	created := time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC)
	sent := []ambidata.Data{
		{
			Created: created,
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
		{
			Created: created,
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
	}
	want := slices.Clone(sent)
	want[0].Hide = true

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

	s.SetHide(ctx, created, true)

	got, errFetch := f.FetchRange(ctx, len(want), 0)
	if errFetch != nil {
		t.Error(errFetch)
		return
	}

	assertCmp(t, "dataarray: ", want, got)
}
