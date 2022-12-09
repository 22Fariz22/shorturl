package db

import (
	"context"
	"github.com/22Fariz22/shorturl/model"
	"log"
	"time"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/jackc/pgx/v5"
)

type db interface {
	AddURL(*model.PackResponse)
	Flush(string, string) error
	repository.Repository
}

type inDBRepository struct {
	conn        *pgx.Conn
	databaseDSN string
	buffer      []model.PackResponse
	ctx         context.Context
}

//func (i *inDBRepository) AddURL(p *model.PackResponse) error {
//	i.buffer = append(i.buffer, *p)
//
//	if cap(i.buffer) == len(i.buffer) {
//		err := i.Flush()
//		if err != nil {
//			return errors.New("cannot add records to the database")
//		}
//	}
//	return nil
//}
//func (i *inDBRepository) Flush(cook string, shortUrl string) error {
//	if i.conn == nil {
//		return errors.New("You haven`t opened the database connection")
//	}
//	tx, err := i.conn.Begin(i.ctx)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback(i.ctx)
//
//	//defer stmt.Close()
//	for _, v := range i.buffer {
//		_, err := tx.Exec(i.ctx, "INSERT INTO urls(cookies,correlation_id, id, longurl) VALUES($1,$2,$3,$4)",
//			cook, v.Correlation_id, shortUrl, v.Original_url)
//		if err != nil {
//			return err
//		}
//	}
//
//	if err := tx.Commit(i.ctx); err != nil {
//		log.Println("update drivers: unable to commit: %v", err)
//		return err
//	}
//
//	i.buffer = i.buffer[:0]
//	return nil
//}

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

	_, err = conn.Exec(ctx, "create table if not exists urls("+
		"cookies text, correlation_id text, id text CONSTRAINT id_pk PRIMARY KEY, longurl text);")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (i *inDBRepository) SaveURL(ctx context.Context, shortID string, longURL string, cook string) error {
	_, err := i.conn.Exec(ctx, "insert into urls (cookies, id, longurl) values($1,$2,$3);", cook, shortID, longURL)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *inDBRepository) GetURL(ctx context.Context, shortID string) (string, bool) {
	var s string
	err := i.conn.QueryRow(ctx, "select longurl from urls where id = $1;", shortID).Scan(&s)
	if err != nil {
		log.Println(err)
		//TODO сделать возврат ошибки
		return "", false
	}

	return s, true
}

//example [map[7PJPPAZ:http://ya.ru] map[JRK5X81:http://ya.ru]]
func (i *inDBRepository) GetAll(ctx context.Context, cook string) ([]map[string]string, error) {
	rows, err := i.conn.Query(ctx, "select id, longurl from urls where cookies = $1;", cook)
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
