package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// get html content

	// req , err := http.NewRequest("GET", "https://ubuntu:optim123@kokoichi0206.mydns.jp/github/", nil)
	// fmt.Printf("req.Header: %v\n", req.Header)

	url, err := url.Parse("http://hoge%40example.com:password@example.jp/secret/")
	if err != nil {
		log.Fatal(err)
	}
	// url.Scheme: http
	fmt.Printf("url.Scheme: %v\n", url.Scheme)
	// url.Host: example.jp
	fmt.Printf("url.Host: %v\n", url.Host)
	// url.User: hoge%40example.com:password
	fmt.Printf("url.User: %v\n", url.User)
	// url.Path: /secret/
	fmt.Printf("url.Path: %v\n", url.Path)

	return

	// res, err := http.Get("http://hoge%40example.com:password@kokoichi0206.mydns.jp/github")
	res, err := http.Get("https://ubuntu:optim123@kokoichi0206.mydns.jp/github/")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// get header
	fmt.Printf("res.Status: %v\n", res.Status)
	io.Copy(os.Stdout, res.Body)
}
