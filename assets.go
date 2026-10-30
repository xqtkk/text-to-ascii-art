package assets

import "embed"

// Files contains the templates and static files required by the web app.
//
//go:embed templates/*.html ui/static
var Files embed.FS
