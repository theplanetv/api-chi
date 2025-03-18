package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DatabaseService(t *testing.T) {
	service := DatabaseService{}

	t.Run("Connection success", func(t *testing.T) {
		err := service.Open()
		assert.NoError(t, err)
	})
}
