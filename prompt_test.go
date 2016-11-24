package sshutil

import (
	"io/ioutil"
	"os"
	"testing"

	"io"

	"github.com/stretchr/testify/assert"
)

func closeAndDelete(f *os.File) {
	if f == nil {
		return
	}

	f.Close()
	os.Remove(f.Name())
}

func TestIOPrompt(t *testing.T) {
	t.Parallel()

	in, _ := ioutil.TempFile("", "in")
	defer closeAndDelete(in)
	out, _ := ioutil.TempFile("", "out")
	defer closeAndDelete(out)

	p := IOPrompt(in, out)
	assert.Equal(t, in, p.in)
	assert.Equal(t, out, p.out)
}

func TestIOPrompter_Prompt_Echo(t *testing.T) {
	t.Parallel()

	in, _ := ioutil.TempFile("", "in_echo")
	defer closeAndDelete(in)
	out, _ := ioutil.TempFile("", "out_echo")
	defer closeAndDelete(out)

	in.WriteString("bar\n")
	in.Seek(0, io.SeekStart)

	p := IOPrompt(in, out)
	res, err := p.Prompt("foo", true)
	assert.NoError(t, err)

	assert.Equal(t, "bar", res)
	out.Seek(0, io.SeekStart)

	actual := make([]byte, 4)
	n, err := out.Read(actual)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "foo", string(actual[:3]))
}
