package exe

/*
func TestRun(t *testing.T) {
	type args struct {
		c       string
		useSudo bool
	}
	execVal, err := exec.Command("dir").CombinedOutput()
	isError := err != nil
	tests := []struct {
		name string
		args args
		want *Execution
	}{
		{
			name: "generic",
			args: args{c: "dir", useSudo: false},
			want: &Execution{
				Input:       "dir",
				Fields:      []string{"dir"},
				UseSudo:     false,
				OutputStr:   string(execVal),
				ErrorStr:    err.Error(),
				IsError:     isError,
				Command:     "dir",
				Args:        []string{},
				regex:       nil,
				regexString: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.args.c, tt.args.useSudo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}*/
