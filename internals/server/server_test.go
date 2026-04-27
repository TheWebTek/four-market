package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	srv := New()
	assert.NotNil(t, srv)
}

func TestNew_WithPort(t *testing.T) {
	srv := New(WithPort("3000"))
	assert.NotNil(t, srv)
}

func TestNew_WithEnvPort(t *testing.T) {
	os.Setenv("PORT", "9999")
	defer os.Unsetenv("PORT")
	srv := New()
	es := srv.(*echoServer)
	assert.Equal(t, "9999", es.config.Port)
}

func TestServerInterface(t *testing.T) {
	var _ Server = (*echoServer)(nil)
}
