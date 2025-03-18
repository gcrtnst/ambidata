package ambidata

import (
	"net/http"
	"testing"
)

func TestAPIErrorError(t *testing.T) {
	tt := []struct {
		name string
		in   *APIError
		want string
	}{
		{
			name: "Nil",
			in:   nil,
			want: "(*ambidata.APIError)(nil)",
		},
		{
			name: "EmptyMethod",
			in:   &APIError{Path: "/", Err: mockError{}},
			want: `&ambidata.APIError{Method:"", Path:"/", Err:ambidata.mockError{}}`,
		},
		{
			name: "EmptyPath",
			in:   &APIError{Method: "GET", Err: mockError{}},
			want: `&ambidata.APIError{Method:"GET", Path:"", Err:ambidata.mockError{}}`,
		},
		{
			name: "EmptyStatusCode",
			in:   &APIError{Method: "GET", Path: "/"},
			want: `&ambidata.APIError{Method:"GET", Path:"/", Err:error(nil)}`,
		},
		{
			name: "Normal",
			in:   &APIError{Method: "GET", Path: "/", Err: mockError{}},
			want: "ambidata: GET /: mock error",
		},
	}

	for _, tc := range tt {
		got := tc.in.Error()
		if got != tc.want {
			t.Errorf("%s: expected %#v, got %#v", tc.name, tc.want, got)
		}
	}
}

func TestStatusCodeError(t *testing.T) {
	tt := []struct {
		name string
		in   *StatusCodeError
		want string
	}{
		{
			name: "Nil",
			in:   nil,
			want: "(*ambidata.StatusCodeError)(nil)",
		},
		{
			name: "EmptyStatusCode",
			in:   &StatusCodeError{},
			want: "&ambidata.StatusCodeError{StatusCode:0}",
		},
		{
			name: "Normal",
			in:   &StatusCodeError{StatusCode: http.StatusNotFound},
			want: "404 Not Found",
		},
		{
			name: "UnknownStatusCode",
			in:   &StatusCodeError{StatusCode: 999},
			want: "999 Unknown Status Code",
		},
	}

	for _, tc := range tt {
		got := tc.in.Error()
		if got != tc.want {
			t.Errorf("%s: expected %#v, got %#v", tc.name, tc.want, got)
		}
	}
}

type mockError struct{}

func (err mockError) Error() string {
	return "mock error"
}
