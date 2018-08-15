package probe

import (
	"net/http"
	"time"
)

type Result struct {
	Timestamp time.Time
	Probe     *Probe

	Code    int
	Body    string
	Headers http.Header

	Error error
}
