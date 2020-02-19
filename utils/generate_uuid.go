package utils

import (
	"github.com/gofrs/uuid"
	"log"
)

func GenerateUUID() string {
	u, err := uuid.NewV4()
	if err != nil {
		log.Println("utils.GenerateUUID:", err)
	}
	return u.String()
}
