package zeroformater

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

var isTTY bool

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
	logger  *zerolog.Logger
	NoColor bool
}

type zeroLogEntry struct {
	log      *zerolog.Logger
	request  *http.Request
	buf      *bytes.Buffer
	useColor bool
}

func New(l *zerolog.Logger) middleware.LogFormatter {
	fi, err := os.Stdout.Stat()
	if err == nil {
		m := os.ModeDevice | os.ModeCharDevice
		isTTY = fi.Mode()&m == m
	}

	color := true
	if runtime.GOOS == "windows" {
		color = false
	}
	return &zeroLogFormatter{
		logger:  l,
		NoColor: !color,
	}
}

func (z *zeroLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {

	switch {
	case status < 200:
		colorWrite(z.buf, z.useColor, bBlue, "%03d", status)
	case status < 300:
		colorWrite(z.buf, z.useColor, bGreen, "%03d", status)
	case status < 400:
		colorWrite(z.buf, z.useColor, bCyan, "%03d", status)
	case status < 500:
		colorWrite(z.buf, z.useColor, bYellow, "%03d", status)
	default:
		colorWrite(z.buf, z.useColor, bRed, "%03d", status)
	}

	colorWrite(z.buf, z.useColor, bBlue, " %dB", bytes)

	z.buf.WriteString(" in ")
	if elapsed < 500*time.Millisecond {
		colorWrite(z.buf, z.useColor, nGreen, "%s", elapsed)
	} else if elapsed < 5*time.Second {
		colorWrite(z.buf, z.useColor, nYellow, "%s", elapsed)
	} else {
		colorWrite(z.buf, z.useColor, nRed, "%s", elapsed)
	}

	z.log.Info().Msg(z.buf.String())
}

func (z *zeroLogEntry) Panic(v interface{}, stack []byte) {
	z.log.Info().Msgf("request failed: %+v", v)
}

func (l zeroLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	useColor := !l.NoColor

	entry := &zeroLogEntry{
		log:      l.logger,
		request:  r,
		buf:      &bytes.Buffer{},
		useColor: useColor,
	}

	colorWrite(entry.buf, useColor, nCyan, "\"")
	colorWrite(entry.buf, useColor, bMagenta, "%s ", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	colorWrite(entry.buf, useColor, nCyan, "%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)

	entry.buf.WriteString("from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteString(" - ")

	return entry
}

func colorWrite(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if isTTY && useColor {
		w.Write(color)
	}
	fmt.Fprintf(w, s, args...)
	if isTTY && useColor {
		w.Write(reset)
	}
}
