package main

import "io"

type MockEncrypter struct{}

func (m MockEncrypter) EncryptPayloadToWriter(_ []byte, _ io.Writer) error {
	return nil
}

func (m MockEncrypter) DecryptPayloadFromReader(_ io.Reader) ([]byte, error) {
	return []byte("hello world!"), nil
}
