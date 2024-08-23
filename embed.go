package accountant

import "embed"

//go:embed .env migrations
var EmbedFs embed.FS
