package feature

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidSchemaURN        = errors.New("invalid schema URN")
	ErrInvalidSchemaVersionURN = errors.New("invalid schema version URN")
)

func getSchemaVersionFromURN(schemaURN string) (int, error) {
	parts := strings.Split(schemaURN, ":")
	if len(parts) < 3 {
		return -1, ErrInvalidSchemaURN
	}
	schemaID := parts[2]

	schemaParts := strings.Split(schemaID, "/")
	if len(schemaParts) < 2 {
		return -1, ErrInvalidSchemaVersionURN
	}

	schemaVersion, err := strconv.Atoi(schemaParts[1])
	if err != nil {
		return -1, ErrInvalidSchemaVersionURN
	}

	return schemaVersion, nil
}
