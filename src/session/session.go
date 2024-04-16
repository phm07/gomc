package session

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gomc/src/connection"
	"gomc/src/encrypt"
	profile "gomc/src/profile"
	"io"
	"net/http"
	"net/url"
)

func ValidateSession(c *connection.Connection) (*profile.Profile, error) {

	h := sha1.New()
	h.Write([]byte(""))
	h.Write(c.Secret)
	h.Write(encrypt.PublicKeyBytes)
	serverId := AuthDigest(h.Sum(nil))

	values := url.Values{}
	values.Add("username", c.Username)
	values.Add("serverId", serverId)

	req, err := http.NewRequest("GET", "https://sessionserver.mojang.com/session/minecraft/hasJoined?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var schema struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Properties []struct {
			Name      string `json:"name"`
			Value     string `json:"value"`
			Signature string `json:"signature"`
		} `json:"properties"`
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, err
	}
	if c.Username != schema.Name {
		return nil, fmt.Errorf("mismatched username: %s != %s", c.Username, schema.Name)
	}

	uid, err := uuid.Parse(
		schema.Id[:8] + "-" +
			schema.Id[8:12] + "-" +
			schema.Id[12:16] + "-" +
			schema.Id[16:20] + "-" +
			schema.Id[20:])
	if err != nil {
		return nil, err
	}

	var properties []profile.Property
	for _, p := range schema.Properties {
		properties = append(properties, profile.Property{
			Name:      p.Name,
			Value:     p.Value,
			Signature: p.Signature,
		})
	}

	return &profile.Profile{
		Id:         uid,
		Name:       schema.Name,
		Properties: properties,
	}, nil
}
