package mt940

import (
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"golang.org/x/text/currency"
)

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func newTransactionDate(t time.Time) TransactionDate {
	return TransactionDate{Time: &t}
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
			want: []Transaction{
				Transaction{
					StatementLine: StatementLine{
						Timestamp: newTransactionDate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)), EntryTime: newTransactionDate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)), Status: "D", FundsCode: "", TransactionTypeID: "NOVB", CustomerReference: "NL47INGB9999999999\nhr gjlm pauli", BankReference: "", ExtraDetails: "ssen", Amount: Amount{65}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR}, TransactionDetails: "NL47INGB9999999999 hr gjlm paulissen\n                                                                 \nBetaling sieraden"},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 4, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Timestamp: newTransactionDate(time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC)), EntryTime: newTransactionDate(time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC)), Status: "D", FundsCode: "", TransactionTypeID: "NIDB", CustomerReference: "NL08ABNA9999999999\ninternational", BankReference: "", ExtraDetails: " card services", Amount: Amount{80155}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 5, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR}, TransactionDetails: "NL08ABNA9999999999 international card services \n                                                                 \n000000000000000000000000000000000 0000000000000000 Betaling aan I\nCS 99999999999 ICS Referentie: 2020-01-05 19:47 000000000000000"},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 7, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 7, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 8, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 8, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 9, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 9, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 10, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 10, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 11, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 11, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 12, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 12, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 13, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 13, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 14, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 14, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 16, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 16, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 17, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 17, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 18, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 18, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 19, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 19, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 20, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 20, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 21, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 21, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 22, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 22, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 23, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 23, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 24, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 24, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Timestamp: newTransactionDate(time.Date(2020, time.January, 25, 0, 0, 0, 0, time.UTC)), EntryTime: newTransactionDate(time.Date(2020, time.January, 25, 0, 0, 0, 0, time.UTC)), Status: "D", FundsCode: "", TransactionTypeID: "NDIV", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{165}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 25, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 25, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR}, TransactionDetails: "Kosten gebruik betaalrekening inclusief 1 betaalpas"},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 26, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 26, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 27, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 27, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 28, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 28, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
				Transaction{
					StatementLine: StatementLine{
						Timestamp: newTransactionDate(time.Date(2020, time.January, 29, 0, 0, 0, 0, time.UTC)), EntryTime: newTransactionDate(time.Date(2020, time.January, 29, 0, 0, 0, 0, time.UTC)), Status: "D", FundsCode: "", TransactionTypeID: "NIDB", CustomerReference: "NL08ABNA9999999999\ninternational", BankReference: "", ExtraDetails: " card services", Amount: Amount{100000}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 29, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 29, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR}, TransactionDetails: "NL08ABNA9999999999 international card services \n                                                                 \n000000000000000000000000000000000 0000000000000000 Betaling aan I\nCS 99999999999 ICS Referentie: 2020-01-29 18:36 000000000000000"},
				Transaction{
					StatementLine: StatementLine{
						Status: "", FundsCode: "", TransactionTypeID: "", CustomerReference: "", BankReference: "", ExtraDetails: "", Amount: Amount{0}},
					TransactionReferenceNumber: "0000000000",
					FinalOpeningBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 30, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					FinalClosingBalance:        Balance{Timestamp: newTransactionDate(time.Date(2020, time.January, 30, 0, 0, 0, 0, time.UTC)), Status: "", Amount: Amount{0}, Currency: currency.EUR},
					TransactionDetails:         "",
				},
			},
			wantErr: false,
		},
		{
			name: "mBank",
			args: args{
				input: must(os.Open("mBank/mt940.sta")),
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
				t.Error(err)

				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Transactions.Parse() len(results) = %v, want %v", len(got), len(tt.want))
				t.Errorf("%#v", got)
			}
			for i, trans := range got {
				if reflect.DeepEqual(trans, tt.want[i]) {
					t.Errorf("Transactions not equal: got[%v] = %v, want[%v] = %v", i, trans, i, tt.want[i])
				}
			}
		})
	}
}
