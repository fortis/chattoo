package server

import (
	"errors"
	"gopkg.in/validator.v2"
)

type person struct {
	Id       int64  `json:"id" validate:"nonzero"`
	Username string `json:"username" validate:"nonzero"`
}

const (
	publicMessageType  = "public"
	privateMessageType = "private"
)

type message struct {
	Type string  `json:"type" validate:"nonzero"`
	From *person `json:"from" validate:"nonzero"`
	To   *person `json:"to"`
	Body string  `json:"body" validate:"nonzero"`
}

func (m message) validate() error {
	if err := validator.Validate(m); err != nil {
		return err
	}

	if m.Type == privateMessageType && m.To == nil {
		return errors.New("recipient required")
	}

	switch m.Type {
	case publicMessageType, privateMessageType:
		return nil
	default:
		return errors.New("unknown message type")
	}
}
