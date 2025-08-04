package base64Utils

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
)

func padBase64String(base64Str string) string {
	for len(base64Str)%4 != 0 {
		base64Str += "="
	}
	return base64Str
}

func DecodeBase64ToUInt64(base64Str string) (uint64, error) {
	data, err := base64.StdEncoding.DecodeString(padBase64String(base64Str))
	if err != nil {
		return 0, err
	}

	if len(data) < 8 {
		// Padding with 0s - little endian should receive the 0 at start
		data = append(make([]byte, 8-len(data)), data...)
	}
	if len(data) > 8 {
		return 0, errors.New("data is too long")
	}

	return binary.LittleEndian.Uint64(data), nil
}

func DecodeBase64ToStr(base64Str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(padBase64String(base64Str))
	if err != nil {
		return "", err
	}

	return string(data), nil
}
