package utils

import (
	"fmt"
	"log/slog"
	"os"
)

// SlogAdapter is used to create a logger for asynq
type SlogAdapter struct {
	Logger *slog.Logger
}

func (s *SlogAdapter) Debug(args ...interface{}) { s.Logger.Debug(fmt.Sprint(args...)) }
func (s *SlogAdapter) Info(args ...interface{})  { s.Logger.Info(fmt.Sprint(args...)) }
func (s *SlogAdapter) Warn(args ...interface{})  { s.Logger.Warn(fmt.Sprint(args...)) }
func (s *SlogAdapter) Error(args ...interface{}) { s.Logger.Error(fmt.Sprint(args...)) }
func (s *SlogAdapter) Fatal(args ...interface{}) {
	s.Logger.Error(fmt.Sprint(args...))
	os.Exit(1)
}
