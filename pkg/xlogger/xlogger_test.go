package xlogger

import (
	"go-multiple-query/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup_Development(t *testing.T) {
	cfg := config.Config{IsDevelopment: true}
	Setup(cfg)
	assert.NotNil(t, Logger)
}

func TestSetup_Production(t *testing.T) {
	cfg := config.Config{}
	Setup(cfg)
	assert.NotNil(t, Logger)
}
