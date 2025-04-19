package wireguard

// import (
// 	"os"
// 	"testing"
// )

// func Test_createWgInterface(t *testing.T) {
// 	type args struct {
// 		f    *os.File
// 		priv string
// 		addr string
// 		port string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{name: "createWgInterface", args: args{f: nil, priv: "priv", addr: "addr", port: "port"}, wantErr: false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := createWgInterface(tt.args.f, tt.args.priv, tt.args.addr, tt.args.port); (err != nil) != tt.wantErr {
// 				t.Errorf("createWgInterface() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func Test_createWgPeer(t *testing.T) {
// 	type args struct {
// 		f   *os.File
// 		pub string
// 		ip  string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{name: "createWgPeer", args: args{pub: "pub", ip: "ip", f: nil}, wantErr: false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := createWgPeer(tt.args.f, tt.args.pub, tt.args.ip); (err != nil) != tt.wantErr {
// 				t.Errorf("createWgPeer() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
