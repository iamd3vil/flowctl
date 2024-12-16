package models

type LogMessage struct {
	Message    string
	Checkpoint string
	Err        error
}
