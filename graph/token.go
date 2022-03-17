package graph

import (
	"auction-back/db"
	"fmt"
)

var validateTokenData = map[db.TokenAction]func(data map[string]interface{}) error{
	db.TokenActionApproveUserEmail: func(data map[string]interface{}) error {
		email, found := data["email"]

		if !found {
			return fmt.Errorf("no email in data")
		}

		_, ok := email.(string)
		if !ok {
			return fmt.Errorf("email in data is not string")
		}

		return nil
	},
	db.TokenActionApproveUserPhone: func(data map[string]interface{}) error {
		phone, found := data["phone"]

		if !found {
			return fmt.Errorf("no phone in data")
		}

		_, ok := phone.(string)
		if !ok {
			return fmt.Errorf("phone in data is not string")
		}

		return nil
	},
}
