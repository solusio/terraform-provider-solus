package solus

// Logger represents a logger for the client.
type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// NullLogger represents a logger which don't log anything.
// Useful when you don't want to see any logs from the client or in tests.
type NullLogger struct{}

// Debugf logs message with debug level.
func (NullLogger) Debugf(string, ...interface{}) {}

// Errorf logs message with error level.
func (NullLogger) Errorf(string, ...interface{}) {}
