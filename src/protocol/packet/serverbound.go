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
	ChatColors     types.Boolean
	SkinParts      types.Byte
	MainHand       types.VarInt
	TextFiltering  types.Boolean
	ServerListings types.Boolean
}

//packet:3:1
type ServerboundConfigurationPluginMessage struct {
	Channel types.String
	Data    types.Data
}

//packet:3:2
type ServerboundConfigurationFinishAck struct{}

//packet:4:0
type ServerboundPlayConfirmTeleport struct {
	TeleportID types.VarInt
}

//packet:4:5
type ServerboundPlayChatMessage struct {
	Message types.String
	Ignored types.Data
}

//packet:4:15
type ServerboundPlayKeepAlive struct {
	KeepAliveId types.Long
}

//packet:4:17
type ServerboundPlayUpdatePosition struct {
	X, Y, Z  types.Double
	OnGround types.Boolean
}

//packet:4:18
type ServerboundPlayUpdatePositionAndRotation struct {
	X, Y, Z    types.Double
	Yaw, Pitch types.Float
	OnGround   types.Boolean
}

//packet:4:19
type ServerboundPlayUpdateRotation struct {
	Yaw, Pitch types.Float
	OnGround   types.Boolean
}
