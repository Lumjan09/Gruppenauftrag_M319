// Package logging stellt strukturiertes Logging mit vier Stufen
// (Debug/Info/Warn/Error) und täglicher Datei-Rotation bereit.
//
// Es werden ausschliesslich Pakete der Standardbibliothek verwendet
// (log/slog). Die Rotation erfolgt über einen eigenen Writer, sodass kein
// externes Abhängigkeitspaket (z. B. lumberjack) benötigt wird.
package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var logger *slog.Logger

// Init initialisiert den globalen Logger.
//
// level ist eine der Stufen "debug", "info", "warn", "error".
// path ist der Pfad zur Logdatei (z. B. ./logs/battle.log). Der Ordner wird
// bei Bedarf erstellt; ist die Datei nicht beschreibbar, wird ein Fehler
// zurückgegeben (der Aufrufer kann daraus eine Panic erzeugen).
func Init(level, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("logverzeichnis konnte nicht erstellt werden: %w", err)
	}
	rw, err := newRotatingWriter(path)
	if err != nil {
		return fmt.Errorf("logdatei nicht beschreibbar: %w", err)
	}
	out := io.MultiWriter(os.Stdout, rw)
	handler := slog.NewTextHandler(out, &slog.HandlerOptions{Level: parseLevel(level)})
	logger = slog.New(handler)
	return nil
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func ensure() *slog.Logger {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}

// Debug protokolliert eine Debug-Nachricht.
func Debug(msg string, args ...any) { ensure().Debug(msg, args...) }

// Info protokolliert eine Info-Nachricht.
func Info(msg string, args ...any) { ensure().Info(msg, args...) }

// Warn protokolliert eine Warn-Nachricht.
func Warn(msg string, args ...any) { ensure().Warn(msg, args...) }

// Error protokolliert eine Error-Nachricht.
func Error(msg string, args ...any) { ensure().Error(msg, args...) }

// rotatingWriter schreibt in eine Logdatei und rotiert sie täglich.
// Beim Datumswechsel wird die alte Datei mit Datums-Suffix umbenannt.
type rotatingWriter struct {
	mu       sync.Mutex
	basePath string
	day      string
	file     *os.File
}

func newRotatingWriter(path string) (*rotatingWriter, error) {
	w := &rotatingWriter{basePath: path}
	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *rotatingWriter) rotateIfNeeded() error {
	today := time.Now().Format("2006-01-02")
	if today == w.day && w.file != nil {
		return nil
	}
	if w.file != nil {
		_ = w.file.Close()
		archived := fmt.Sprintf("%s.%s", w.basePath, w.day)
		_ = os.Rename(w.basePath, archived)
	}
	f, err := os.OpenFile(w.basePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	w.file = f
	w.day = today
	return nil
}

// Write implementiert io.Writer mit täglicher Rotation.
func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.rotateIfNeeded(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}
