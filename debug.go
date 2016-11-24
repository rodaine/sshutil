package sshutil

import "sync"

var debug = &debugger{}

// Logger defines the interface for the debug logger. This can be set defined
// via SetDebug.
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// SetDebug activates the sshutil debugger with l. By default, no logger is
// specified. This method is thread-safe, but the passed in Logger must also be
// (ie, log.Logger).
func SetDebug(l Logger) { debug.setLog(l) }

func l(format string, v ...interface{}) { debug.log(format, v...) }

type debugger struct {
	sync.RWMutex
	l Logger
}

func (d *debugger) setLog(l Logger) {
	d.Lock()
	d.l = l
	d.Unlock()
}

func (d *debugger) log(format string, v ...interface{}) {
	d.RLock()
	l := d.l
	d.RUnlock()

	if l == nil {
		return
	}

	if len(v) == 0 {
		l.Print(format)
		return
	}

	l.Printf(format, v...)
}
