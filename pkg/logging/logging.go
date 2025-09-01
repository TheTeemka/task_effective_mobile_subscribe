package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func SetSlog(stage string) {

	var l slog.Level
	var w io.Writer = os.Stdout

	switch stage {
	case slog.LevelDebug.String():
		l = slog.LevelDebug
	case slog.LevelWarn.String():
		l = slog.LevelWarn
	case slog.LevelInfo.String():
		l = slog.LevelInfo
	case slog.LevelError.String():
		l = slog.LevelError
	default:
		panic("Unknown stage")
	}

	h := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: l,
		// AddSource: stage == StageDev,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.File = filepath.Base(source.File)
				}
			}
			return a
		},
	})

	slog.SetDefault(slog.New(h))
}
