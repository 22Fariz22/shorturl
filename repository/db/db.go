package db

import (
	"context"
	"log"

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

func (i *inDBRepository) SaveURL(ctx context.Context, shortID string, longURL string, cook string) error {
	_, err := i.conn.Exec(ctx, "create table if not exists urls(cookies text, id text,longurl text);")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = i.conn.Exec(ctx, "insert into urls (cookies, id, longurl) values($1,$2,$3);", cook, shortID, longURL)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *inDBRepository) GetURL(ctx context.Context, shortID string) (string, bool) {
	_, err := i.conn.Exec(ctx, "create table if not exists urls(cookies text, id text,longurl text);")
	if err != nil {
		log.Println(err)
		return "", false
	}

	var s string

	err = i.conn.QueryRow(ctx, "select longurl from urls where id = $1;", shortID).Scan(&s)
	if err != nil {
		log.Println(err)
		return "", false
	}
	return s, true
}

//example [map[7PJPPAZ:http://ya.ru] map[JRK5X81:http://ya.ru]]
func (i *inDBRepository) GetAll(ctx context.Context, cook string) []map[string]string {
	_, err := i.conn.Exec(ctx, "create table if not exists urls(cookies text, id text,longurl text);")
	if err != nil {
		log.Println(err)
	}

	rows, err := i.conn.Query(ctx, "select id, longurl from urls where cookies = $1;", cook)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	list := make([]map[string]string, 0)

	for rows.Next() {
		var id string
		var longurl string
		err = rows.Scan(&id, &longurl)
		if err != nil {
			log.Println(err)
			return nil
		}

		aMap := map[string]string{id: longurl}

		list = append(list, aMap)

	}
	return list
}

func (i *inDBRepository) Init() error {
	conn, err := pgx.Connect(context.Background(), i.databaseDSN)
	if err != nil {
		return err
	}
	i.conn = conn

	return nil
}

func (i *inDBRepository) Ping(ctx context.Context) error {
	//
	//conn, err := pgx.Connect(context.Background(), i.databaseDSN)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	return i.conn.Ping(ctx)
}
