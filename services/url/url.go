package url

import (
	"encoding/base32"
	"errors"
	urlRepo "github.com/h3isenbug/url-panel/repositories/url"
	"math/rand"
)

const shortPathDefaultLenght = 8

type Service interface {
	CreateShortURL(email, recommended, longURL string) (string, error)
	DeleteShortURL(email, shortPath string) error
	GetMyURLS(email string) ([]*urlRepo.URL, error)
}

type ServiceV1 struct {
	urlRepository urlRepo.URLRepository
}

func NewURLServiceV1(urlRepository urlRepo.URLRepository) Service {
	return &ServiceV1{
		urlRepository: urlRepository,
	}
}

func (service ServiceV1) CreateShortURL(email, recommended, longURL string) (string, error) {
	var randomBytes = make([]byte, shortPathDefaultLenght*2)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	var randomString = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	var shortPath = recommended
	var postfix = randomString
	if recommended == "" {
		shortPath = randomString[:shortPathDefaultLenght]
		postfix = randomString[shortPathDefaultLenght:]
	}

	for tries := 0; tries < shortPathDefaultLenght; tries++ {
		if err := service.urlRepository.SaveShortPath(email, longURL, shortPath); err != nil {
			shortPath += postfix[tries : tries+1]
		} else {
			return shortPath, nil
		}
	}

	return "", errors.New("could not create unique random string")
}

func (service ServiceV1) DeleteShortURL(email, shortPath string) error {
	if err := service.urlRepository.DeleteURL(email, shortPath); err != nil {
		return err
	}

	return nil
}

func (service ServiceV1) GetMyURLS(email string) ([]*urlRepo.URL, error) {
	return service.urlRepository.GetMyURLS(email)
}
