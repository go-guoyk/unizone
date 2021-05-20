package sugar

// Logger main interface, basically abstracted from zap.SugaredLogger
type Logger interface {
	Debug(message string, items ...interface{})
	Info(message string, items ...interface{})
	Warn(message string, items ...interface{})
	Error(message string, items ...interface{})
	Panic(message string, items ...interface{})
	Fatal(message string, items ...interface{})
}
