package profile

import (
	"crypto/md5"
	"github.com/google/uuid"
)

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

func OfflineProfile(name string) *Profile {
	return &Profile{
		Id:   offlineUUIDFromString(name),
		Name: name,
	}
}

func offlineUUIDFromString(name string) uuid.UUID {
	h := md5.New()
	h.Write([]byte("OfflinePlayer:" + name))
	s := h.Sum(nil)
	var id uuid.UUID
	copy(id[:], s)
	id[6] = (id[6] & 0x0f) | uint8(3<<4)
	id[8] = (id[8] & 0x3f) | 0x80
	return id
}
