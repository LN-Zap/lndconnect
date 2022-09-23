package main

import "io"

// MockEncrypter is a mocked up structure implementing tor.NewOnionFile required encrypter methods
type MockEncrypter struct{}

// EncryptPayloadToWriter is implementing expected signature
func (m MockEncrypter) EncryptPayloadToWriter(_ []byte, _ io.Writer) error {
	return nil
}

// DecryptPayloadFromReader is implementing expected signature
func (m MockEncrypter) DecryptPayloadFromReader(_ io.Reader) ([]byte, error) {
	return []byte("hello world!"), nil
}
