package storage

import (
	"github.com/22Fariz22/shorturl/internal/entity"
	"reflect"
	"testing"
)

// Test_memoryStorage_Insert тестируем стораж инмемори и случай когда такой урл не существует в мапе
//
//	Insert(key, value string, cook string, deleted bool) (string, error)
func Test_memoryStorage_Insert(t *testing.T) {
	type fields struct {
		storage map[string]entity.URL
		//mutex   sync.RWMutex
	}

	type args struct {
		key     string //shortID
		value   string //longURL
		cook    string //cook
		deleted bool   //delete status
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "long url not exist in map",
			fields: fields{
				storage: make(map[string]entity.URL, 1),
				//mutex:   sync.RWMutex{},
			},
			args: args{
				key:     "some-shorturl",
				value:   "https://ya.ru",
				cook:    "123456789",
				deleted: false,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryStorage{
				storage: tt.fields.storage,
				//mutex:   tt.fields.mutex,
			}

			got, err := m.Insert(tt.args.key, tt.args.value, tt.args.cook, tt.args.deleted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Insert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_memoryStorage_InsertIfExist тестируем стораж инмемори и случай когда такой урл уже существует
func Test_memoryStorage_InsertExist(t *testing.T) {
	type fields struct {
		storage map[string]entity.URL
		//mutex   sync.RWMutex
	}

	type args struct {
		key     string //shortID
		value   string //longURL
		cook    string //cook
		deleted bool   //delete status
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "long url exist in map",
			fields: fields{
				storage: make(map[string]entity.URL, 1),
				//mutex:   sync.RWMutex{},
			},
			args: args{
				key:     "some-shorturl",
				value:   "https://ya.ru",
				cook:    "123456789",
				deleted: false,
			},
			want:    "some-shorturl",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryStorage{
				storage: tt.fields.storage,
				//mutex:   tt.fields.mutex,
			}
			m.storage["https://ya.ru"] = entity.URL{
				Cookies:       "123456789",
				ID:            "some-shorturl",
				LongURL:       "https//:ya.ru",
				CorrelationID: "",
				Deleted:       false,
			}

			got, err := m.Insert(tt.args.key, tt.args.value, tt.args.cook, tt.args.deleted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Insert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_memoryStorage_Get  сигнатура:Get(key string) (entity.URL, bool)
// // когда вводим существующий шортурл
func Test_memoryStorage_Get(t *testing.T) {
	type fields struct {
		storage map[string]entity.URL
		//mutex   sync.RWMutex
	}
	type args struct {
		key string //shortID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   entity.URL
		want1  bool
	}{
		{
			name: "получаем существующий лонг урл",
			fields: fields{
				storage: make(map[string]entity.URL),
				//mutex:   sync.RWMutex{},
			},
			args: args{"some-short-url"},
			want: entity.URL{
				Cookies:       "123456",
				ID:            "some-short-url",
				LongURL:       "https://ya.ru",
				CorrelationID: "",
				Deleted:       false,
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryStorage{
				storage: tt.fields.storage,
				//mutex:   tt.fields.mutex,
			}
			m.storage["https://ya.ru"] = entity.URL{
				Cookies:       "123456",
				ID:            "some-short-url",
				LongURL:       "https://ya.ru",
				CorrelationID: "",
				Deleted:       false,
			}
			got, got1 := m.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

// Test_memoryStorage_GetNotExist  сигнатура:Get(key string) (entity.URL, bool)
// когда вводим не существующий шортурл
func Test_memoryStorage_GetNotExist(t *testing.T) {
	type fields struct {
		storage map[string]entity.URL
		//mutex   sync.RWMutex
	}
	type args struct {
		key string //shortID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   entity.URL
		want1  bool
	}{
		{
			name: "попытка ввести не существующий шортурл",
			fields: fields{
				storage: make(map[string]entity.URL),
				//mutex:   sync.RWMutex{},
			},
			args:  args{"some-short-url"},
			want:  entity.URL{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryStorage{
				storage: tt.fields.storage,
				//mutex:   tt.fields.mutex,
			}
			m.storage["https://google.com"] = entity.URL{
				Cookies:       "123456",
				ID:            "another-short-url",
				LongURL:       "https://google.com",
				CorrelationID: "",
				Deleted:       false,
			}
			got, got1 := m.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

// Test_memoryStorage_GetAllStorageURL сигнатура: GetAllStorageURL(cook string) []map[string]string
func Test_memoryStorage_GetAllStorageURL(t *testing.T) {
	type fields struct {
		storage map[string]entity.URL
		//mutex   sync.RWMutex
	}
	type args struct {
		cook string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []map[string]string
	}{
		{
			name: "получаем список своих урлов",
			fields: fields{
				storage: make(map[string]entity.URL, 0),
				//mutex:   sync.RWMutex{},
			},
			args: args{cook: "123456"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryStorage{
				storage: tt.fields.storage,
				//mutex:   tt.fields.mutex,
			}

			tt.want = append(tt.want, map[string]string{"another-short-url": "https://yahoo.com"})
			tt.want = append(tt.want, map[string]string{"some-short-url": "https://google.com"})

			m.storage["https://yahoo.com"] = entity.URL{
				Cookies:       "123456",
				ID:            "another-short-url",
				LongURL:       "https://yahoo.com",
				CorrelationID: "",
				Deleted:       false,
			}
			m.storage["https://google.com"] = entity.URL{
				Cookies:       "123456",
				ID:            "some-short-url",
				LongURL:       "https://google.com",
				CorrelationID: "",
				Deleted:       false,
			}
			if got := m.GetAllStorageURL(tt.args.cook); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllStorageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_memoryStorage_DeleteStorage сигнатура:DeleteStorage(listShorts []string, cookies string) error
func Test_memoryStorage_DeleteStorage(t *testing.T) {
	type fields struct {
		storage map[string]entity.URL
		//mutex   sync.RWMutex
	}
	type args struct {
		listShorts []string
		cookies    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "удаляем без ошибки",
			fields: fields{
				storage: map[string]entity.URL{},
				//mutex:   sync.RWMutex{},
			},
			args: args{
				listShorts: []string{"some_short_url", "another_short_url"},
				cookies:    "123456",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &memoryStorage{
				storage: tt.fields.storage,
				//mutex:   tt.fields.mutex,
			}

			m.storage["ya.ru"] = entity.URL{
				Cookies:       "123456",
				ID:            "some_short_url",
				LongURL:       "ya.ru",
				CorrelationID: "",
				Deleted:       false,
			}
			m.storage["yaho.ru"] = entity.URL{
				Cookies:       "123456",
				ID:            "another_short_url",
				LongURL:       "yaho.ru",
				CorrelationID: "",
				Deleted:       false,
			}

			if err := m.DeleteStorage(tt.args.listShorts, tt.args.cookies); (err != nil) != tt.wantErr {
				t.Errorf("DeleteStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
