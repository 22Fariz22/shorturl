package file

import (
	"bufio"
	"context"
	"github.com/22Fariz22/shorturl/internal/config"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"os"
	"reflect"
	"testing"

	"github.com/22Fariz22/shorturl/internal/storage"
)

//Test_inFileRepository_SaveURL
//сигнатура: SaveURL(ctx context.Context, shortID string, longURL string, cook string) (string, error)
func Test_inFileRepository_SaveURL(t *testing.T) {
	type fields struct {
		inFileRepository
	}

	type args struct {
		ctx     context.Context
		shortID string
		longURL string
		cook    string
	}
	file, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0644)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "сохраняем в файл",

			fields: fields{inFileRepository{
				file:          file,
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			}},
			args: args{
				ctx:     context.Background(),
				shortID: "some_short_url",
				longURL: "https://ya.ru",
				cook:    "123456",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inFileRepository{
				file:          tt.fields.file,
				memoryStorage: tt.fields.memoryStorage,
				reader:        tt.fields.reader,
			}

			got, err := f.SaveURL(tt.args.ctx, tt.args.shortID, tt.args.longURL, tt.args.cook)
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

func Test_inFileRepository_GetURL(t *testing.T) {
	type fields struct {
		inFileRepository
	}
	type args struct {
		ctx     context.Context
		shortID string
	}
	file, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0644)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   entity.URL
		want1  bool
	}{
		{
			name: "get url not ok",
			fields: fields{inFileRepository{
				file:          file,
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			}},
			args: args{
				ctx:     context.Background(),
				shortID: "shorturl",
			},
			want:  entity.URL{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inFileRepository{
				file:          tt.fields.file,
				memoryStorage: tt.fields.memoryStorage,
				reader:        tt.fields.reader,
			}
			got, got1 := f.GetURL(tt.args.ctx, tt.args.shortID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_inFileRepository_GetAll(t *testing.T) {
	type fields struct {
		inFileRepository
	}
	type args struct {
		ctx  context.Context
		cook string
	}
	file, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0644)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]string
		wantErr bool
	}{
		{
			name: "empty list",
			fields: fields{inFileRepository{
				file:          file,
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			}},
			args: args{
				ctx:  context.Background(),
				cook: "123456",
			},
			want:    []map[string]string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inFileRepository{
				file:          tt.fields.file,
				memoryStorage: tt.fields.memoryStorage,
				reader:        tt.fields.reader,
			}
			got, err := f.GetAll(tt.args.ctx, tt.args.cook)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inFileRepository_Init(t *testing.T) {
	type fields struct {
		inFileRepository
	}
	file, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0644)

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "init ok",
			fields: fields{inFileRepository{
				file:          file,
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inFileRepository{
				file:          tt.fields.file,
				memoryStorage: tt.fields.memoryStorage,
				reader:        tt.fields.reader,
			}
			if err := f.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inFileRepository_RepoBatch(t *testing.T) {
	type fields struct {
		inFileRepository
	}
	type args struct {
		ctx       context.Context
		cook      string
		batchList []entity.PackReq
	}
	file, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0644)

	var batchList []entity.PackReq

	batchList = append(batchList, entity.PackReq{
		CorrelationID: "1",
		OriginalURL:   "bllabla",
		ShortURL:      "bllabla",
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "repo ok",
			fields: fields{inFileRepository{
				file:          file,
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			}},
			args: args{
				ctx:       context.Background(),
				cook:      "12345",
				batchList: batchList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inFileRepository{
				file:          tt.fields.file,
				memoryStorage: tt.fields.memoryStorage,
				reader:        tt.fields.reader,
			}
			if err := f.RepoBatch(tt.args.ctx, tt.args.cook, tt.args.batchList); (err != nil) != tt.wantErr {
				t.Errorf("RepoBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_inFileRepository_Delete(t *testing.T) {
	type fields struct {
		inFileRepository
	}
	type args struct {
		list   []string
		cookie string
	}
	file, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0644)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "del ok",
			fields: fields{inFileRepository{
				file:          file,
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			}},
			args: args{
				list:   []string{},
				cookie: "123456",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &inFileRepository{
				file:          tt.fields.file,
				memoryStorage: tt.fields.memoryStorage,
				reader:        tt.fields.reader,
			}
			if err := f.Delete(tt.args.list, tt.args.cookie); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewConsumer(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	cfg := config.Config{}
	cfg.FileStoragePath = "test"

	//var file os.File

	tests := []struct {
		name    string
		args    args
		want    *Consumer
		wantErr bool
	}{
		{
			name: "new consumer",
			args: args{cfg: cfg},
			want: &Consumer{
				File:   &os.File{},
				reader: new(bufio.Reader),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConsumer(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewConsumer() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

func _TestNew(t *testing.T) {
	type args struct {
		cfg *config.Config
	}

	cfg := config.NewConfig()

	//consumer, err := NewConsumer(*cfg)
	//if err != nil {
	//	log.Println(err)
	//}

	tests := []struct {
		name string
		args args
		want usecase.Repository
	}{
		{
			name: "new",
			args: args{cfg: cfg},
			want: &inFileRepository{
				file:          &os.File{},
				memoryStorage: storage.New(),
				reader:        new(bufio.Reader),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
