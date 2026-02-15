package assets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type HandlerHttp struct {
	srv *Service
}

func NewHandler(srv *Service) http.Handler {
	h := &HandlerHttp{
		srv: srv,
	}

	admin := user.RequireRole(user.TechPriest)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{gameId}/t/{assetType}", h.getAssetFromQuery)
	mux.HandleFunc("GET /{gameId}/{assetPath}", h.getAsset)

	// only admins may upload/modify
	mux.Handle("POST /upload", admin(http.HandlerFunc(h.uploadAsset)))
	mux.Handle("GET /types", admin(http.HandlerFunc(h.getAssetTypes)))
	mux.Handle("GET /delete/{id}", admin(http.HandlerFunc(h.deleteAsset)))

	return mux
}

func (h *HandlerHttp) getAssetTypes(w http.ResponseWriter, r *http.Request) {
	res := AssetTypeStrings()

	marshal, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Could not marshal JASON"+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.Warn().Err(err).Msg("Could not write response")
		return
	}
}

func (h *HandlerHttp) uploadAsset(w http.ResponseWriter, r *http.Request) {
	var TenMB int64 = 10 << 20
	if err := r.ParseMultipartForm(TenMB); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer fileutil.Close(file)

	assetTypeStr := r.FormValue("type")
	assetType, err := AssetTypeString(assetTypeStr)
	if err != nil {
		http.Error(w, "Invalid asset type: "+err.Error(), http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	gameIDStr := r.FormValue("gameId")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid GameID format", http.StatusBadRequest)
		return
	}

	asset := &Asset{
		Model: gorm.Model{
			ID: uint(id),
		},
		GameID: uint(gameID),
		Type:   assetType,
	}

	err = h.srv.SaveLocalAsset(r.Context(), asset, file)
	if err != nil {
		http.Error(w, "Failed to upload: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerHttp) getAssetFromQuery(w http.ResponseWriter, r *http.Request) {
	gameID, err := getStrId(r.PathValue("gameId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assetTypeStr := r.PathValue("assetType")
	if assetTypeStr == "" {
		http.Error(w, "Missing type", http.StatusBadRequest)
		return
	}

	assetType, err := AssetTypeString(assetTypeStr)
	if err != nil {
		http.Error(w, "Invalid asset type: "+err.Error(), http.StatusBadRequest)
		return
	}

	assetPath, err := h.srv.GetAssetByType(r.Context(), uint(gameID), assetType)
	if err != nil {
		http.Error(w, "Failed to get asset: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, assetPath)
}

func (h *HandlerHttp) getAsset(w http.ResponseWriter, r *http.Request) {
	gameID, err := getStrId(r.PathValue("gameId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assetPath := r.PathValue("assetPath")
	if assetPath == "" {
		http.Error(w, fmt.Sprintf("invalid asset path provided"), http.StatusBadRequest)
		return
	}

	asset, err := h.srv.GetAsset(
		r.Context(),
		gameID,
		assetPath,
	)
	if err != nil {
		return
	}

	http.ServeFile(w, r, asset)
}

func (h *HandlerHttp) deleteAsset(w http.ResponseWriter, r *http.Request) {
	assetId, err := getStrId(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.delete(r.Context(), uint(assetId))
	if err != nil {
		http.Error(w, "Failed to delete asset: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getStrId(strId string) (int, error) {
	if strId == "" {
		return 0, fmt.Errorf("invalid id")
	}

	return strconv.Atoi(strId)
}
