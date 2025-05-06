package web

import "embed"

//go:embed webui/*
var EmbeddedWebFiles embed.FS
