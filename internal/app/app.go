package app

import (
	"fmt"

	"github.com/qsoulior/auth-server/internal/repo"
	"github.com/qsoulior/auth-server/internal/usecase"
	"github.com/qsoulior/auth-server/pkg/db"
)

func Run() {
	pg, err := db.NewPostgres("postgres://postgres:test123@localhost:5432/postgres?search_path=app")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pg.Close()

	tu := usecase.NewToken(repo.NewTokenPostgres(pg))
	uu := usecase.NewUser(tu, repo.NewUserPostgres(pg))

	server := NewServer(uu, tu)
	server.ListenAndServe()
}
