package sshutil

import (
	"io/ioutil"
	"log"
	"sync"
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestDebugger_SetLog(t *testing.T) {
	t.Parallel()

	d := &debugger{}
	assert.Nil(t, d.l)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		d.setLog(log.New(ioutil.Discard, "fizz", 0))
		wg.Done()
	}()

	go func() {
		d.setLog(log.New(ioutil.Discard, "buzz", 0))
		wg.Done()
	}()

	wg.Wait()

	assert.NotNil(t, d.l)
}

func TestDebugger_Log(t *testing.T) {
	t.Parallel()

	d := &debugger{}

	assert.NotPanics(t, func() { d.log("sshutil") })

	wg := &sync.WaitGroup{}
	wg.Add(1)

	buf := new(bytes.Buffer)

	go func() {
		d.setLog(log.New(buf, "", 0))
		d.log("fizz %s", "buzz")
		wg.Done()
	}()

	assert.NotPanics(t, func() { d.log("foo") })

	wg.Wait()
	assert.Contains(t, buf.String(), "fizz buzz\n")

	d.log("bar")
	assert.Contains(t, buf.String(), "bar")
}

func TestL(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() { l("foo") })
}

func TestSetDebug(t *testing.T) {
	t.Parallel()

	l := log.New(ioutil.Discard, "", 0)
	SetDebug(l)
	assert.Equal(t, l, debug.l)
}
