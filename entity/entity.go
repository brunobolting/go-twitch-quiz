package entity

import "github.com/google/uuid"

type ID = string

func NewID() ID {
	return ID(uuid.New().String())
}

func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id.String()), err
}
