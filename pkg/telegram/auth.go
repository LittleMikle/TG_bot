package telegram

import (
	"context"
	"fmt"
	"github.com/LittleMikle/TG_bot/pkg/repository"
)

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

func (b *Bot) generateAuthorizationLink(chatId int64) (string, error) {
	redirectURL := b.generateRedirectLink(chatId)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	err = b.tokenRepository.Save(chatId, requestToken, repository.RequestTokens)
	if err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

func (b *Bot) generateRedirectLink(chatId int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatId)
}
