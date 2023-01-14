package mutation

import "testing"

func TestWriteInlineTemplate(t *testing.T) {
	type args struct {
		filename string
		tpl      string
		content  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteInlineTemplate(tt.args.filename, tt.args.tpl, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("WriteInlineTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
