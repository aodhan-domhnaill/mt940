package mt940

import "testing"

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
