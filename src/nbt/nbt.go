package nbt

import (
	"bytes"
	"fmt"
	"gomc/src/util"
	"math"
	"strings"
)

type Tag interface {
	GetName() string
	Marshal() []byte
	String() string
	ID() byte
}

type ByteTag struct {
	Name string
	Data byte
}

func (t *ByteTag) GetName() string {
	return t.Name
}

func (t *ByteTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.WriteByte(t.Data)
	return buf.Bytes()
}

func (t *ByteTag) String() string {
	return fmt.Sprintf("%s: %db", t.Name, t.Data)
}

func (*ByteTag) ID() byte {
	return 0x01
}

type ShortTag struct {
	Name string
	Data int16
}

func (t *ShortTag) GetName() string {
	return t.Name
}

func (t *ShortTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x02)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.Write(util.Int16ToBytes(t.Data))
	return buf.Bytes()
}

func (t *ShortTag) String() string {
	return fmt.Sprintf("%s: %ds", t.Name, t.Data)
}

func (*ShortTag) ID() byte {
	return 0x02
}

type IntTag struct {
	Name string
	Data int32
}

func (t *IntTag) GetName() string {
	return t.Name
}

func (t *IntTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x03)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.Write(util.Int32ToBytes(t.Data))
	return buf.Bytes()
}

func (t *IntTag) String() string {
	return fmt.Sprintf("%s: %d", t.Name, t.Data)
}

func (*IntTag) ID() byte {
	return 0x03
}

type LongTag struct {
	Name string
	Data int64
}

func (t *LongTag) GetName() string {
	return t.Name
}

func (t *LongTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x04)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.Write(util.Int64ToBytes(t.Data))
	return buf.Bytes()
}

func (t *LongTag) String() string {
	return fmt.Sprintf("%s: %dl", t.Name, t.Data)
}

func (*LongTag) ID() byte {
	return 0x04
}

type FloatTag struct {
	Name string
	Data float32
}

func (t *FloatTag) GetName() string {
	return t.Name
}

func (t *FloatTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x05)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.Write(util.Uint32ToBytes(math.Float32bits(t.Data)))
	return buf.Bytes()
}

func (t *FloatTag) String() string {
	return fmt.Sprintf("%s: %ff", t.Name, t.Data)
}

func (*FloatTag) ID() byte {
	return 0x05
}

type DoubleTag struct {
	Name string
	Data float64
}

func (t *DoubleTag) GetName() string {
	return t.Name
}

func (t *DoubleTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x06)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.Write(util.Uint64ToBytes(math.Float64bits(t.Data)))
	return buf.Bytes()
}

func (t *DoubleTag) String() string {
	return fmt.Sprintf("%s: %fd", t.Name, t.Data)
}

func (*DoubleTag) ID() byte {
	return 0x06
}

type ByteArrayTag struct {
	Name string
	Data []byte
}

func (t *ByteArrayTag) GetName() string {
	return t.Name
}

func (t *ByteArrayTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x07)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	buf.Write(util.Int32ToBytes(int32(len(t.Data))))
	buf.Write(t.Data)
	return buf.Bytes()
}

func (t *ByteArrayTag) String() string {
	var sb strings.Builder
	sb.WriteString(t.Name)
	sb.WriteString(": [B;")
	for i, b := range t.Data {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%db", b))
	}
	sb.WriteByte(']')
	return sb.String()
}

func (*ByteArrayTag) ID() byte {
	return 0x07
}

type StringTag struct {
	Name string
	Data string
}

func (t *StringTag) GetName() string {
	return t.Name
}

func (t *StringTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x08)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	writeString(&buf, t.Data)
	return buf.Bytes()
}

func (t *StringTag) String() string {
	return fmt.Sprintf(`%s: "%s"`, t.Name, t.Data)
}

func (*StringTag) ID() byte {
	return 0x08
}

type ListTag[T Tag] struct {
	Name string
	Data []T
}

func (t *ListTag[T]) GetName() string {
	return t.Name
}

func (t *ListTag[T]) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x09)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	if len(t.Data) > 0 {
		buf.WriteByte(t.Data[0].ID())
	} else {
		buf.WriteByte(0)
	}
	buf.Write(util.Int32ToBytes(int32(len(t.Data))))
	for _, e := range t.Data {
		buf.Write(e.Marshal()[1:])
	}
	return buf.Bytes()
}

func (t *ListTag[T]) String() string {
	var sb strings.Builder
	sb.WriteString(t.Name)
	sb.WriteString(": [")
	for i, e := range t.Data {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(e.String())
	}
	sb.WriteByte(']')
	return sb.String()
}

func (*ListTag[T]) ID() byte {
	return 0x09
}

type CompoundTag struct {
	Name string
	Data []Tag
}

func (t *CompoundTag) GetName() string {
	return t.Name
}

func (t *CompoundTag) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x0a)
	if t.Name != "" {
		writeString(&buf, t.Name)
	}
	for _, c := range t.Data {
		buf.Write(c.Marshal())
	}
	buf.WriteByte(0)
	return buf.Bytes()
}

func (t *CompoundTag) String() string {
	var sb strings.Builder
	if t.Name != "" {
		sb.WriteString(t.Name)
		sb.WriteString(": ")
	}
	sb.WriteByte('{')
	for i, c := range t.Data {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(c.String())
	}
	sb.WriteByte('}')
	return sb.String()
}

func (*CompoundTag) ID() byte {
	return 0x0a
}

func writeString(buf *bytes.Buffer, s string) {
	buf.Write(util.Uint16ToBytes(uint16(len(s))))
	buf.WriteString(s)
}
