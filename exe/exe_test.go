package exe

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestExecution_ContainsLC(t *testing.T) {
	type fields struct {
		input       string
		fields      []string
		useSudo     bool
		outputStr   string
		isError     bool
		cmnd        string
		args        []string
		regex       *regexp.Regexp
		regexString string
		err         error
	}
	type args struct {
		c string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "hasOutputMatch",
			want:   true,
			args:   args{c: "hello"},
			fields: fields{outputStr: "hello"},
		},
		{
			name:   "hasErrorMatch",
			want:   true,
			args:   args{c: "hello"},
			fields: fields{err: fmt.Errorf("hello")},
		},
		{
			name:   "hasNoMatch",
			want:   false,
			args:   args{c: "fubar"},
			fields: fields{err: fmt.Errorf("hello")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Execution{
				input:       tt.fields.input,
				fields:      tt.fields.fields,
				useSudo:     tt.fields.useSudo,
				outputStr:   tt.fields.outputStr,
				isError:     tt.fields.isError,
				cmnd:        tt.fields.cmnd,
				args:        tt.fields.args,
				regex:       tt.fields.regex,
				regexString: tt.fields.regexString,
				err:         tt.fields.err,
			}
			if got := e.ContainsLC(tt.args.c); got != tt.want {
				t.Errorf("ContainsLC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecution_GetErr(t *testing.T) {
	type fields struct {
		input       string
		fields      []string
		useSudo     bool
		outputStr   string
		isError     bool
		cmnd        string
		args        []string
		regex       *regexp.Regexp
		regexString string
		err         error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "testIfErrSet",
			fields:  fields{err: fmt.Errorf("testme")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Execution{
				input:       tt.fields.input,
				fields:      tt.fields.fields,
				useSudo:     tt.fields.useSudo,
				outputStr:   tt.fields.outputStr,
				isError:     tt.fields.isError,
				cmnd:        tt.fields.cmnd,
				args:        tt.fields.args,
				regex:       tt.fields.regex,
				regexString: tt.fields.regexString,
				err:         tt.fields.err,
			}
			if err := e.GetErr(); (err != nil) != tt.wantErr {
				t.Errorf("GetErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecution_SetRegex(t *testing.T) {
	type fields struct {
		input       string
		fields      []string
		useSudo     bool
		outputStr   string
		isError     bool
		cmnd        string
		args        []string
		regex       *regexp.Regexp
		regexString string
		err         error
	}
	type args struct {
		re string
	}
	regx, _ := regexp.Compile("p([a-z]+)ch")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *regexp.Regexp
		wantErr bool
	}{
		{
			name:    "regexSuccess",
			fields:  fields{},
			wantErr: false,
			args:    args{re: "p([a-z]+)ch"},
			want:    regx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Execution{
				input:       tt.fields.input,
				fields:      tt.fields.fields,
				useSudo:     tt.fields.useSudo,
				outputStr:   tt.fields.outputStr,
				isError:     tt.fields.isError,
				cmnd:        tt.fields.cmnd,
				args:        tt.fields.args,
				regex:       tt.fields.regex,
				regexString: tt.fields.regexString,
				err:         tt.fields.err,
			}
			got, err := e.SetRegex(tt.args.re)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetRegex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecution_RegexMatch(t *testing.T) {
	type fields struct {
		input       string
		fields      []string
		useSudo     bool
		outputStr   string
		isError     bool
		cmnd        string
		args        []string
		regex       *regexp.Regexp
		regexString string
		err         error
	}
	r, _ := regexp.Compile("p([a-z]+)ch")
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "outputSuccess",
			fields: fields{outputStr: "peachy", regex: r},
			want:   true,
		},
		{
			name:   "regexNotCompiled",
			fields: fields{outputStr: "peachy", regex: nil},
			want:   false,
		},
		{
			name:   "errorMatched",
			fields: fields{err: fmt.Errorf("peachy"), regex: r},
			want:   true,
		},
		{
			name:   "noMatchCompile",
			fields: fields{err: fmt.Errorf("channels"), regex: r},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Execution{
				input:       tt.fields.input,
				fields:      tt.fields.fields,
				useSudo:     tt.fields.useSudo,
				outputStr:   tt.fields.outputStr,
				isError:     tt.fields.isError,
				cmnd:        tt.fields.cmnd,
				args:        tt.fields.args,
				regex:       tt.fields.regex,
				regexString: tt.fields.regexString,
				err:         tt.fields.err,
			}
			if got := e.RegexMatch(); got != tt.want {
				t.Errorf("RegexMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecution_Failed(t *testing.T) {
	type fields struct {
		input       string
		fields      []string
		useSudo     bool
		outputStr   string
		isError     bool
		cmnd        string
		args        []string
		regex       *regexp.Regexp
		regexString string
		err         error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "success",
			want:   true,
			fields: fields{err: fmt.Errorf("hi"), isError: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Execution{
				input:       tt.fields.input,
				fields:      tt.fields.fields,
				useSudo:     tt.fields.useSudo,
				outputStr:   tt.fields.outputStr,
				isError:     tt.fields.isError,
				cmnd:        tt.fields.cmnd,
				args:        tt.fields.args,
				regex:       tt.fields.regex,
				regexString: tt.fields.regexString,
				err:         tt.fields.err,
			}
			if got := e.Failed(); got != tt.want {
				t.Errorf("Failed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		c       string
		useSudo bool
	}
	tests := []struct {
		name string
		args args
		want *Execution
	}{
		{
			name: "success",
			args: args{c: "dir", useSudo: false},
			want: &Execution{
				input:   "dir",
				fields:  []string{"dir"},
				isError: true,
				cmnd:    "dir",
				args:    []string{},
				err:     errors.New(`exec: "dir": executable file not found in %PATH%`),
			},
		},
		{
			name: "withSudo",
			args: args{c: "dir", useSudo: true},
			want: &Execution{
				input:   "dir",
				fields:  []string{"dir"},
				isError: true,
				cmnd:    "sudo",
				useSudo: true,
				args:    []string{"dir"},
				err:     errors.New(`exec: "sudo": executable file not found in %PATH%`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.args.c, tt.args.useSudo); !reflect.DeepEqual(got) {
				t.Errorf("Function Run()\n Got=%#v\n want%#v\nError: %s", got, tt.want)
			}
		})
	}
}

func TestHasExecInPath(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "linux",
			args: args{name: "foo"},
			want: "INFO: Could not find files for the given pattern(s).\r",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasExecInPath(tt.args.name); got != tt.want {
				t.Errorf("HasExecInPath() = \ngots %#v\nwant %#v", got, tt.want)
			}
		})
	}
}
