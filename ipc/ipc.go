package ipc

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
)

var socket net.Conn

// Choose the right directory to the ipc socket and return it
func GetIpcPath() string {
	variablesnames := []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

	for _, variablename := range variablesnames {
		path, exists := os.LookupEnv(variablename)

		if exists {
			return path
		}
	}

	return "/tmp"
}

func CloseSocket() error {
	if socket != nil {
		socket.Close()
		socket = nil
	}
	return nil
}

// Read the socket response
func Read() (string, error) {
	buf := make([]byte, 512)
	payloadlength, err := socket.Read(buf)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	for i := 8; i < payloadlength; i++ {
		err := buffer.WriteByte(buf[i])
		if err != nil {
			return "", err
		}
	}

	return buffer.String(), nil
}

// Send opcode and payload to the unix socket
func Send(opcode int, payload string) (string, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int32(opcode))
	if err != nil {
		return "", err
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(payload)))
	if err != nil {
		return "", err
	}

	buf.Write([]byte(payload))
	_, err = socket.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	return Read()
}
