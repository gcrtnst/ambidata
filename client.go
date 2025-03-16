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
	Method     string
	Path       string
	StatusCode int
}

func (err *APIError) Error() string {
	if err == nil || err.Method == "" || err.Path == "" || err.StatusCode == 0 {
		return fmt.Sprintf("%#v", err)
	}

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString("ambidata: ")
	b.WriteString(err.Method)
	b.WriteByte(' ')
	b.WriteString(err.Path)
	b.WriteString(": ")

	statusCode := err.StatusCode
	statusText := http.StatusText(statusCode)
	b.WriteString(strconv.Itoa(statusCode))
	if statusText != "" {
		b.WriteByte(' ')
		b.WriteString(statusText)
	}

	return b.String()
}
