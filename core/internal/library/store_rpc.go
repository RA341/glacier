package library

import (
	"time"

	v1 "github.com/ra341/glacier/generated/library/v1"
	types2 "github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/internal/metadata/types"
	"github.com/rs/zerolog/log"
)

func (g *Game) ToProto() *v1.Game {
	return &v1.Game{
		ID:            uint64(g.ID),
		CreatedAt:     g.CreatedAt.Format(time.RFC3339),
		EditedAt:      g.UpdatedAt.Format(time.RFC3339),
		Meta:          g.Meta.ToProto(),
		GameType:      g.GameType.String(),
		DownloadState: g.Download.ToProto(),
	}
}

func (g *Game) FromProto(rpcGame *v1.Game) {
	// do not update this it should be handle by DB
	//g.UpdatedAt
	//g.CreatedAt

	meta := &types.Meta{}
	meta.FromProto(rpcGame.Meta)

	gameType, err := GameTypeString(rpcGame.GameType)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse game type")
		gameType = GameTypeUnknown
	}

	down := &types2.Download{}
	down.FromProto(rpcGame.DownloadState)

	g.ID = uint(rpcGame.ID)
	g.GameType = gameType
	g.Meta = *meta
	g.Download = *down
}
