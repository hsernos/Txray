package protocols

import (
	"encoding/base64"
	"strings"
)

func base64Encode(str string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(str))
}

func base64EncodeWithEq(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) string {
	i := len(str) % 4
	switch i {
	case 1:
		str = str[:len(str)-1]
	case 2:
		str += "=="
	case 3:
		str += "="
	}
	var data []byte
	var err error
	if strings.Contains(str, "-") || strings.Contains(str, "_") {
		data, err = base64.URLEncoding.DecodeString(str)
	} else {
		data, err = base64.StdEncoding.DecodeString(str)
	}
	if err != nil {
		return ""
	}
	return string(data)
}
