package api

import (
	"os"
	"testing"
	"time"

	db "github.com/MathHRM/simple_bank/db/sqlc"
	"github.com/MathHRM/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config {
		TokenSimetricKey: util.RandomString(32),
		AcessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}