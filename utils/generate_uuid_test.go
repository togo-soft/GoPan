package utils

import (
	"fmt"
	uuid "github.com/gofrs/uuid"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	fmt.Println(u.String())

}
