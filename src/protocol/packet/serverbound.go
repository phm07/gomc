package packet

import (
	"gomc/src/protocol/types"
)

//go:generate go run generation/gen.go -- $GOFILE

//packet:0:0
type ServerboundHandshake struct {
	ProtocolVersion types.VarInt
	ServerAddress   types.String
	ServerPort      types.UShort
	NextState       types.VarInt
}

//packet:1:0
type ServerboundStatusRequest struct{}

//packet:1:1
type ServerboundStatusPing struct {
	Payload types.Long
}

//packet:2:0
type ServerboundLoginStart struct {
	Username types.String
	UUID     types.UUID
}

//packet:2:1
type ServerboundLoginEncryptionResponse struct {
	SharedSecret types.ByteBuf
	VerifyToken  types.ByteBuf
}
