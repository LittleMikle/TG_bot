package server

import (
	"fmt"
	"github.com/LittleMikle/TG_bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/http"
	"strconv"
)

type AuthorizationServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

func NewAuthorizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectURL:     redirectURL}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":8081",
		Handler: s,
	}

	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIDParam := r.URL.Query().Get("chat_id")

	if chatIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatID, err := strconv.ParseInt(chatIDParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestToken, err := s.tokenRepository.Get(chatID, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println(requestToken)

	authRes, err := s.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		fmt.Println("FAILED WITH AUTHRES", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(authRes.AccessToken, authRes.Username)

	err = s.tokenRepository.Save(chatID, authRes.AccessToken, repository.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(chatID)
	fmt.Println(requestToken)
	fmt.Println(authRes.AccessToken)

	w.Header().Add("Location", s.redirectURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
