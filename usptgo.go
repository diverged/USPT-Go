package usptgo

import (
	"github.com/diverged/uspt-go/internal"
	"github.com/diverged/uspt-go/types"
)

// USPTGo accepts a USPTGoConfig and returns one USPTGoDoc channel and one USPTGoError channel
func USPTGo(cfg *types.USPTGoConfig) (<-chan *types.USPTGoDoc, <-chan error, error) {

	// Defaults to a no-op logger which does nothing with log messages
	if cfg.Logger == nil {
		cfg.Logger = noOpLogger{}
	}

	docChan, errChan, err := internal.Dispatcher(cfg)
	if err != nil {
		return nil, nil, err
	}

	return docChan, errChan, nil

}

type noOpLogger struct{}

func (l noOpLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (l noOpLogger) Info(msg string, keysAndValues ...interface{})  {}
func (l noOpLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (l noOpLogger) Error(msg string, keysAndValues ...interface{}) {}
