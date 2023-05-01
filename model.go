package mt940

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/text/currency"
)

var (
	ErrTagDoesNotApply   = errors.New("tag doesn't apply to this struct")
	ErrNoTagsFound       = errors.New("no tags found")
	ErrTagResultsMissing = errors.New("missing expected tag results fields")
)

type TransactionDate struct {
	*time.Time
}

type Balance struct {
	Timestamp TransactionDate
	Status    string
	Amount    float64 // TODO make fixed point
	Currency  currency.Unit
}

type StatementLine struct {
	Timestamp         TransactionDate
	EntryTime         TransactionDate
	Status            string
	FundsCode         string
	TransactionTypeID string
	CustomerReference string
	BankReference     string
	ExtraDetails      string
	Amount            float64 // TODO make fixed point
}

type Transaction struct {
	StatementLine
	TransactionReferenceNumber string
	FinalOpeningBalance        Balance
	FinalClosingBalance        Balance
	TransactionDetails         string
}

type Transactions struct {
	transactions          []Transaction
	AccountIdentification string
	StatementNumber       string
	StatementSeqNumber    string
}

type TagParser interface {
	AddTag(t *Tag, r TagResults) *TagError
}

func (td *TransactionDate) Parse(year, month, day string) error {
	if len(year) == 2 && year <= "69" {
		year = "20" + year
	} else if len(year) == 2 {
		year = "19" + year
	}
	t, err := time.Parse("2006-01-02", year+"-"+month+"-"+day)
	if err != nil {
		return err
	}
	td.Time = &t
	return nil
}

func (b *Balance) AddTag(t *Tag, r TagResults) *TagError {
	if err := b.Timestamp.Parse(r["year"], r["month"], r["day"]); err != nil {
		return &TagError{err, t, ""}
	}

	return nil
}

func (sl *StatementLine) AddTag(t *Tag, r TagResults) *TagError {
	if err := sl.Timestamp.Parse(r["year"], r["month"], r["day"]); err != nil {
		return &TagError{err, t, ""}
	}
	if err := sl.EntryTime.Parse(r["year"], r["entry_month"], r["entry_day"]); err != nil {
		return &TagError{err, t, ""}
	}
	sl.Status = r["status"]
	sl.FundsCode = r["funds_code"]
	_, err := fmt.Sscanf(
		strings.Replace(r["amount"], ",", ".", -1),
		"%f", &sl.Amount)

	if err != nil {
		return &TagError{err, t, ""}
	}

	sl.TransactionTypeID = r["id"]
	sl.CustomerReference = r["customer_reference"]
	sl.BankReference = r["bank_reference"]
	sl.ExtraDetails = r["extra_details"]

	return nil
}

func (ts *Transactions) AddTag(t *Tag, r TagResults) *TagError {
	switch t.id {
	case "25":
		ts.AccountIdentification = r["account_identification"]
	case "28C":
		ts.StatementNumber = r["statement_number"]
		ts.StatementSeqNumber = r["sequence_number"]
	default:
		return &TagError{ErrTagDoesNotApply, t, ""}
	}
	return nil
}

func (tr *Transaction) AddTag(t *Tag, r TagResults) *TagError {
	switch t.id {
	case "20":
		tr.TransactionReferenceNumber = r["transaction_reference"]
	case "60F":
		return tr.FinalOpeningBalance.AddTag(t, r)
	case "61":
		tr.StatementLine.AddTag(t, r)
	case "62F":
		return tr.FinalClosingBalance.AddTag(t, r)
	case "86":
		tr.TransactionDetails = r["transaction_details"]
	default:
		return &TagError{ErrTagDoesNotApply, t, ""}
	}
	return nil
}

func (t *Transactions) Parse(input io.Reader) ([]Transaction, error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	tagIndexes := tagRegex.FindAllIndex(data, -1)
	if len(tagIndexes) == 0 {
		return nil, ErrNoTagsFound
	}
	tr := &Transaction{}
	for i, inds := range tagIndexes {
		start := tagIndexes[i][0]
		end := len(data)
		if i+1 < len(tagIndexes) {
			end = tagIndexes[i+1][0]
		}
		block := data[start:end]
		// strip : off beginning and end
		id := string(data[inds[0]+1 : inds[1]-1])
		tag, ok := Tags[id]
		if !ok {
			return nil, fmt.Errorf("%v for tag %v", ErrNotExist, id)
		}

		result, err := tag.Parse(string(block))
		if err != nil {
			return nil, err
		}

		{
			parsers := []TagParser{
				tr, t,
			}

			var err *TagError
			for _, p := range parsers {
				err = p.AddTag(&tag, result)
				if err != nil && err.error == ErrTagDoesNotApply {
					continue
				} else {
					break
				}
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return t.transactions, nil
}
