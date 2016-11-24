package sshutil

import "testing"

import "net"
import "github.com/stretchr/testify/assert"
import "os"

func TestAgentWithSocket(t *testing.T) {
	t.Parallel()

	sock := "/tmp/sshutil_TestAgentWithSocket.sock"
	defer os.Remove(sock)

	l, err := net.Listen("unix", sock)
	assert.NoError(t, err)
	defer l.Close()

	agent, err := AgentWithSocket(sock)
	assert.NoError(t, err)
	assert.NotNil(t, agent)
}

func TestAgentWithSocket_BadSocket(t *testing.T) {
	t.Parallel()

	agent, err := AgentWithSocket("/tmp/sshutil_TestAgentWithSocket_BadSocket.sock")
	assert.Error(t, err)
	assert.Nil(t, agent)
}

func TestDefaultAgent(t *testing.T) {
	t.Parallel()

	sock := "/tmp/sshutil_TestDefaultAgent.sock"
	defer os.Remove(sock)

	prevSock := os.Getenv(sshAgentSocket)
	defer os.Setenv(sshAgentSocket, prevSock)
	os.Setenv(sshAgentSocket, sock)

	l, err := net.Listen("unix", sock)
	assert.NoError(t, err)
	defer l.Close()

	agent, err := StdAgent()
	assert.NoError(t, err)
	assert.NotNil(t, agent)
}
