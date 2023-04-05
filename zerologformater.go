package github.com/lastbyte32/chi-zerologformater

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

var (
	nRed     = []byte{'\033', '[', '3', '1', 'm'}
	nGreen   = []byte{'\033', '[', '3', '2', 'm'}
	nYellow  = []byte{'\033', '[', '3', '3', 'm'}
	nCyan    = []byte{'\033', '[', '3', '6', 'm'}
	bRed     = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen   = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	bYellow  = []byte{'\033', '[', '3', '3', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	reset    = []byte{'\033', '[', '0', 'm'}
)

var (
	_ middleware.LogFormatter = (*zeroLogFormatter)(nil)
)

type zeroLogFormatter struct {
	logger *zerolog.Logger
}

type zeroLogEntry struct {
	log     *zerolog.Logger
	request *http.Request
	buf     *bytes.Buffer
}

func New(l *zerolog.Logger) middleware.LogFormatter {
	return &zeroLogFormatter{
		logger: l,
	}
}

func (z *zeroLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {

	switch {
	case status < 200:
		cW(z.buf, bBlue, "%03d", status)
	case status < 300:
		cW(z.buf, bGreen, "%03d", status)
	case status < 400:
		cW(z.buf, bCyan, "%03d", status)
	case status < 500:
		cW(z.buf, bYellow, "%03d", status)
	default:
		cW(z.buf, bRed, "%03d", status)
	}

	cW(z.buf, bBlue, " %dB", bytes)

	z.buf.WriteString(" in ")
	if elapsed < 500*time.Millisecond {
		cW(z.buf, nGreen, "%s", elapsed)
	} else if elapsed < 5*time.Second {
		cW(z.buf, nYellow, "%s", elapsed)
	} else {
		cW(z.buf, nRed, "%s", elapsed)
	}

	z.log.Info().
		//Str("method", z.request.Method).
		Msg(z.buf.String())

}

func (z *zeroLogEntry) Panic(v interface{}, stack []byte) {
	z.log.Info().Msgf("request failed: %+v", v)
}

func (l zeroLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &zeroLogEntry{
		log:     l.logger,
		request: r,
		buf:     &bytes.Buffer{},
	}
	cW(entry.buf, nCyan, "\"")
	cW(entry.buf, bMagenta, "%s ", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	cW(entry.buf, nCyan, "%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)

	entry.buf.WriteString("from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteString(" - ")

	return entry
}

func cW(w io.Writer, color []byte, s string, args ...interface{}) {
	w.Write(color)
	fmt.Fprintf(w, s, args...)
	w.Write(reset)
}
