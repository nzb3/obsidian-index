package app

import (
	"log/slog"
	"os"

	"github.com/nzb3/obsidian-index/internal/indexator"
)

type config interface {
	GetVaultDir() string
	IsVerbose() bool
	IsDryRun() bool
	IsBackup() bool
	GetExcludeDirs() []string
}

type App struct {
	cfg       config
	indexator *indexator.Indexator
}

func New(cfg config) *App {
	app := &App{
		cfg: cfg,
	}

	app.initLogger()
	app.initIndexator()

	return app
}

func (app *App) initLogger() {
	var level slog.Level
	if app.cfg.IsVerbose() {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func (app *App) initIndexator() *indexator.Indexator {
	if app.indexator != nil {
		return app.indexator
	}

	app.indexator = indexator.NewIndexatorWithOptions(
		app.cfg.GetVaultDir(),
		app.cfg.IsDryRun(),
		app.cfg.IsBackup(),
		app.cfg.GetExcludeDirs(),
	)
	return app.indexator
}

func (app *App) Run() error {
	return app.indexator.Start()
}
