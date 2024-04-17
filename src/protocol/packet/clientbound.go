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

//packet:4:5
type ClientboundPlayAckBlockChange struct {
	SequenceID types.VarInt
}

//packet:4:9
type ClientboundPlayBlockUpdate struct {
	Location types.Position
	BlockID  types.VarInt
}

const (
	ClientboundPlayGameEventWaitForChunks types.Byte = 13
)

//packet:4:20
type ClientboundPlayGameEvent struct {
	Event types.Byte
	Value types.Float
}

//packet:4:24
type ClientboundPlayKeepAlive struct {
	KeepAliveId types.Long
}

//packet:4:25
type ClientboundPlayChunkData struct {
	ChunkX, ChunkZ       types.Int
	Heightmaps           types.Data
	Data                 types.ByteBuf
	NumBlockEntities     types.VarInt
	SkyLightMask         types.BitSet
	BlockLightMask       types.BitSet
	EmptySkyLightMask    types.BitSet
	EmptyBlockLightMask  types.BitSet
	SkyLight             types.Data
	BlockLightArrayCount types.Data
}

//packet:4:25
type ClientboundPlayChunkData2 struct {
	Data types.Data
}

//packet:4:29
type ClientboundPlayLogin struct {
	EntityID            types.Int
	IsHardcore          types.Boolean
	DimensionNames      []types.String
	MaxPlayers          types.VarInt
	ViewDistance        types.VarInt
	SimulationDistance  types.VarInt
	ReducedDebugInfo    types.Boolean
	EnableRespawnScreen types.Boolean
	LimitedCrafting     types.Boolean
	DimensionType       types.String
	DimensionName       types.String
	HashedSeed          types.Long
	GameMode            types.Byte
	PreviousGameMode    types.Byte
	IsDebug             types.Boolean
	IsFlat              types.Boolean
	HasDeathLocation    types.Boolean
	PortalCooldown      types.VarInt
}

const (
	PlayerCapabilityFlagInvulnerable types.Byte = 1 << iota
	PlayerCapabilityFlagFlying
	PlayerCapabilityFlagAllowFlying
	PlayerCapabilityFlagCreativeMode
)

//packet:4:36
type ClientboundPlayPlayerCapabilities struct {
	Flags        types.Byte
	FlyingSpeed  types.Float
	WalkingSpeed types.Float
}

//packet:4:3b
type ClientboundPlayPlayerInfoRemove struct {
	UUIDs []types.UUID
}

//packet:4:3e
type ClientboundPlaySynchronizePosition struct {
	X, Y, Z    types.Double
	Yaw, Pitch types.Float
	Flags      types.Byte
	TeleportID types.VarInt
}

//packet:4:52
type ClientboundPlaySetCenterChunk struct {
	ChunkX, ChunkZ types.VarInt
}

//packet:4:69
type ClientboundPlaySystemMessage struct {
	Content types.Data
	Overlay types.Boolean
}
