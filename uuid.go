package main

import (
	"github.com/gofrs/uuid"
)

func UUID() string {
	id, _ := uuid.NewV4()
	return id.String()
}

func ValidUUIDString(sid string) bool {
	_, err := uuid.FromString(sid)
	return err == nil
}
