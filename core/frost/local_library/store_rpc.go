package download

import (
	"github.com/ra341/glacier/frost/local_library/download"
	v1 "github.com/ra341/glacier/generated/frost_library/v1"
	"github.com/ra341/glacier/internal/library"
)

func (g *LocalGame) ToProto() *v1.LocalGame {
	return &v1.LocalGame{
		ID:       uint64(g.ID),
		GameID:   uint64(g.GameId),
		Game:     g.Game.ToProto(),
		Download: g.Download.ToProto(),
		Play:     g.Play.ToProto(),
	}
}

func (g *LocalGame) FromProto(gameRpc *v1.LocalGame) {
	var glacierGame library.Game
	glacierGame.FromProto(gameRpc.Game)

	var lldown download.LocalDownload
	lldown.FromProto(gameRpc.Download)

	var llPlay GamePlay
	llPlay.FromProto(gameRpc.Play)

	g.ID = uint(gameRpc.ID)
	g.GameId = int(gameRpc.GameID)
	g.Game = glacierGame
	g.Download = lldown
	g.Play = llPlay
}

func (p *GamePlay) ToProto() *v1.GamePlay {
	return &v1.GamePlay{
		LaunchExe: p.LaunchExe,
	}
}

func (p *GamePlay) FromProto(play *v1.GamePlay) {
	if play == nil {
		return
	}

	p.LaunchExe = play.LaunchExe
}
