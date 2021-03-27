// Package assets contains static assets for the website
package assets

import "embed"

// Assets embeds the assets into the go binary.
//go:embed index.html
//go:embed js/*
var Assets embed.FS

// Index is the embedded index page.
var Index = func() []byte {
	data, err := Assets.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	return data
}()
