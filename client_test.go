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
			in:   &APIError{Path: "/", StatusCode: http.StatusNotFound},
			want: `&ambidata.APIError{Method:"", Path:"/", StatusCode:404}`,
		},
		{
			name: "EmptyPath",
			in:   &APIError{Method: "GET", StatusCode: http.StatusNotFound},
			want: `&ambidata.APIError{Method:"GET", Path:"", StatusCode:404}`,
		},
		{
			name: "EmptyStatusCode",
			in:   &APIError{Method: "GET", Path: "/"},
			want: `&ambidata.APIError{Method:"GET", Path:"/", StatusCode:0}`,
		},
		{
			name: "Normal",
			in:   &APIError{Method: "GET", Path: "/", StatusCode: http.StatusNotFound},
			want: "ambidata: GET /: 404 Not Found",
		},
		{
			name: "UnknownStatusCode",
			in:   &APIError{Method: "GET", Path: "/", StatusCode: 999},
			want: "ambidata: GET /: 999",
		},
	}

	for _, tc := range tt {
		got := tc.in.Error()
		if got != tc.want {
			t.Errorf("%s: expected %#v, got %#v", tc.name, tc.want, got)
		}
	}
}
