package url

import (
	"github.com/go-redis/redis"
)

type RedisCacheWriteRepository struct {
	redisClient *redis.Client
	nextLayer   URLRepository
}

func NewRedisCacheWriteRepositoryV1(redisClient *redis.Client, nextLayer URLRepository) URLRepository {
	return &RedisCacheWriteRepository{redisClient: redisClient, nextLayer: nextLayer}
}

func (repo RedisCacheWriteRepository) SaveShortPath(email, longURL, shortPath string) error {
	if err := repo.nextLayer.SaveShortPath(email, longURL, shortPath); err != nil {
		return err
	}

	return repo.redisClient.Set(shortPath, longURL, 0).Err()
}

func (repo RedisCacheWriteRepository) DeleteURL(email, shortPath string) error {
	if err := repo.nextLayer.DeleteURL(email, shortPath); err != nil {
		return err
	}
	return repo.redisClient.Del(shortPath).Err()
}

func (repo RedisCacheWriteRepository) UserOwnsURL(email, shortPath string) (bool, error) {
	return repo.nextLayer.UserOwnsURL(email, shortPath)
}

func (repo RedisCacheWriteRepository) GetMyURLS(email string) ([]*URL, error) {
	return repo.nextLayer.GetMyURLS(email)
}
