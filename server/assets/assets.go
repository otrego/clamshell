// Package assets contains
package assets

import "embed"

//go:embed index.html
var Assets embed.FS

var Index = func() []byte {
	data, err := Assets.ReadFile("index.html")
	if err != nil {
		panic(err)
	}
	return data
}()
