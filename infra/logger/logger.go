package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func NewLogger(logDir string) (*zerolog.Logger, error) {
	// to print stack trace
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// log writers
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	fileWriter, err := setupLogFileWriter(logDir)
	if err != nil {
		return nil, err
	}

	logWriters := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	logger := zerolog.New(logWriters).
		With().
		Caller().
		Timestamp().
		Logger()

	// TODO: Add logging file
	log.Info().Msg("Hello World!")

	return &logger, nil
}

func setupLogFileWriter(logDir string) (*os.File, error) {
	if logDir == "" {
		logDir = "logs"
	}

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	logFilePath := logDir + "/" + time.Now().Format("2006-01-02") + ".log"
	fileWriter, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return fileWriter, nil
}
