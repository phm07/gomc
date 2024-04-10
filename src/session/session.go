package session

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gomc/src/connection"
	"gomc/src/encrypt"
	profile "gomc/src/profile"
	"net/http"
	"net/url"
)

func ValidateSession(c *connection.Connection) error {

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
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
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
	if err := json.NewDecoder(resp.Body).Decode(&schema); err != nil {
		return err
	}
	if c.Username != schema.Name {
		return fmt.Errorf("mismatched username: %s != %s", c.Username, schema.Name)
	}

	uid, err := uuid.Parse(
		schema.Id[:8] + "-" +
			schema.Id[8:12] + "-" +
			schema.Id[12:16] + "-" +
			schema.Id[16:20] + "-" +
			schema.Id[20:])
	if err != nil {
		return err
	}

	var properties []profile.Property
	for _, p := range schema.Properties {
		properties = append(properties, profile.Property{
			Name:      p.Name,
			Value:     p.Value,
			Signature: p.Signature,
		})
	}

	c.Profile = &profile.Profile{
		Id:         uid,
		Name:       schema.Name,
		Properties: properties,
	}
	return nil
}
