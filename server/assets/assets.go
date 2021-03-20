// Package assets contains static assets for the website
package assets

import "embed"

//go:embed index.html
// Assets embeds the assets into the go binary.
var Assets embed.FS

// Index is the embedded index page.
var Index = func() []byte {
	data, err := Assets.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	return data
}()
