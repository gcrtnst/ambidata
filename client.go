package ambidata

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Config struct {
	Scheme string
	Host   string
	Client *http.Client
}

type APIError struct {
	Method string
	Path   string
	Err    error
}

func (err *APIError) Error() string {
	if err == nil || err.Method == "" || err.Path == "" || err.Err == nil {
		return fmt.Sprintf("%#v", err)
	}

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString("ambidata: ")
	b.WriteString(err.Method)
	b.WriteByte(' ')
	b.WriteString(err.Path)
	b.WriteString(": ")
	b.WriteString(err.Err.Error())
	return b.String()
}
