package models

import (
	"fmt"
)

// SocketHandler represents socket handler model.
type SocketHandler struct {
	tableName struct{} `pg:"?schema.sockets_sockethandler,discard_unknown_columns"` // nolint

	ID          int
	Metadata    JSON
	SocketID    int
	Socket      *Socket
	HandlerName string
	Handler     string
}

func (m *SocketHandler) String() string {
	return fmt.Sprintf("SocketHandler<ID=%d, HandlerName=%q>", m.ID, m.HandlerName)
}

// VerboseName returns verbose name for model.
func (m *SocketHandler) VerboseName() string {
	return "SocketHandler"
}
