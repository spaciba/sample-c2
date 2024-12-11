package listeners

type ListenerState int

type ListenerType int

const (
	STOPPED ListenerState = iota
	STOPPING
	STARTING
	STARTED

	TCP ListenerType = iota
	HTTP
)

type Listener interface {
	Init( listener_address string) (listenerID string, message string, error_code int32) 
}



