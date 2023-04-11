package main

import "github.com/qsoulior/auth-server/internal/app"

func main() {
	s := app.New()
	s.ListenAndServe()
}
