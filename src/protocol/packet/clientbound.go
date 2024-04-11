package packet

import (
	"gomc/src/protocol/types"
)

//go:generate go run generation/gen.go -- $GOFILE

//packet:1:0
type ClientboundStatusResponse struct {
	Json types.String
}

//packet:1:1
type ClientboundStatusPong struct {
	Payload types.Long
}

//packet:2:0
type ClientboundLoginDisconnect struct {
	Reason types.String
}

//packet:2:1
type ClientboundLoginEncryptionRequest struct {
	ServerID    types.String
	PublicKey   types.ByteBuf
	VerifyToken types.ByteBuf
}

//packet:2:2
type ClientboundLoginSuccess struct {
	UUID     types.UUID
	Username types.String
	Zero     types.VarInt
}

//packet:3:2
type ClientboundConfigurationFinish struct{}

//packet:3:5
type ClientboundConfigurationRegistryData struct {
	RegistryDataNBT types.Data
}
