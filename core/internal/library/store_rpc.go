package library

import (
	"time"

	v1 "github.com/ra341/glacier/generated/library/v1"
	downloaderTypes "github.com/ra341/glacier/internal/downloader/types"
	indexTypes "github.com/ra341/glacier/internal/indexer/types"
	metaTypes "github.com/ra341/glacier/internal/metadata/types"
)

func (g *GameFile) ToProto() *v1.GameFile {
	return &v1.GameFile{
		Exe: g.Exe,
	}
}

func (g *Game) ToProto() *v1.Game {
	return &v1.Game{
		ID:            uint64(g.ID),
		CreatedAt:     g.CreatedAt.Format(time.RFC3339),
		EditedAt:      g.UpdatedAt.Format(time.RFC3339),
		Meta:          g.Meta.ToProto(),
		DownloadState: g.Download.ToProto(),
		Source:        g.Source.ToProto(),
		File:          g.File.ToProto(),
	}
}

func (g *GameFile) FromProto(file *v1.GameFile) {
	if file == nil {
		return
	}

	g.Exe = file.Exe
}

func (g *Game) FromProto(rpcGame *v1.Game) {
	if rpcGame == nil {
		return
	}

	// do not update this it should be handled by DB
	//g.UpdatedAt
	//g.CreatedAt

	meta := metaTypes.Meta{}
	meta.FromProto(rpcGame.Meta)

	down := downloaderTypes.Download{}
	down.FromProto(rpcGame.DownloadState)

	src := indexTypes.Source{}
	src.FromProto(rpcGame.Source)

	file := GameFile{}
	file.FromProto(rpcGame.File)

	g.ID = uint(rpcGame.ID)
	g.File = file
	g.Meta = meta
	g.Download = down
	g.Source = src
}
