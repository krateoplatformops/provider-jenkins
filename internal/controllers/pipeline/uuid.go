package pipeline

import (
	"strings"

	"github.com/google/uuid"
)

func fillWithUUID(str string, marker string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	res := strings.Replace(str, marker, id.String(), -1)
	return res, nil
}
