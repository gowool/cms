package migrations

import (
	"embed"
	"io/fs"

	"github.com/gowool/cms/internal"
)

//go:embed *
var FS embed.FS

var PgFS = internal.Must(fs.Sub(FS, "pg"))
