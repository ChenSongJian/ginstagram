package main

import (
	"github.com/ChenSongJian/ginstagram/utils/db"
	"github.com/ChenSongJian/ginstagram/web"
)

func main() {
	db.InitDB()

	r := web.NewRouter()
	r.Run()
}
