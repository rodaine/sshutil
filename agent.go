package sshutil

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh/agent"
)

// AgentWithSocket connects to an SSH Agent socket s. An error is returned if
// the socket cannot be dialed. To use the default Agent, DefaultAgent can be
// used instead.
func AgentWithSocket(s string) (agent.Agent, error) {
	sock, err := net.Dial("unix", s)
	if err != nil {
		return nil, err
	}

	return agent.NewClient(sock), nil
}

const sshAgentSocket = "SSH_AUTH_SOCK"

// StdAgent returns the default SSH Agent taken from the environment
// (SSH_AUTH_SOCK). An error is returned if the socket cannot be dialed or the
// agent isn't running.
func StdAgent() (agent.Agent, error) {
	sock := os.Getenv(sshAgentSocket)
	return AgentWithSocket(sock)
}
