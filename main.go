package main

import (
	"github.com/zitryss/perfmon/delivery/http"
	"github.com/zitryss/perfmon/domain"
	"github.com/zitryss/perfmon/infrastructure/database"
	"github.com/zitryss/perfmon/internal/log"
)

func main() {
	log.Setup()
	db := database.NewPsql("prod")
	use := domain.NewUsecase(db)
	svr := http.NewServer(use)
	svr.Start()
}
