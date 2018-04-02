package crypto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
)

// #cgo CFLAGS: -masm=intel
// #include "native.h"
//import "C"

var (
	keyTable = make([]byte, 0)
)

func ImportKeyTable(path string) error {
	if len(path) > 1 {
		if buffer, err := ioutil.ReadFile(path); err != nil {
			return err
		} else {
			keyTable = buffer
			return nil
		}
	}

	return errors.New("EmptyPath")
}

func Encrypt(buffer []byte) []byte {

	return nil
}

func Decrypt(buffer []byte) []byte {
	if len(buffer) > 4 && len(keyTable) > 1 {
		var keyIdentifier uint8
		var size uint16
		var encryptedBuffer []byte
		var sender uint8

		reader := bytes.NewReader(buffer)
		binary.Read(reader, binary.LittleEndian, &keyIdentifier)
		binary.Read(reader, binary.LittleEndian, &size)
		binary.Read(reader, binary.LittleEndian, &sender)
		encryptedBuffer = buffer[5:len(buffer)]
		outBuff := make([]byte, len(buffer))

		copy(outBuff[0:5], buffer[0:5])
		for i := uint16(0); i < uint16(len(encryptedBuffer)); i++ {
			outBuff[5+i] = encryptedBuffer[i] ^ keyTable[4*uint16(keyIdentifier)-3*(i/3)+i]
		}

		return outBuff
	}
	return nil
}
