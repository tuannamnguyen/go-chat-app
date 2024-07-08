package main

import (
	"bytes"
	"errors"
)

type message struct {
	bytes  []byte
	author *user
}

func (m message) prepareMsg() ([]byte, error) {
	buffer := bytes.NewBufferString(m.author.name + ": ")
	bufLen, _ := buffer.Write(m.bytes)

	if bufLen != len(m.bytes) {
		return []byte{}, errors.New("error creating message")
	}

	return buffer.Bytes(), nil
}
