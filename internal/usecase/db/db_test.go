package db

import (
	"context"
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/handler"
	repoMock "github.com/22Fariz22/shorturl/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func _TestInDBRepositoryGetURL(t *testing.T) {
	ctl := gomock.NewController(t)
	ctl.Finish()

	repo := repoMock.NewMockRepository(ctl)

	ctx := context.Background()
	short := handler.GenUlid()

	mockResp := entity.URL{LongURL: "https://ya.ru"}

	expected := entity.URL{LongURL: "https://ya.ru"}

	repo.EXPECT().GetURL(ctx, short).Return(mockResp, true).Times(1)
	long, ok := repo.GetURL(ctx, short)

	require.Equal(t, expected, long)
	require.EqualValues(t, true, ok)
}

func _Test_inDBRepository_SaveURL(t *testing.T) {
	cfg := config.Config{}
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
	}

	type fields struct {
		pool        *pgxpool.Pool
		databaseDSN string
	}
	type args struct {
		ctx      context.Context
		shortURL string
		longURL  string
		cook     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "save in db",
			fields: fields{
				pool:        db,
				databaseDSN: "postgres://postgres:55555@127.0.0.1:5432/dburl",
			},
			args: args{
				ctx:      nil,
				shortURL: "some_short_url",
				longURL:  "https://ya.ru",
				cook:     "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inDBRepository{
				pool:        tt.fields.pool,
				databaseDSN: tt.fields.databaseDSN,
			}
			got, err := i.SaveURL(tt.args.ctx, tt.args.shortURL, tt.args.longURL, tt.args.cook)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SaveURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
