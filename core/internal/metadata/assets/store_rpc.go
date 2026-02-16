package assets

import (
	v1 "github.com/ra341/glacier/generated/assets/v1"
)

func (m *Asset) ToProto() *v1.Asset {
	return &v1.Asset{
		ID:        uint64(m.ID),
		GameId:    uint64(m.GameID),
		Type:      m.Type.String(),
		RemoteUrl: m.RemoteURL,
		LocalPath: m.LocalPath,
	}
}

func (m *Asset) FromProto(rpcMeta *v1.Asset) {
	if rpcMeta == nil {
		return
	}

	m.ID = uint(rpcMeta.ID)
	m.GameID = uint(rpcMeta.GameId)
	m.LocalPath = rpcMeta.LocalPath
	m.RemoteURL = rpcMeta.RemoteUrl

	typeString, err := AssetTypeString(rpcMeta.Type)
	if err != nil {
		typeString = AssetUnknown
	}
	m.Type = typeString
}
