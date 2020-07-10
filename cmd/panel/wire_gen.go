// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"net/http"
)

// Injectors from inject_app.go:

func wireUp() (*http.Server, func(), error) {
	analyticsHandler, err := provideAnalyticsHandler()
	if err != nil {
		return nil, nil, err
	}
	userRepository, cleanup, err := provideUserRepository()
	if err != nil {
		return nil, nil, err
	}
	authService := provideAuthService(userRepository)
	logService := provideLogService()
	authHandler := provideAuthHandler(authService, logService)
	client, cleanup2 := provideRedisClient()
	urlRepository, cleanup3, err := provideURLRepository(client)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	service := provideURLService(urlRepository)
	panelHandler := providePanelHandler(service, logService)
	router := provideMuxRouter(analyticsHandler, authHandler, panelHandler, authService)
	server, cleanup4 := provideHTTPServer(router)
	return server, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
