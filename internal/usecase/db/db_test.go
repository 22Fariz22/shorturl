package db

import (
	"context"
	"encoding/hex"
	"github.com/22Fariz22/shorturl/internal/config"
	mock_usecase "github.com/22Fariz22/shorturl/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"log"
	"testing"
)

func Test_inDBRepository_SaveURLvarFirst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mock_usecase.NewMockRepository(ctrl)

	//value := []byte("Some value")

	ctx := context.Background()

	s.EXPECT().SaveURL(ctx, "12", "34", "56").Return("78", nil)

	s.SaveURL(ctx, "12", "34", "56")

	//type args struct {
	//	ctx      context.Context
	//	shortURL string
	//	longURL  string
	//	cook     string
	//}
	//tests := []struct {
	//	name    string
	//	fields  fields
	//	args    args
	//	want    string
	//	wantErr bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		i := &inDBRepository{
	//			pool:        tt.fields.pool,
	//			databaseDSN: tt.fields.databaseDSN,
	//		}
	//		got, err := i.SaveURL(tt.args.ctx, tt.args.shortURL, tt.args.longURL, tt.args.cook)
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("SaveURL() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if got != tt.want {
	//			t.Errorf("SaveURL() got = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func Test_inDBRepository_SaveURLvarSecond(t *testing.T) {
	type mockBehavior func(s *mock_usecase.MockRepositoryMockRecorder)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		expect       string
		error        error
	}{
		{
			name: "normal save",
			mockBehavior: func(s *mock_usecase.MockRepositoryMockRecorder) {
				s.SaveURL(context.Background(), "", "", "").Return("", nil)
			},
			expect: "",
			error:  nil,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

		})
	}
}

func Test_inDBRepository_SaveURL(t *testing.T) {
	// встраиваем мок-объекты вместо интерфейса, чтобы установить ожидания
	type fields struct {
		repo *mock_usecase.MockRepository
	}
	//type args struct {
	//	justNice bool
	//}
	tests := []struct {
		name string
		// «prepare» позволяет инициализировать наши моки в рамках конкретного теста
		prepare func(f *fields)
		//args    args
		wantErr bool
	}{
		{
			name: "some save",
			prepare: func(f *fields) {
				f.repo.EXPECT().SaveURL(context.Background(), "", "", "").Return("", nil)
			},
			//args:    args{justNice: false},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_usecase.NewMockRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}
			secretKey, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
			if err != nil {
				log.Fatal(err)
			}
			s := New(&config.Config{
				ServerAddress: "localhost:8080",
				BaseURL:       "http://localhost:8080",
				SecretKey:     secretKey,
				DatabaseDSN:   "postgres://postgres:55555@127.0.0.1:5432/dburl",
			})

			if _, err := s.SaveURL(context.Background(), "", "", ""); (err != nil) != tt.wantErr {
				t.Errorf("GreetVisitors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
