//+build wireinject

package main

import (
	"github.com/google/wire"
	"net/http"
)

func wireUp() (*http.Server, func(), error) {
	wire.Build(provideHTTPServer, provideAnalyticsHandler, provideLogService, provideMuxRouter, provideRedisClient, provideAuthHandler, providePanelHandler, provideURLService, provideAuthService, provideUserRepository, provideURLRepository)

	return &http.Server{}, func() {}, nil
}
