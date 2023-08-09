// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-realworld/internal/biz"
	"kratos-realworld/internal/conf"
	"kratos-realworld/internal/data"
	"kratos-realworld/internal/server"
	"kratos-realworld/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	database, err := data.NewDB(confData)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(confData, logger, database)
	if err != nil {
		return nil, nil, err
	}
	backendRepo := data.NewBackendRepo(dataData, logger)
	backendUsecase := biz.NewBackendUsecase(backendRepo, logger)
	userRepo := data.NewUserRepo(dataData, logger)
	transaction := data.NewTransaction(dataData)
	userUsecase := biz.NewUserUsecase(userRepo, logger, jwt, transaction)
	fileLocalRepo := data.NewFileLocalRepo(dataData, confData, logger)
	fileRepo := data.NewFileRepo(dataData, logger)
	fileUsecase := biz.NewFileUsecase(fileLocalRepo, fileRepo, confData, transaction, logger)
	backendService := service.NewBackendService(backendUsecase, userUsecase, fileUsecase)
	grpcServer := server.NewGRPCServer(confServer, backendService, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, backendService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
