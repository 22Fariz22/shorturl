package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"

	"github.com/22Fariz22/shorturl/model"

	"github.com/22Fariz22/shorturl/config"
	"github.com/22Fariz22/shorturl/repository"
	"github.com/jackc/pgx/v5"
)

type inDBRepository struct {
	conn        *pgx.Conn
	pool        *pgxpool.Conn
	databaseDSN string
	buffer      []model.PackResponse
	ctx         context.Context
}

func New(cfg *config.Config) repository.Repository {
	return &inDBRepository{
		databaseDSN: cfg.DatabaseDSN,
		buffer:      make([]model.PackResponse, 0, 1000),
	}
}

func (i *inDBRepository) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, i.databaseDSN)
	if err != nil {
		return err
	}
	i.conn = conn

	_, err = conn.Exec(ctx,
		"CREATE TABLE if not exists urls(id_pk SERIAL PRIMARY KEY, cookies TEXT, correlation_id TEXT,"+
			" short_url TEXT, long_url TEXT , deleted boolean default false);")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *inDBRepository) Delete(list []string, cookie string) error {
	fmt.Println("cookie in db del", cookie)
	log.Println("begin Delete")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := i.conn.Begin(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Prepare(ctx,
		"UPDATE", "UPDATE urls SET deleted = true WHERE short_url = $1 and cookies=$2;")
	log.Println("after prepare")
	if err != nil {
		log.Println("log in db del(2):", err)
		return err
	}

	for i := range list {
		log.Println("before Exec")
		_, err = tx.Exec(ctx,
			"UPDATE urls SET deleted = true WHERE short_url = $1 and cookies=$2;",
			list[i], cookie)
		log.Println("after exec")
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = tx.Commit(ctx)
	log.Println("after commit")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *inDBRepository) SaveURL(ctx context.Context, shortURL string, longURL string, cook string) (string, error) {
	var ErrAlreadyExists = errors.New("this URL already exists")
	var id int8
	fmt.Println("cookie in db saveurl", cook)
	fmt.Println("lu in db SaveURL", longURL)
	// проверяем количество строк, если есть то значит такой урл существует
	row := i.conn.QueryRow(ctx, `select count(*) from urls where long_url = $1 and cookies=$2`,
		longURL, cook)
	err := row.Scan(&id)
	if err != nil {
		log.Println("log in db SaveURL(0):", err)
	}
	if id == 0 {
		// добавляем новую запись
		_, err = i.conn.Exec(ctx, "insert into urls (cookies, short_url, long_url) values($1,$2,$3);",
			cook, shortURL, longURL)
		if err != nil {
			log.Println("log in db SaveURL(1):", err)
			return "", err
		}
		return "", nil
	}

	// делаем запрос на существующую запись и выдаем шортурл
	var su string

	err = i.conn.QueryRow(ctx, "select short_url from urls where long_url = $1 and cookies = $2;",
		longURL, cook).Scan(&su)
	if err != nil {
		log.Println("log in SaveURL(2):", err)
		return "", err
	}
	fmt.Println("id", id)
	fmt.Println("short_url:", su)
	return su, ErrAlreadyExists
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
