package graph

import (
	"github.com/brharrelldev/BlockListAPI/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	DB *database.BlockListDB
}