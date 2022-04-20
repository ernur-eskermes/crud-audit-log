package domain

import (
	"errors"
	"time"
)

const (
	EntityUser = "USER"
	EntityBook = "BOOK"

	ActionCreate   = "CREATE"
	ActionUpdate   = "UPDATE"
	ActionGet      = "GET"
	ActionDelete   = "DELETE"
	ActionRegister = "REGISTER"
	ActionLogin    = "LOGIN"
)

var (
	entities = map[string]LogRequest_Entities{
		EntityUser: LogRequest_USER,
		EntityBook: LogRequest_BOOK,
	}

	actions = map[string]LogRequest_Actions{
		ActionCreate:   LogRequest_CREATE,
		ActionUpdate:   LogRequest_UPDATE,
		ActionGet:      LogRequest_GET,
		ActionDelete:   LogRequest_DELETE,
		ActionRegister: LogRequest_REGISTER,
		ActionLogin:    LogRequest_LOGIN,
	}
)

type LogItem struct {
	Entity    string    `json:"entity" bson:"entity"`
	Action    string    `json:"action" bson:"action"`
	EntityID  string    `json:"entity_id" bson:"entity_id"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

func ToPbEntity(entity string) (LogRequest_Entities, error) {
	val, ex := entities[entity]
	if !ex {
		return 0, errors.New("invalid entity")
	}

	return val, nil
}

func ToPbAction(action string) (LogRequest_Actions, error) {
	val, ex := actions[action]
	if !ex {
		return 0, errors.New("invalid action")
	}

	return val, nil
}
