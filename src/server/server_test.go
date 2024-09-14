package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewServerSuccess(t *testing.T) {
	server, err := CreateNewServerUDP("1234", 1)
	assert.NoError(t, err)
	assert.Contains(t, server.Address(), ":1234")
}

func TestCreateNewServerNoPortProvided(t *testing.T) {
	server, err := CreateNewServerUDP("", 1)
	assert.Nil(t, server)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "No port specified")
}

func TestCreateNewServerNegativeLifespan(t *testing.T) {
	server, err := CreateNewServerUDP("1234", -1)
	assert.Nil(t, server)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "Lifetime cannot be negative")
}
