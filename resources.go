package main

import "embed"

//go:embed resources/fonts/SourceHanSansCN-Regular.ttf
var chineseFontFS embed.FS

//go:embed resources/images/enemy.svg resources/images/player.svg resources/images/enemy.png resources/images/player.png
var imagesFS embed.FS
