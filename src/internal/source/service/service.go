package service

import (
	"context"
	"xyz_multifinance/src/internal/shared/mysql"
	"xyz_multifinance/src/internal/source/repository"
	"xyz_multifinance/src/internal/source/usecase"
)

func NewApplication(ctx context.Context) usecase.Usecase {
	db, err := mysql.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepository(db)

	return usecase.NewUsecase(repo)
}
