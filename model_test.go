package mt940

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func TestTransactions_Parse(t *testing.T) {
	type args struct {
		input io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Transaction
		wantErr bool
	}{
		{
			name: "ASNB",
			args: args{
				input: must(os.Open("ASNB/0708271685_09022020_164516.940.txt")),
			},
			want:    []Transaction{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transactions{}
			got, err := tr.Parse(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transactions.Parse() error = %v for tag %v value '%v', wantErr %v",
					err, err.(*TagError).id, err.(*TagError).Value, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transactions.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
