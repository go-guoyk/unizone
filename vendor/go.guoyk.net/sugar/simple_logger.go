package sugar

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type simpleLogger struct {
	debug bool
	lock  sync.Locker
}

func (c *simpleLogger) Log(level string, out *os.File, message string, items ...interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	extras := make([]string, 0, len(items))
	for _, item := range items {
		extras = append(extras, fmt.Sprintf("%v", item))
	}
	sb := &strings.Builder{}
	sb.WriteString(time.Now().Format(time.RFC3339))
	sb.WriteString(" [")
	sb.WriteString(level)
	sb.WriteString("] ")
	sb.WriteString(message)
	sb.WriteRune(' ')
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("%v", item))
		if i < len(items)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteRune('\n')
	_, _ = out.WriteString(sb.String())
	_ = out.Sync()
}

func (c *simpleLogger) Debug(message string, items ...interface{}) {
	if !c.debug {
		return
	}
	c.Log("debug", os.Stdout, message, items...)
}

func (c *simpleLogger) Info(message string, items ...interface{}) {
	c.Log("info", os.Stdout, message, items...)
}

func (c *simpleLogger) Warn(message string, items ...interface{}) {
	c.Log("warn", os.Stdout, message, items...)
}

func (c *simpleLogger) Error(message string, items ...interface{}) {
	c.Log("error", os.Stderr, message, items...)
}

func (c *simpleLogger) Panic(message string, items ...interface{}) {
	c.Log("panic", os.Stderr, message, items...)
	panic("panic from sugar.simpleLogger")
}

func (c *simpleLogger) Fatal(message string, items ...interface{}) {
	c.Log("fatal", os.Stderr, message, items...)
	os.Exit(1)
}

// NewSimpleLogger create a simple Logger, shall only use for library development
func NewSimpleLogger(debug bool) Logger {
	return &simpleLogger{debug: debug, lock: &sync.Mutex{}}
}
