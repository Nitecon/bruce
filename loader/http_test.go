package loader

import (
	"io"
	"reflect"
	"testing"
)

func TestReaderFromHttp(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    io.ReadCloser
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ReaderFromHttp(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReaderFromHttp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReaderFromHttp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ReaderFromHttp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
