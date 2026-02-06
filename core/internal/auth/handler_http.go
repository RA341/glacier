package auth

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type HandlerHttp struct {
	srv *Service
}

func NewHandlerHttp(srv *Service) http.Handler {
	hand := &HandlerHttp{srv: srv}

	subMux := http.NewServeMux()
	subMux.HandleFunc("GET /oidc", hand.OIDCLogin)
	subMux.HandleFunc("GET /oidc/callback", hand.OIDCCallback)

	return subMux
}

const StateCookieName = "oidc_state"

// OIDCLogin GET /auth/login/google
func (h *HandlerHttp) OIDCLogin(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()

	url, err := h.srv.GetOIDCLoginURL(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     StateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// OIDCCallback GET /auth/callback
func (h *HandlerHttp) OIDCCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	oidcErr := query.Get("error")
	if oidcErr != "" {
		mess := fmt.Sprintf("error: %s, desc: %s", oidcErr, query.Get("error_description"))
		http.Error(w, mess, http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie(StateCookieName)
	if err != nil {
		http.Error(w, "State cookie missing", http.StatusBadRequest)
		return
	}

	returnedState := query.Get("state")
	if returnedState == "" || returnedState != cookie.Value {
		http.Error(w, "Invalid state parameter (CSRF detected)", http.StatusUnauthorized)
		return
	}

	code := query.Get("code")
	ctx := r.Context()

	normalLogin := func(sessionT SessionType) (session Session, sessionToken string, refreshToken string, err error) {
		return h.srv.LoginOIDC(ctx, code, sessionT)
	}

	err = loginSessionHandler(
		normalLogin,
		w,
		r.Header,
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("error Failed to verify OIDC token details %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	http.Redirect(w, r, "/library", http.StatusPermanentRedirect)
}
