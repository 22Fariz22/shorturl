package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/22Fariz22/shorturl/model"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/jackc/pgx/v5"
)

type inDBRepository struct {
	conn        *pgx.Conn
	databaseDSN string
	buffer      []model.PackResponse
	ctx         context.Context
}

func (i *inDBRepository) RepoBatch(ctx context.Context, cook string, batchList []model.PackReq) error {
	for b := range batchList {
		_, err := i.conn.Exec(ctx, "insert into urls (cookies,correlation_id, short_url, long_url) values($1,$2,$3,$4);",
			cook, batchList[b].CorrelationID, batchList[b].ShortURL, batchList[b].OriginalURL)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func New(cfg *config.Config) repository.Repository {
	return &inDBRepository{
		databaseDSN: cfg.DatabaseDSN,
		buffer:      make([]model.PackResponse, 0, 1000),
	}
}

func (i *inDBRepository) Init() error {
	conn, err := pgx.Connect(context.Background(), i.databaseDSN)
	if err != nil {
		return err
	}
	i.conn = conn

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = conn.Exec(ctx,
		"CREATE TABLE if not exists urls(id_pk SERIAL PRIMARY KEY, cookies TEXT, correlation_id TEXT,"+
			" short_url TEXT, long_url TEXT UNIQUE, deleted boolean default false);")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (i *inDBRepository) Delete(list []string, cookie string) error {
	log.Println("begin Delete")
	tx, err := i.conn.Begin(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Prepare(context.Background(), "UPDATE", "UPDATE urls SET deleted = true WHERE short_url = $1;")
	log.Println("after prepare")
	if err != nil {
		log.Println(err)
		return err
	}

	for i := range list {
		log.Println("before Exec")
		_, err = tx.Exec(context.Background(), "UPDATE urls SET deleted = true WHERE short_url = $1 ;", list[i])
		log.Println("after exec")
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = tx.Commit(context.Background())
	log.Println("after commit")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *inDBRepository) SaveURL(ctx context.Context, shortURL string, longURL string, cook string) (string, error) {
	var s string
	var ErrAlreadyExists = errors.New("this URL already exists")
	err := i.conn.QueryRow(ctx, `
		   				WITH e AS(
		   				INSERT INTO urls (cookies, long_url, short_url)
		   					   VALUES ($1, $2, $3)
		   				ON CONFLICT("long_url") DO NOTHING
		   				RETURNING long_url
						)
						SELECT long_url FROM e
						Union
						SELECT short_url FROM urls where long_url=$3
		   ;`, cook, longURL, shortURL).Scan(&s)
	if err != nil {
		fmt.Println("err in queryrow", err)
	}
	if s != longURL {
		return s, ErrAlreadyExists
	}

	_, err = i.conn.Exec(ctx, "insert into urls (cookies, short_url, long_url) values($1,$2,$3);", cook, shortURL, longURL)
	if err != nil {
		log.Println(err)
	}
	return "", nil
}

// GetURL return long url, deleted status, exist row
func (i *inDBRepository) GetURL(ctx context.Context, shortID string) (model.URL, bool) {
	row, err := i.conn.Query(ctx, "select cookies,long_url,deleted from urls where short_url = $1;", shortID)
	if err != nil {
		log.Println(err)
		//TODO сделать возврат ошибки
		//	http.NotFound(w, r) ?
		return model.URL{}, false
	}
	defer row.Close()

	//type longAndDeleted struct {
	//	long    string
	//	deleted bool
	//}
	rows := make([]model.URL, 0)

	for row.Next() {
		var s model.URL
		err := row.Scan(&s.Cookies, &s.LongURL, &s.Deleted)
		if err != nil {
			return s, false
		}

		rows = append(rows, s)
	}

	if len(rows) == 0 {
		return model.URL{}, false
	}

	return rows[0], true
}

//example [map[7PJPPAZ:http://ya.ru] map[JRK5X81:http://ya.ru]]
func (i *inDBRepository) GetAll(ctx context.Context, cook string) ([]map[string]string, error) {
	rows, err := i.conn.Query(ctx, "select short_url, long_url from urls where cookies = $1;", cook)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	list := make([]map[string]string, 0)

	for rows.Next() {
		var id string
		var longurl string
		err = rows.Scan(&id, &longurl)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		aMap := map[string]string{id: longurl}

		list = append(list, aMap)
	}
	return list, nil
}

func (i *inDBRepository) Ping(ctx context.Context) error {
	return i.conn.Ping(ctx)
}
