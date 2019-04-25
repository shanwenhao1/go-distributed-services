package utils

import "github.com/satori/go.uuid"

func NewUuid() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}
