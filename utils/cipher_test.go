package utils

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestPBKDF2Base64(t *testing.T) {
	str := base64.StdEncoding.EncodeToString([]byte("ebPDWd3PtkybuV9"))
	fmt.Println(str, "len:", len(str))
}

func TestUnix2Base62(t *testing.T) {
	fmt.Println(Unix2Base62())
}
