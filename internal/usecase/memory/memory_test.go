package memory

import (
	"context"
	"github.com/22Fariz22/shorturl/internal/entity"
	"github.com/22Fariz22/shorturl/internal/storage"
	"github.com/22Fariz22/shorturl/internal/usecase"
	"reflect"
	"testing"
)

//Test_inMemoryRepository_GetURL получить урл
func Test_inMemoryRepository_GetURL(t *testing.T) {
	type fields struct {
		memoryStorage storage.MemoryStorage
	}
	type args struct {
		ctx     context.Context
		shortID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   entity.URL
		want1  bool
	}{
		{
			name:   "получаем урл",
			fields: fields{memoryStorage: storage.New()},
			args: args{
				ctx:     nil,
				shortID: "some_short_url",
			},
			want:  entity.URL{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inMemoryRepository{
				memoryStorage: tt.fields.memoryStorage,
			}

			got, got1 := m.GetURL(tt.args.ctx, tt.args.shortID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

//Test_inMemoryRepository_SaveURL получить урлы
func Test_inMemoryRepository_SaveURL(t *testing.T) {
	type fields struct {
		memoryStorage storage.MemoryStorage
	}
	type args struct {
		ctx     context.Context
		shortID string
		longURL string
		cook    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "save url",
			fields: fields{memoryStorage: storage.New()},
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
			m := &inMemoryRepository{
				memoryStorage: tt.fields.memoryStorage,
			}
			got, err := m.SaveURL(tt.args.ctx, tt.args.shortID, tt.args.longURL, tt.args.cook)
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

//Test_inMemoryRepository_GetAll получить все урлы
func Test_inMemoryRepository_GetAll(t *testing.T) {
	type fields struct {
		memoryStorage storage.MemoryStorage
	}
	type args struct {
		ctx  context.Context
		cook string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]string
		wantErr bool
	}{
		{
			name:   "get all urls",
			fields: fields{memoryStorage: storage.New()},
			args: args{
				ctx:  nil,
				cook: "123456",
			},
			want:    make([]map[string]string, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inMemoryRepository{
				memoryStorage: tt.fields.memoryStorage,
			}
			got, err := m.GetAll(tt.args.ctx, tt.args.cook)
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

//Test_inMemoryRepository_RepoBatch
//signature: RepoBatch(ctx context.Context, cook string, batchList []entity.PackReq) error
func Test_inMemoryRepository_RepoBatch(t *testing.T) {
	type fields struct {
		memoryStorage storage.MemoryStorage
	}
	type args struct {
		ctx       context.Context
		cook      string
		batchList []entity.PackReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "create batch urls",
			fields: fields{memoryStorage: storage.New()},
			args: args{
				ctx:       nil,
				cook:      "1234567",
				batchList: make([]entity.PackReq, 0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inMemoryRepository{
				memoryStorage: tt.fields.memoryStorage,
			}
			tt.args.batchList = append(tt.args.batchList, entity.PackReq{
				CorrelationID: "1",
				OriginalURL:   "https://ya.ru",
				ShortURL:      "some_short_url",
			})
			tt.args.batchList = append(tt.args.batchList, entity.PackReq{
				CorrelationID: "2",
				OriginalURL:   "https://yahoo.ru",
				ShortURL:      "another_short_url",
			})

			if err := m.RepoBatch(tt.args.ctx, tt.args.cook, tt.args.batchList); (err != nil) != tt.wantErr {
				t.Errorf("RepoBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//Test_inMemoryRepository_Delete delete batch
func Test_inMemoryRepository_Delete(t *testing.T) {
	type fields struct {
		memoryStorage storage.MemoryStorage
	}
	type args struct {
		list   []string
		cookie string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "delete url",
			fields: fields{memoryStorage: storage.New()},
			args: args{
				list:   []string{"sdfasd", "lsdgjn"},
				cookie: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &inMemoryRepository{
				memoryStorage: tt.fields.memoryStorage,
			}
			if err := m.Delete(tt.args.list, tt.args.cookie); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//TestNew create New repo inmemory
func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want usecase.Repository
	}{
		{
			name: "create repo",
			want: New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
