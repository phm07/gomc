package profile

import "github.com/google/uuid"

type Property struct {
	Name      string
	Value     string
	Signature string
}

type Profile struct {
	Id         uuid.UUID
	Name       string
	Properties []Property
}
