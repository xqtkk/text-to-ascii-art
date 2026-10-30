package main

import (
	"flag"
	"text-to-ascii-art/web"
)

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	flag.Parse()

	web.Web(addr)
}
