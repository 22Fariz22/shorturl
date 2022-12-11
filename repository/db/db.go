package db

import (
	"context"
	"fmt"
	"github.com/22Fariz22/shorturl/model"
	"log"
	"time"

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
		_, err := i.conn.Exec(ctx, "insert into urls (cookies,correlation_id, id, longurl) values($1,$2,$3,$4);",
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

	//_, err = conn.Exec(ctx, "create table if not exists urls(cookies text, correlation_id text, id text CONSTRAINT id_pk PRIMARY KEY UNIQUE, longurl text);")
	_, err = conn.Exec(ctx,
		"create table if not exists urls(id_pk SERIAL PRIMARY KEY, cookies TEXT, correlation_id TEXT, short_url TEXT, long_url TEXT UNIQUE);")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (i *inDBRepository) SaveURL(ctx context.Context, shortID string, longURL string, cook string) (string, error) {

	var s string
	_ = i.conn.QueryRow(ctx, `
				WITH e AS(
				INSERT INTO urls (cookies, short_url, long_url) 
					   VALUES ($1, $2, $3)
				ON CONFLICT("long_url") DO NOTHING
				RETURNING long_url
			)
			SELECT long_url FROM e
			Union 
			SELECT short_url FROM urls where long_url=$3
;`, cook, shortID, longURL).Scan(&s)
	//if err != nil {
	//	log.Println(err)
	//	return "", err
	//}
	//fmt.Println("err:",err)
	if s != longURL {
		fmt.Println("такой есть. longurl: ", longURL, " s:", s)
		return s, nil

	}

	fmt.Println("такого нету. longurl: ", longURL, " s:", s)
	_, err := i.conn.Exec(ctx, "insert into urls (cookies, short_url, long_url) values($1,$2,$3);", cook, shortID, longURL)
	fmt.Println("err:", err)

	//fmt.Println("err in db after insert:\n", err)
	//if err != nil {
	//	log.Println(err)
	//	return "", err
	//}
	return "", err
}

//fmt.Println("s:\n", s)
//return "", nil

func (i *inDBRepository) GetURL(ctx context.Context, shortID string, cook string) (string, bool) {
	var s string
	err := i.conn.QueryRow(ctx, "select long_url from urls where short_url = $1 and cookies=$2 ;", shortID, cook).Scan(&s)
	if err != nil {
		log.Println(err)
		//TODO сделать возврат ошибки
		return "", false
	}
	return s, true
}

//example [map[7PJPPAZ:http://ya.ru] map[JRK5X81:http://ya.ru]]
func (i *inDBRepository) GetAll(ctx context.Context, cook string) ([]map[string]string, error) {
	rows, err := i.conn.Query(ctx, "select short_url, long_url from urls where cookies = $1;", cook)
	if err != nil {
		//log.Println(err)
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
