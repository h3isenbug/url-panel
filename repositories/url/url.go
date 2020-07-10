package url

import "time"

type URL struct {
	ShortPath string    `json:"shortPath" db:"short_path"`
	LongURL   string    `json:"longURL" db:"long_url"`
	EMail     string    `json:"-" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type URLRepository interface {
	SaveShortPath(email, longURL, shortPath string) error
	DeleteURL(email, shortPath string) error
	UserOwnsURL(email, shortPath string) (bool, error)
	GetMyURLS(email string) ([]*URL, error)
}
