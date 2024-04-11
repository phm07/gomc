package nbt_test

import (
	"github.com/stretchr/testify/assert"
	"gomc/src/nbt"
	"os"
	"testing"
)

func TestNBT(t *testing.T) {
	bytes, err := os.ReadFile("testdata/test.nbt")
	if err != nil {
		t.Fatal(err)
	}

	tag := nbt.ShortTag{Name: "shortTest", Data: 32767}

	assert.Equal(t, bytes, tag.Marshal())
	assert.Equal(t, "shortTest: 32767", tag.String())
}

func TestBigNBT(t *testing.T) {

	bytes, err := os.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		t.Fatal(err)
	}

	byteArrayTest := make([]byte, 1000)
	for i := 0; i < len(byteArrayTest); i++ {
		byteArrayTest[i] = byte((i*i*255 + i*7) % 100)
	}

	tag := nbt.CompoundTag{
		Name: "Level",
		Data: []nbt.Tag{
			&nbt.ShortTag{Name: "shortTest", Data: 32767},
			&nbt.LongTag{Name: "longTest", Data: 9223372036854775807},
			&nbt.ByteTag{Name: "byteTest", Data: 127},
			&nbt.ByteArrayTag{
				Name: "byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))",
				Data: byteArrayTest,
			},
			&nbt.ListTag[*nbt.LongTag]{
				Name: "listTest (long)",
				Data: []*nbt.LongTag{{Data: 11}, {Data: 12}, {Data: 13}, {Data: 14}, {Data: 15}},
			},
			&nbt.FloatTag{Name: "floatTest", Data: 0.49823147058486938},
			&nbt.DoubleTag{Name: "doubleTest", Data: 0.49312871321823148},
			&nbt.IntTag{Name: "intTest", Data: 2147483647},
			&nbt.ListTag[*nbt.CompoundTag]{
				Name: "listTest (compound)",
				Data: []*nbt.CompoundTag{
					{Data: []nbt.Tag{
						&nbt.LongTag{Name: "created-on", Data: 1264099775885},
						&nbt.StringTag{Name: "name", Data: "Compound tag #0"},
					}},
					{Data: []nbt.Tag{
						&nbt.LongTag{Name: "created-on", Data: 1264099775885},
						&nbt.StringTag{Name: "name", Data: "Compound tag #1"},
					}},
				},
			},
			&nbt.CompoundTag{
				Name: "nested compound test",
				Data: []nbt.Tag{
					&nbt.CompoundTag{
						Name: "egg",
						Data: []nbt.Tag{
							&nbt.StringTag{Name: "name", Data: "Eggbert"},
							&nbt.FloatTag{Name: "value", Data: 0.5},
						},
					},
					&nbt.CompoundTag{
						Name: "ham",
						Data: []nbt.Tag{
							&nbt.StringTag{Name: "name", Data: "Hampus"},
							&nbt.FloatTag{Name: "value", Data: 0.75},
						},
					},
				},
			},
			&nbt.StringTag{Name: "stringTest", Data: "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!"},
		},
	}

	assert.Equal(t, bytes, tag.Marshal())
}
