package logium

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type AlignedTextOptions struct {
	Level      slog.Level
	TimeFormat string

	MsgWidth int

	Colors   bool
	SortKeys bool
}

type TextHandler struct {
	out  io.Writer
	opts AlignedTextOptions

	attrs []slog.Attr
	group string

	mu *sync.Mutex
}

func NewAlignedTextHandler(out io.Writer, opts AlignedTextOptions) slog.Handler {
	if opts.TimeFormat == "" {
		opts.TimeFormat = "2006-01-02 15:04:05"
	}
	if opts.MsgWidth <= 0 {
		opts.MsgWidth = 55
	}

	if opts.Colors {
		if shouldDisableColors(out) {
			opts.Colors = false
		}
	}

	return &TextHandler{
		out:  out,
		opts: opts,
		mu:   &sync.Mutex{},
	}
}

func (h *TextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

func (h *TextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	cp := *h
	cp.attrs = append(append([]slog.Attr(nil), h.attrs...), attrs...)
	return &cp
}

func (h *TextHandler) WithGroup(name string) slog.Handler {
	cp := *h
	if cp.group == "" {
		cp.group = name
	} else if name != "" {
		cp.group = cp.group + "." + name
	}
	return &cp
}

func (h *TextHandler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time
	if ts.IsZero() {
		ts = time.Now()
	}

	levelText := level4(r.Level)
	if h.opts.Colors {
		levelText = colorize(levelText, r.Level)
	}

	head := levelText + "[" + ts.Format(h.opts.TimeFormat) + "] "

	msg := padRightRunes(r.Message, h.opts.MsgWidth)

	pairs := make([]kv, 0, 16)

	for _, a := range h.attrs {
		h.collectAttr(&pairs, a)
	}
	r.Attrs(func(a slog.Attr) bool {
		h.collectAttr(&pairs, a)
		return true
	})

	if h.opts.SortKeys {
		sort.Slice(pairs, func(i, j int) bool { return pairs[i].k < pairs[j].k })
	}

	var b strings.Builder
	b.Grow(256)

	b.WriteString(head)
	b.WriteString(msg)

	for _, p := range pairs {
		b.WriteByte(' ')
		if h.opts.Colors {
			b.WriteString(colorize(p.k, r.Level))
		} else {
			b.WriteString(p.k)
		}
		b.WriteByte('=')
		b.WriteString(formatValue(p.v))
	}

	b.WriteByte('\n')

	h.mu.Lock()
	_, err := io.WriteString(h.out, b.String())
	h.mu.Unlock()
	return err
}

type kv struct {
	k string
	v slog.Value
}

func (h *TextHandler) collectAttr(dst *[]kv, a slog.Attr) {
	if a.Equal(slog.Attr{}) {
		return
	}

	key := a.Key
	if h.group != "" {
		key = h.group + "." + key
	}

	v := a.Value.Resolve()

	if key == "error" && v.Kind() == slog.KindAny && v.Any() == nil {
		return
	}

	if v.Kind() == slog.KindGroup {
		for _, ga := range v.Group() {
			if ga.Equal(slog.Attr{}) {
				continue
			}
			gk := key + "." + ga.Key
			*dst = append(*dst, kv{k: gk, v: ga.Value.Resolve()})
		}
		return
	}

	*dst = append(*dst, kv{k: key, v: v})
}

func level4(l slog.Level) string {
	switch {
	case l <= slog.LevelDebug:
		return "DEBU"
	case l == slog.LevelInfo:
		return "INFO"
	case l == slog.LevelWarn:
		return "WARN"
	default:
		return "ERRO"
	}
}

func padRightRunes(s string, width int) string {
	if width <= 0 {
		return s
	}
	n := utf8.RuneCountInString(s)
	if n >= width {
		return s
	}
	return s + strings.Repeat(" ", width-n)
}

func formatValue(v slog.Value) string {
	v = v.Resolve()

	switch v.Kind() {
	case slog.KindString:
		s := v.String()
		if needsQuotes(s) {
			return strconv.Quote(s)
		}
		return s
	case slog.KindBool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'f', -1, 64)
	case slog.KindTime:
		return v.Time().Format(time.RFC3339Nano)
	case slog.KindDuration:
		return v.Duration().String()
	default:
		return fmt.Sprint(v.Any())
	}
}

func needsQuotes(s string) bool {
	for _, r := range s {
		switch r {
		case ' ', '\t', '\n', '"', '=':
			return true
		}
	}
	return false
}

func colorize(text string, lvl slog.Level) string {
	const (
		reset = "\x1b[0m"
		blue  = "\x1b[34m"
		green = "\x1b[92m"
		yell  = "\x1b[33m"
		red   = "\x1b[31m"
	)

	var c string
	switch {
	case lvl <= slog.LevelDebug:
		c = blue
	case lvl == slog.LevelInfo:
		c = green
	case lvl == slog.LevelWarn:
		c = yell
	default:
		c = red
	}

	return c + text + reset
}

func shouldDisableColors(out io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return true
	}

	f, ok := out.(*os.File)
	if !ok {
		return true
	}

	st, err := f.Stat()
	if err != nil {
		return true
	}

	return (st.Mode() & os.ModeCharDevice) == 0
}
