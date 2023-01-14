package mutation

import (
	"io"
	"reflect"
	"testing"
)

func Test_useGzipReader(t *testing.T) {
	type args struct {
		filename   string
		fileReader io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want io.ReadCloser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := useGzipReader(tt.args.filename, tt.args.fileReader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("useGzipReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
