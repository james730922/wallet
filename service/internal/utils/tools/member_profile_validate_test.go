package tools

import "testing"

func Test_memberProfileValidate_BankAccount(t *testing.T) {
	type args struct {
		bankAccount string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正確",
			args: args{
				bankAccount: "9558889550942166",
			},
			wantErr: false,
		},
		{
			name: "錯誤",
			args: args{
				bankAccount: "955888955094216",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemberProfileValidate
			if err := m.BankCardNumber(tt.args.bankAccount); (err != nil) != tt.wantErr {
				t.Errorf("BankAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memberProfileValidate_Email(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正確",
			args: args{
				email: "simon.666@gmail.com",
			},
			wantErr: false,
		},
		{
			name: "錯誤",
			args: args{
				email: "asdflal@@al",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemberProfileValidate
			if err := m.Email(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memberProfileValidate_Mobile(t *testing.T) {
	type args struct {
		countryCode string
		mobile      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正確",
			args: args{
				countryCode: "+86",
				mobile:      "15125683493",
			},
			wantErr: false,
		},
		{
			name: "錯誤-countryCode",
			args: args{
				countryCode: "866",
				mobile:      "15125683493",
			},
			wantErr: true,
		},
		{
			name: "錯誤-mobile",
			args: args{
				countryCode: "+86",
				mobile:      "1512568349",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemberProfileValidate
			if err := m.Mobile(tt.args.countryCode, tt.args.mobile); (err != nil) != tt.wantErr {
				t.Errorf("Mobile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memberProfileValidate_Name(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正確",
			args: args{
				name: "賽門",
			},
			wantErr: false,
		},
		{
			name: "正確-有點",
			args: args{
				name: "剛·賽門",
			},
			wantErr: false,
		},
		{
			name: "錯誤",
			args: args{
				name: "賽門123",
			},
			wantErr: true,
		},
		{
			name: "錯誤-有點",
			args: args{
				name: "賽門·123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemberProfileValidate
			if err := m.Name(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Name() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memberProfileValidate_Passwd(t *testing.T) {
	type args struct {
		passwd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正確",
			args: args{
				passwd: "1joj2o4jgnv2",
			},
			wantErr: false,
		},
		{
			name: "錯誤",
			args: args{
				passwd: "賽門2k2nrk",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemberProfileValidate
			if err := m.Passwd(tt.args.passwd); (err != nil) != tt.wantErr {
				t.Errorf("Passwd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memberProfileValidate_QQ(t *testing.T) {
	type args struct {
		qq string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "正確",
			args: args{
				qq: "1244079287597",
			},
			wantErr: false,
		},
		{
			name: "錯誤-太短",
			args: args{
				qq: "123",
			},
			wantErr: true,
		},
		{
			name: "錯誤-開頭為0",
			args: args{
				qq: "0123245",
			},
			wantErr: true,
		},
		{
			name: "錯誤-其它字元",
			args: args{
				qq: "0123245阿阿",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemberProfileValidate
			if err := m.QQ(tt.args.qq); (err != nil) != tt.wantErr {
				t.Errorf("QQ() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
