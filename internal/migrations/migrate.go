package migrations

import "embed"

//go:embed scripts/*.sql
var Migrations embed.FS
