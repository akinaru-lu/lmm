package log

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type kafkaWriter struct{}

func newKafkaSyncWriter() zapcore.WriteSyncer {
	w := zapcore.AddSync(new(kafkaWriter))
	return zapcore.Lock(w)
}

func (w *kafkaWriter) Write(b []byte) (int, error) {
	return fmt.Println("TODO: send logger data to kafka")
}

// Init initializes logger for std log and zap global logger
// Expected to be called in main() first
func Init() func() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "msg",
		NameKey:      "logger",
		TimeKey:      "ts",
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		LineEnding:   zapcore.DefaultLineEnding,
	}

	kafkaEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	globalEnabler := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return true
	})

	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, newKafkaSyncWriter(), globalEnabler),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), globalEnabler),
	)

	logger := zap.New(core).Named("logger")

	undo := zap.ReplaceGlobals(logger)

	return func() {
		logger.Sync()
		undo()
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}
