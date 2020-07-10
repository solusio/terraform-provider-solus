package solus

type Logger interface {
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type NullLogger struct{}

func (NullLogger) Debugf(string, ...interface{}) {}
func (NullLogger) Errorf(string, ...interface{}) {}
