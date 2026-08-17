package tools

import "testing"

func TestMaskPhoneNumber(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "len 4",
			args: args{number: "1234"},
			want: "*****",
		},
		{
			name: "len 5",
			args: args{number: "12345"},
			want: "*****",
		},
		{
			name: "len 6",
			args: args{number: "123456"},
			want: "1*****",
		},
		{
			name: "len 11",
			args: args{number: "12345678901"},
			want: "12*****8901",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskPhoneNumber(tt.args.number); got != tt.want {
				t.Errorf("MaskPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskCardNumber(t *testing.T) {
	type args struct {
		cardNumber string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "len 0", args: args{cardNumber: ""}, want: ""},
		{name: "len 1", args: args{cardNumber: "1"}, want: "1"},
		{name: "len 5", args: args{cardNumber: "12345"}, want: "**345"},
		{name: "len 6", args: args{cardNumber: "123456"}, want: "***456"},
		{name: "len > 6", args: args{cardNumber: "123456789"}, want: "***456789"},
		{name: "chinese", args: args{cardNumber: "哇哈哈哇哈哈哇哈哈"}, want: "***哇哈哈哇哈哈"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskCardNumber(tt.args.cardNumber); got != tt.want {
				t.Errorf("MaskCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskName(tt.args.name); got != tt.want {
				t.Errorf("MaskName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskPhoneNumber1(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskPhoneNumber(tt.args.number); got != tt.want {
				t.Errorf("MaskPhoneNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskQQ(t *testing.T) {
	type args struct {
		qq string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskQQ(tt.args.qq); got != tt.want {
				t.Errorf("MaskQQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maskTimes(t *testing.T) {
	type args struct {
		nameLen      int
		maskLenLimit int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskTimes(tt.args.nameLen, tt.args.maskLenLimit); got != tt.want {
				t.Errorf("maskTimes() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestAppMaskCardNumber(t *testing.T) {
	type args struct {
		cardNumber string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "len 0", args: args{cardNumber: ""}, want: ""},
		{name: "len 1", args: args{cardNumber: "1"}, want: "1"},
		{name: "len 4", args: args{cardNumber: "1234"}, want: "**34"},
		{name: "len 5", args: args{cardNumber: "12345"}, want: "*****2345"},
		{name: "len 6", args: args{cardNumber: "123456"}, want: "*****3456"},
		{name: "len > 6", args: args{cardNumber: "123456789"}, want: "*****6789"},
		{name: "chinese", args: args{cardNumber: "哇哈哈哇哈哈哇哈哈"}, want: "*****哈哇哈哈"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppMaskCardNumber(tt.args.cardNumber); got != tt.want {
				t.Errorf("MaskCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppMaskName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "len 0", args: args{name: ""}, want: ""},
		{name: "len 1", args: args{name: "a"}, want: "a"},
		{name: "len 2", args: args{name: "ab"}, want: "*b"},
		{name: "len 3", args: args{name: "abc"}, want: "**c"},
		{name: "len 4", args: args{name: "abcd"}, want: "**d"},
		{name: "len 4ch", args: args{name: "一二三四"}, want: "**四"},
		{name: "len 7ch", args: args{name: "一二三四五六七"}, want: "**七"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppMaskName(tt.args.name); got != tt.want {
				t.Errorf("MaskName() = %v, want %v", got, tt.want)
			}
		})
	}
}