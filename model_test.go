package mt940

import (
	"testing"
	"time"
)

func TestAmount_Parse(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		amount  int64
		wantErr bool
	}{
		{
			"Basic",
			"123,23",
			12323,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amt := &Amount{}
			if err := amt.Parse(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Amount.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if amt.int64 != tt.amount {
				t.Errorf("Amount.Parse() got %v, wanted %v", amt.int64, tt.amount)
			}
		})
	}
}

func TestTransactionDate_Parse(t *testing.T) {

	type args struct {
		year  string
		month string
		day   string
	}
	tests := []struct {
		args    args
		t       time.Time
		wantErr bool
	}{
		{args{"70", "12", "01"}, time.Date(1970, 12, 1, 0, 0, 0, 0, time.UTC), false},
		{args{"20", "12", "01"}, time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC), false},
		{args{"50", "11", "01"}, time.Date(2050, 11, 1, 0, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run("date", func(t *testing.T) {
			td := &TransactionDate{}
			if err := td.Parse(tt.args.year, tt.args.month, tt.args.day); (err != nil) != tt.wantErr {
				t.Errorf("TransactionDate.Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if *td.Time != tt.t {
				t.Errorf("TransactionDate.Parse() got %v, want %v", td.Time, tt.t)
			}
		})
	}
}
