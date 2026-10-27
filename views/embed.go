// Package views はフロントエンド (template/) のビルド成果物をバイナリに埋め込む。
// vite の outDir がこのディレクトリなので、vite.config.ts で emptyOutDir を有効にしないこと (このファイルが消える)。
package views

import "embed"

//go:embed index.html favicon.ico vite.svg assets
var FS embed.FS
