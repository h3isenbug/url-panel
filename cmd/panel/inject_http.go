package main

import (
	"context"
	"github.com/gorilla/mux"
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"github.com/h3isenbug/url-panel/config"
	"github.com/h3isenbug/url-panel/handlers"
	"github.com/h3isenbug/url-panel/handlers/analytics"
	"github.com/h3isenbug/url-panel/handlers/auth"
	"github.com/h3isenbug/url-panel/handlers/panel"
	auth2 "github.com/h3isenbug/url-panel/services/auth"
	urlService "github.com/h3isenbug/url-panel/services/url"
	"net/http"
)

func provideHTTPServer(router *mux.Router) (*http.Server, func()) {
	server := &http.Server{
		Addr:    ":" + config.Config.Port,
		Handler: router,
	}
	return server, func() { server.Shutdown(context.Background()) }
}

func provideAnalyticsHandler() (analytics.AnalyticsHandler, error) {
	return analytics.NewAnalyticsReverseProxyHandler(config.Config.AnalyticsServer)
}

func provideAuthHandler(authService auth2.AuthService, logService log2.LogService) auth.AuthHandler {
	return auth.NewAuthHandlerV1(authService, logService)
}

func providePanelHandler(urlService urlService.Service, logService log2.LogService) panel.PanelHandler {
	return panel.NewPanelHandlerV1(urlService, logService)
}

func provideMuxRouter(analyticsHandler analytics.AnalyticsHandler, authHandler auth.AuthHandler, panelHandler panel.PanelHandler, authService auth2.AuthService) *mux.Router {
	router := mux.NewRouter()
	router.Use(handlers.GorillaMuxURLParamMiddleware)
	router.Methods("POST").Path("/login").HandlerFunc(authHandler.Login)
	router.Methods("POST").Path("/register").HandlerFunc(authHandler.Register)

	panelRouter := router.PathPrefix("/panel").Subrouter()
	panelRouter.Use(handlers.AuthTokenMiddleware(authService))
	panelRouter.Methods("GET").Path("/analytics/{short_path}").HandlerFunc(analyticsHandler.GetURLAnalytics)

	panelRouter.Methods("GET").Path("/urls").HandlerFunc(panelHandler.GetMyURLS)
	panelRouter.Methods("POST").Path("/urls").HandlerFunc(panelHandler.CreateShortURL)
	panelRouter.Methods("DELETE").Path("/urls/{short_path}").HandlerFunc(panelHandler.DeleteShortURL)
	return router
}
