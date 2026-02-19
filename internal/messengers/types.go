package messengers

import (
	"context"
	"encoding/json"
)

// EventType identifies the kind of event a Message carries.
type EventType string

const (
	EventFlowExecution EventType = "flow.execution"
)

// FlowExecutionEvent carries structured data about a flow execution state change.
type FlowExecutionEvent struct {
	FlowID    string `json:"flow_id"`
	FlowName  string `json:"flow_name"`
	ExecID    string `json:"exec_id"`
	Status    string `json:"status"`
	Error     string `json:"error"`
	Namespace string `json:"namespace"`
	RootURL   string `json:"-"`
}

// Message is the generic struct passed to messengers.
type Message struct {
	Event  EventType
	Data   any
	Config map[string]any
}

type Messenger interface {
	Send(ctx context.Context, message Message) error
	Close()
}

// GroupResolver resolves a group name to a list of member email addresses.
type GroupResolver interface {
	ResolveGroupEmails(ctx context.Context, groupName string) ([]string, error)
}

// configStringSlice extracts a []string value from a config map.
func configStringSlice(cfg map[string]any, key string) []string {
	v, ok := cfg[key]
	if !ok {
		return nil
	}
	b, _ := json.Marshal(v)
	var s []string
	json.Unmarshal(b, &s)
	return s
}
