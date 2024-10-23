package port

type Detail struct {
	Key        string
	Value      interface{}
	DetailType int
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Alarm(msg string, args ...interface{})
	WithError(err error) Logger
	WithDetails(detail ...Detail) Logger
}
