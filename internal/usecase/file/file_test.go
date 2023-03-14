package file

import (
	"bufio"
	"context"
	"github.com/22Fariz22/shorturl/internal/storage"
	"io"
	"testing"
)

func Test_inFileRepository_SaveURL(t *testing.T) {
	type fields struct {
		file          io.ReadWriteCloser
		memoryStorage storage.MemoryStorage
		reader        *bufio.Reader
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
		// TODO: Add test cases.
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
