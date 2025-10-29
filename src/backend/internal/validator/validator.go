package validator

import (
	"backend/internal/model"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateOrder(raw []byte) (*model.Order, error) {
	var o model.Order
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, err
	}
	return &o, v.Struct(&o)
}
