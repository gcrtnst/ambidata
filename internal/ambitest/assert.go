package main

import (
	"cmp"
	"slices"
	"strconv"
	"time"

	"github.com/gcrtnst/ambidata"
	gocmp "github.com/google/go-cmp/cmp"
)

func assertEqual[U comparable](t *T, prefix string, want U, got U) {
	if got != want {
		t.Errorf("%sexpected %#v, got %#v", prefix, want, got)
	}
}

func assertNotEqual[U comparable](t *T, prefix string, want U, got U) {
	if got == want {
		t.Errorf("%sexpected != %#v, got %#v", prefix, want, got)
	}
}

func assertGreaterOrEqual[U cmp.Ordered](t *T, prefix string, want U, got U) {
	if cmp.Compare(got, want) < 0 {
		t.Errorf("%sexpected >= %#v, got %#v", prefix, want, got)
	}
}

func assertEmptySlice[E any](t *T, prefix string, got []E) {
	if len(got) > 0 {
		t.Errorf("%sexpected an empty slice, got a slice of length %d", prefix, len(got))
	}
}

func assertContains[E comparable](t *T, prefix string, want E, got []E) {
	if !slices.Contains(got, want) {
		t.Errorf("%sexpected to contain %#v, got %#v", prefix, want, got)
	}
}

func assertCmp(t *T, prefix string, want any, got any) {
	if diff := gocmp.Diff(want, got); diff != "" {
		t.Errorf("%smismatch (-want, +got)\n%s", prefix, diff)
	}
}

func assertAtoi(t *T, prefix string, got string) {
	if _, err := strconv.Atoi(got); err != nil {
		t.Errorf("%sexpected a decimal number string, got %#v", prefix, got)
	}
}

func assertTimeIsBetween(t *T, prefix string, stt time.Time, end time.Time, got time.Time) {
	if got.Before(stt) || got.After(end) {
		t.Errorf("%sexpected to be between %s, %s, got %s", prefix, stt.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano), got.Format(time.RFC3339Nano))
	}
}

func assertFieldColor(t *T, prefix string, got ambidata.FieldColor) {
	if _, ok := got.ToRGBA(); !ok {
		t.Errorf("%sexpected a known color index, got %#v", prefix, got)
	}
}

func assertNonZeroLocation(t *T, prefix string, got ambidata.Maybe[ambidata.Location]) {
	if !got.OK {
		t.Errorf("%sexpected non-zero location, got null", prefix)
	} else if got.V.Lat == 0 || got.V.Lng == 0 {
		t.Errorf("%sexpected non-zero location, got (%f, %f)", prefix, got.V.Lat, got.V.Lng)
	}
}
