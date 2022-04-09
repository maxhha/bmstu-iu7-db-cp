package token

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
)

type Interface interface {
	Create(action models.TokenAction, viewer models.User, data map[string]interface{}) error
	Activate(action models.TokenAction, token_code string, viewer models.User) (models.Token, error)
}

type SenderInterface interface {
	Name() string
	Send(models.Token) (bool, error)
}

type TokenPort struct {
	db      ports.DB
	senders []SenderInterface
}

func New(db ports.DB, senders []SenderInterface) TokenPort {
	return TokenPort{db, senders}
}

func (t *TokenPort) send(token models.Token) error {
	var errors error
	hasNotifier := false

	for _, n := range t.senders {
		sent, err := n.Send(token)

		if err != nil {
			hasNotifier = true
			errors = multierror.Append(
				errors,
				fmt.Errorf("%s send: %w", n.Name(), err),
			)
		} else if sent {
			hasNotifier = true
		}
	}

	if !hasNotifier {
		errors = multierror.Append(
			errors,
			fmt.Errorf("acton %s does not have any notifier", token.Action),
		)
	}

	return errors
}

func (t *TokenPort) Create(action models.TokenAction, viewer models.User, data map[string]interface{}) error {
	token := models.Token{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)),
		Action:    action,
		Data:      data,
		UserID:    viewer.ID,
	}

	if err := t.db.Token().Create(&token); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if err := t.send(token); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (t *TokenPort) Activate(action models.TokenAction, token_code string, viewer models.User) (models.Token, error) {
	token, err := t.db.Token().Take(ports.TokenTakeConfig{
		IDs:     []string{token_code},
		UserIDs: []string{viewer.ID},
	})
	if err != nil {
		return token, fmt.Errorf("db take: %w", err)
	}

	if token.Action != action {
		return token, fmt.Errorf("action not match")
	}

	if token.ActivatedAt.Valid {
		return token, fmt.Errorf("already activated")
	}

	token.ActivatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err := t.db.Token().Update(&token); err != nil {
		return token, fmt.Errorf("save: %w", err)
	}

	return token, nil
}
