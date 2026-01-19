package types

import (
	v1 "github.com/ra341/glacier/generated/search/v1"
	"github.com/rs/zerolog/log"
)

func (ig *Source) ToProto() *v1.GameSource {
	return &v1.GameSource{
		IndexerType: ig.IndexerType.String(),
		GameType:    ig.GameType.String(),
		Title:       ig.Title,
		DownloadUrl: ig.DownloadUrl,
		ImageURL:    ig.ImageURL,
		FileSize:    ig.FileSize,
		CreatedISO:  ig.CreatedISO,
	}
}

func (ig *Source) FromProto(rpcGame *v1.GameSource) {
	indexerType, err := IndexerTypeString(rpcGame.IndexerType)
	if err != nil {
		log.Warn().Err(err).Str("name", rpcGame.IndexerType).Msg("invalid indexer type")
		indexerType = IndexerUnknown
	}

	gameType, err := GameTypeString(rpcGame.GameType)
	if err != nil {
		log.Warn().Err(err).Str("name", rpcGame.GameType).Msg("invalid gametype type")
		gameType = GameTypeUnknown
	}

	ig.IndexerType = indexerType
	ig.GameType = gameType
	ig.Title = rpcGame.Title
	ig.DownloadUrl = rpcGame.DownloadUrl
	ig.ImageURL = rpcGame.ImageURL
	ig.FileSize = rpcGame.FileSize
	ig.CreatedISO = rpcGame.CreatedISO
}
