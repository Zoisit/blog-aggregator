package cli

import (
	"github.com/Zoisit/blog-aggregator/internal/config"
	"github.com/Zoisit/blog-aggregator/internal/database"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}
