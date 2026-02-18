package messengers

import "context"

// EventType identifies the kind of event a Message carries.
type EventType string

const (
	EventFlowExecution EventType = "flow.execution"
)

// FlowExecutionEvent carries structured data about a flow execution state change.
type FlowExecutionEvent struct {
	FlowID    string
	FlowName  string
	ExecID    string
	Status    string // "completed", "errored", "cancelled", "pending_approval"
	Error     string
	Namespace string
	RootURL   string
}

// Message is the generic envelope passed to messengers.
// Data is `any` so different event types can use different structs.
type Message struct {
	Event  EventType      // e.g. EventFlowExecution
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
	s, _ := v.([]string)
	return s
}
