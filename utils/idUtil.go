package utils

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GenerateID() (string, error) {
	idObj, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	idStr := idObj.String()
	idStr = strings.ReplaceAll(idStr, "-", "")
	return idStr, nil
}
