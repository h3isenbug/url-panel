package main

import (
	log2 "github.com/h3isenbug/url-common/pkg/services/log"
	"github.com/h3isenbug/url-panel/config"
	url2 "github.com/h3isenbug/url-panel/repositories/url"
	"github.com/h3isenbug/url-panel/repositories/user"
	"github.com/h3isenbug/url-panel/services/auth"
	"github.com/h3isenbug/url-panel/services/url"
	"os"
)

func provideURLService(urlRepo url2.URLRepository) url.Service {
	return url.NewURLServiceV1(urlRepo)
}

func provideAuthService(userRepo user.UserRepository) auth.AuthService {
	return auth.NewAuthServiceV1(userRepo, config.Config.JWTKey)
}

func provideLogService() log2.LogService {
	return log2.NewLogServiceV1(os.Stdout)
}
