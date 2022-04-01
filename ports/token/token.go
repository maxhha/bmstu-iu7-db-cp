package token

import (
	"auction-back/db"
	"auction-back/grpc/notifier"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

type Interface interface {
	Create(action db.TokenAction, viewer *db.User, data map[string]interface{}) error
	Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error)
}

type TokenPort struct {
	db            *gorm.DB
	emailNotifier notifier.NotifierClient
}

// TODO: make array of senders with interface{ Send(db.Token) error }
// they would decide to send or not to send token
func New(db *gorm.DB, email notifier.NotifierClient) TokenPort {
	return TokenPort{db, email}
}

var actionsForEmailNotification = map[db.TokenAction]struct{}{
	db.TokenActionSetUserEmail: {},
}

var actionToDataGetter = map[db.TokenAction]func(db *gorm.DB, token db.Token) (map[string]string, error){
	db.TokenActionSetUserEmail: func(_ *gorm.DB, token db.Token) (map[string]string, error) {
		return map[string]string{
			"token": fmt.Sprintf("%06d", token.ID),
		}, nil
	},
}

var actionToEmailReceiverGetter = map[db.TokenAction]func(db *gorm.DB, token db.Token) (string, error){
	db.TokenActionSetUserEmail: func(_ *gorm.DB, token db.Token) (string, error) {
		email, ok := token.Data["email"]
		if !ok {
			return "", fmt.Errorf("no email in token data")
		}

		str, ok := email.(string)
		if !ok {
			return "", fmt.Errorf("fail convert to string of %v", email)
		}

		return str, nil
	},
}

func (t *TokenPort) sendEmail(token db.Token) error {
	var receiver string
	var data map[string]string

	if getData, has := actionToDataGetter[token.Action]; has {
		var err error
		if data, err = getData(t.db, token); err != nil {
			return fmt.Errorf("get data for %s: %w", token.Action, err)
		}
	}

	if getReceiver, has := actionToEmailReceiverGetter[token.Action]; has {
		var err error
		if receiver, err = getReceiver(t.db, token); err != nil {
			return fmt.Errorf("get receiver for %s: %w", token.Action, err)
		}
	} else {
		return fmt.Errorf("no receiver getter for %s", token.Action)
	}

	input := notifier.SendInput{
		Receivers: []string{receiver},
		Action:    string(token.Action),
		Data:      data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.emailNotifier.Send(ctx, &input)
	if err != nil {
		return fmt.Errorf("email notifier: %w", err)
	}

	if result.Status != "OK" {
		return fmt.Errorf("email notifier status: %v", result.Status)
	}

	return nil
}

func (t *TokenPort) send(token db.Token) error {
	action := token.Action
	var errors error
	hasNotifier := false

	if _, shouldSendEmail := actionsForEmailNotification[action]; shouldSendEmail {
		hasNotifier = true
		if err := t.sendEmail(token); err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	if !hasNotifier {
		errors = multierror.Append(
			errors,
			fmt.Errorf("acton %s does not have any notifier", action),
		)
	}

	return errors
}

func (t *TokenPort) Create(action db.TokenAction, viewer *db.User, data map[string]interface{}) error {
	token := db.Token{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)),
		Action:    action,
		Data:      data,
		UserID:    viewer.ID,
	}

	if err := t.db.Create(&token).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if err := t.send(token); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (t *TokenPort) Activate(action db.TokenAction, token_code string, viewer *db.User) (db.Token, error) {
	if viewer == nil {
		return db.Token{}, fmt.Errorf("unauthorized")
	}

	token := db.Token{}

	if err := t.db.Take(&token, "id = ? and user_id = ?", token_code, viewer.ID).Error; err != nil {
		return db.Token{}, fmt.Errorf("take: %w", err)
	}

	if token.Action != action {
		return db.Token{}, fmt.Errorf("action not match")
	}

	if token.ActivatedAt.Valid {
		return token, fmt.Errorf("already activated")
	}

	token.ActivatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err := t.db.Save(&token).Error; err != nil {
		return db.Token{}, fmt.Errorf("save: %w", err)
	}

	return token, nil
}
