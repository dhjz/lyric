package main

import (
	"dlrc/service/atyy"
	"dlrc/service/qqmusic"
	"embed"
	"io/fs"
	"log"

	"flag"
	"fmt"
	"net/http"
)

//go:embed app
var f embed.FS

func main() {

	port := flag.Int("p", 40008, "server port")
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	st, _ := fs.Sub(f, "app")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(st))))

	http.HandleFunc("/songs1", atyy.HandleSongs)
	http.HandleFunc("/lyric1", atyy.HandleLyric)
	http.HandleFunc("/cookie1", atyy.HandleCookie)
	http.HandleFunc("/songs", qqmusic.HandleSongs)
	http.HandleFunc("/lyric", qqmusic.HandleLyric)

	fmt.Println("********app run on http://localhost" + addr + "/ 启动参数-p 指定端口 ********")
	log.Fatal(http.ListenAndServe(addr, nil))

}
