package main

import (
	"aio/writer"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

func createLogHandler(w io.Writer, lvl slog.Level, typ string) (slog.Handler, error) {
	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: lvl <= slog.LevelInfo,
	}

	switch strings.ToLower(typ) {
	case "text":
		return slog.NewTextHandler(w, opts), nil
	case "json":
		return slog.NewJSONHandler(w, opts), nil
	default:
		return nil, fmt.Errorf("unexpected log format %q", typ)
	}
}

func main() {
	handler, err := createLogHandler(writer.NewWriter(os.Stdout, 10000000), slog.LevelInfo, "json")
	if err != nil {
		return
	}
	l := slog.New(handler)
	for {
		l.Info(`{"some": "message"}`)
	}
}
