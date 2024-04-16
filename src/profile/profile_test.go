package profile

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOfflineProfile(t *testing.T) {
	assert.Equal(t, OfflineProfile("Notch"), &Profile{
		Name: "Notch",
		Id:   uuid.Must(uuid.Parse("b50ad385-829d-3141-a216-7e7d7539ba7f")),
	})
}
