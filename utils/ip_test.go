package utils

import "testing"

func TestIncrIP(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "IncrIP", args: args{ip: "10.9.0.1/24"}, want: "10.9.0.2/24", wantErr: false},
		{name: "IncrIPErrorOverflow", args: args{ip: "10.255.255.255/24"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IncrIP(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IncrIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IncrIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
