package sshutil

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// A KeyPair can generate a ssh.Signer for authenticating with an SSH server.
type KeyPair interface {
	Signer() (ssh.Signer, error)
}

// Key generates a KeyPair from file. This method should be used if the key is
// not password-encrypted.
func Key(file string) KeyPair {
	return keyPair{File: file}
}

// EncryptedKey generates a KeyPair from a password-encrypted key file.
func EncryptedKey(file, password string) KeyPair {
	return keyPair{
		File:     file,
		Password: []byte(password),
	}
}

// UserKey generates a KeyPair from a file path relative to the current user's
// SSH directory (ie, ~/.ssh). This method should be used if the key is not
// password-encrypted.
func UserKey(file string) KeyPair { return Key(sshDir(file)) }

// EncryptedUserKey generates a KeyPair from a password-encrypted file relative
// to the current user's SSH directory (ie, ~/.ssh).
func EncryptedUserKey(file, password string) KeyPair {
	return EncryptedKey(sshDir(file), password)
}

type keyPair struct {
	File     string
	Password []byte
}

func (kp keyPair) Signer() (s ssh.Signer, err error) {
	b, err := ioutil.ReadFile(kp.File)
	if err != nil {
		return nil, ke(kp, err)
	}

	p, _ := pem.Decode(b)
	if p == nil {
		return nil, ke(kp, errors.New("no PEM block found"))
	}

	if x509.IsEncryptedPEMBlock(p) {
		s, err = kp.parseEncryptedKey(p)
	} else {
		s, err = ssh.ParsePrivateKey(b)
	}

	if err != nil {
		return nil, ke(kp, err)
	}

	return
}

func (kp keyPair) parseEncryptedKey(p *pem.Block) (ssh.Signer, error) {
	if kp.Password == nil {
		return nil, errors.New("no password provided for encrypted key")
	}

	der, err := x509.DecryptPEMBlock(p, kp.Password)
	if err != nil {
		return nil, err
	}

	var s crypto.PrivateKey
	switch p.Type {
	case "RSA PRIVATE KEY":
		s, err = x509.ParsePKCS1PrivateKey(der)
	case "EC PRIVATE KEY":
		s, err = x509.ParseECPrivateKey(der)
	case "DSA PRIVATE KEY":
		s, err = ssh.ParseDSAPrivateKey(der)
	default:
		s, err = nil, fmt.Errorf("unsupported private key type %q", p.Type)
	}

	if err != nil {
		return nil, err
	}

	return ssh.NewSignerFromKey(s)
}

func sshDir(file string) string {
	u, err := user.Current()

	if err != nil {
		l("unable to find current user: %v", err)
		return file
	}

	return filepath.Clean(filepath.Join(u.HomeDir, ".ssh", file))
}

// A KeyError is returned from any of these key-related methods, composing
// around the source error, and including the filename for the key that
// generated the error.
type KeyError struct {
	Err  error
	File string
}

func ke(kp keyPair, err error) KeyError {
	return KeyError{
		Err:  err,
		File: kp.File,
	}
}

// Error satisfies the error interface
func (e KeyError) Error() string {
	return fmt.Sprintf("sshutil: key error %q: %v", e.File, e.Err)
}

var _ KeyPair = keyPair{}
var _ error = KeyError{}
