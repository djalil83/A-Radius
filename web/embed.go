package web

import "embed"

// Assets berisi asset dashboard web yang ikut masuk ke binary API.
//
//go:embed dashboards/pelanggan/index.html dashboards/pelanggan/dashboard.js
var Assets embed.FS
