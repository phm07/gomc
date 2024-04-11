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

//packet:2:3
type ServerboundLoginAck struct{}

//packet:3:0
type ServerboundConfigurationClientInformation struct {
	Language       types.String
	ViewDistance   types.Byte
	ChatMode       types.VarInt
	ChatColors     types.Byte
	SkinParts      types.Byte
	MainHand       types.VarInt
	TextFiltering  types.Byte
	ServerListings types.Byte
}

//packet:3:1
type ServerboundConfigurationPluginMessage struct {
	Channel types.String
	Data    types.Data
}

//packet:3:2
type ServerboundConfigurationFinishAck struct{}
