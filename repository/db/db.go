package db

import (
	"context"
	"net/http"
	"time"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/jackc/pgx/v5"
)

type inDBRepository struct {
	conn        *pgx.Conn
	databaseDSN string
}

func New(cfg *config.Config) repository.Repository {

	return &inDBRepository{
		databaseDSN: cfg.DatabaseDSN,
	}
}

func (i *inDBRepository) SaveURL(shortID string, longURL string, cook string) error {
	return nil

}

func (i *inDBRepository) GetURL(shortID string) (string, bool) {
	return "", false

}

func (i *inDBRepository) GetAll(s string) []map[string]string {
	return nil
}

func (i *inDBRepository) Init() error {
	conn, err := pgx.Connect(context.Background(), i.databaseDSN)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}
	i.conn = conn
	return nil
}

func (i *inDBRepository) Ping() int {
	//
	//conn, err := pgx.Connect(context.Background(), i.databaseDSN)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := i.conn.Ping(ctx)

	status := http.StatusOK

	if err != nil {
		status = http.StatusInternalServerError
	}
	return status
}