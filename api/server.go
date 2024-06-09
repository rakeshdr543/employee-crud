package api

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Server struct {
	store *sql.DB
}

func NewServer(db *sql.DB) *Server {
	server := &Server{
		store: db,
	}

	return server
}
