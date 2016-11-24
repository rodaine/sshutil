package sshutil

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyError_Error(t *testing.T) {
	t.Parallel()

	e := KeyError{
		Err:  errors.New("bar"),
		File: "foo",
	}

	assert.Equal(t, `sshutil: key error "foo": bar`, e.Error())
}

func TestKey(t *testing.T) {
	t.Parallel()

	kp := Key("foo").(keyPair)
	assert.Equal(t, "foo", kp.File)
	assert.Empty(t, kp.Password)
}

func TestEncryptedKey(t *testing.T) {
	t.Parallel()

	kp := EncryptedKey("foo", "bar").(keyPair)
	assert.Equal(t, "foo", kp.File)
	assert.Equal(t, []byte("bar"), kp.Password)
}

func TestSSHDir(t *testing.T) {
	t.Parallel()

	u, _ := user.Current()
	expected := filepath.Join(u.HomeDir, ".ssh", "foo")
	assert.Equal(t, expected, sshDir(filepath.Join(".", "foo")))
}

func TestUserKey(t *testing.T) {
	t.Parallel()

	u, _ := user.Current()
	expected := filepath.Join(u.HomeDir, ".ssh", "foo")

	kp := UserKey("foo").(keyPair)
	assert.Equal(t, expected, kp.File)
	assert.Empty(t, kp.Password)
}

func TestEncryptedUserKey(t *testing.T) {
	t.Parallel()

	u, _ := user.Current()
	expected := filepath.Join(u.HomeDir, ".ssh", "foo")

	kp := EncryptedUserKey("foo", "bar").(keyPair)
	assert.Equal(t, expected, kp.File)
	assert.Equal(t, []byte("bar"), kp.Password)
}

func TestKeyPair_Signer(t *testing.T) {
	t.Parallel()

	keyDir := filepath.Join("testdata", "keys")
	pw, err := ioutil.ReadFile(filepath.Join(keyDir, "password"))
	assert.NoError(t, err)
	pw = bytes.TrimSpace(pw)

	tests := []struct {
		File      string
		Encrypted bool
		Error     bool
		Desc      string
	}{
		{File: "rsa", Desc: "unencrypted RSA"},
		{File: "dsa", Desc: "unencrypted DSA"},
		{File: "ecdsa", Desc: "unencrypted ECDSA"},
		{File: "openssh", Desc: "unencrypted OPENSSH (ed25519)"},

		{File: "rsa_enc", Encrypted: true, Desc: "encrypted RSA"},
		{File: "dsa_enc", Encrypted: true, Desc: "encrypted DSA"},
		{File: "ecdsa_enc", Encrypted: true, Desc: "encrypted ECDSA"},
		{File: "openssh_enc", Encrypted: true, Error: true, Desc: "encrypted OPENSSH (ed25519 - unsupported)"},

		{File: "rsa_enc", Error: true, Desc: "encrypted RSA - no password"},
		{File: "rsa_enc_unknown", Error: true, Encrypted: true, Desc: "encrypted RSA - wrong password"},

		{File: "nonexistant", Error: true, Desc: "nonexistant file"},
		{File: "malformed", Error: true, Desc: "malformed file"},
		{File: "malformed_enc", Error: true, Encrypted: true, Desc: "encrypted PEM, but malformed/unsupported DEK info"},
		{File: "unsupported", Error: true, Desc: "unsupported key type"},
		{File: "unsupported_enc", Error: true, Encrypted: true, Desc: "unsupported encrypted key type"},
	}

	for _, test := range tests {
		fn := filepath.Join(keyDir, test.File)
		var kp KeyPair

		if test.Encrypted {
			kp = EncryptedKey(fn, string(pw))
		} else {
			kp = Key(fn)
		}

		s, err := kp.Signer()

		if test.Error {
			assert.Error(t, err, test.Desc)
			continue
		}
		assert.NoError(t, err, test.Desc)

		b, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", fn))
		assert.NoError(t, err, test.Desc)

		parts := strings.Split(string(b), " ")
		assert.True(t, len(parts) >= 2, test.Desc)

		pub := s.PublicKey()
		assert.Equal(t, parts[0], pub.Type(), test.Desc)
		assert.Equal(t, parts[1], base64.StdEncoding.EncodeToString(s.PublicKey().Marshal()), test.Desc)
	}
}
