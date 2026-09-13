package utils

import (
	"context"
	"log/slog"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// 初始日志配置：多路输出+轮转
func initLogger() {
	// 控制台handler
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	// 文件 handler
	writer, err := rotatelogs.New(
		"./logs/app-%Y%m%d.log",
		rotatelogs.WithLinkName("./logs/app-latest.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic("Faild to initialize log rotator: " + err.Error())
	}

	fileHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	})

	logger := slog.New(NewMultihander(consoleHandler, fileHandler))
	slog.SetDefault(logger)
}

type multiHandler []slog.Handler

func NewMultihander(handlers ...slog.Handler) slog.Handler {
	return multiHandler(handlers)
}

// 实现enable方法，只要有一个handler 启用了这个级别，就返回true
func (h multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, l) {
			return true
		}
	}

	return false
}

// 实现handler方法，遍历所有hanlder,如果该handler启用了当前级别，则处理日志
func (h multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h {
		if handler.Enabled(ctx, r.Level) {
			// 注意：实际生产中可能需要 Clone Record 以防 handler 修改原始记录
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}

	return nil
}

// 实现 withAttrs 和 withGroup 保证实现接口
func (h multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h))
	for i, handler := range h {
		handlers[i] = handler.WithAttrs(attrs)
	}

	return multiHandler(handlers)
}

func (h multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h))
	for i, handler := range h {
		handlers[i] = handler.WithGroup(name)
	}

	return multiHandler(handlers)
}
