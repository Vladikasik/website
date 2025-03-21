package handler

import (
	"log"

	"aynshteyn.dev/backend/internal/store"
)

// Application holds the application dependencies and configuration
type Application struct {
	Config string
	Logger *log.Logger
	Store  *store.SQLiteStore
} 