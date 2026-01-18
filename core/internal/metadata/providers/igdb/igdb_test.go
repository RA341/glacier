package igdb

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	err := godotenv.Load()
	require.NoError(t, err)
	cliId := os.Getenv("CLIID")
	cliSecret := os.Getenv("CLISECRET")

	config := map[string]any{
		"clientId":     cliId,
		"clientSecret": cliSecret,
	}

	srv, err := New(config)
	require.NoError(t, err)

	games, err := srv.GetMatches("warhammer")
	require.NoError(t, err)

	//time.Sleep(1 * time.Second)
	//games, err = srv.searchGames("forza")
	//require.NoError(t, err)

	t.Log(games)
}
