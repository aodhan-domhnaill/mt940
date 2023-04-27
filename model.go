package mt940

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

var (
	ErrTagDoesNotApply = errors.New("tag doesn't apply to this struct")
	ErrNoTagsFound     = errors.New("no tags found")
)

type Transaction struct {
	TransactionReferenceNumber string
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

func (ts *Transactions) AddTag(t *Tag, r TagResults) *TagError {
	switch t.id {
	case "25":
		ts.AccountIdentification = r["account_identification"]
	case "28C":
		ts.StatementNumber = r["statement_number"]
		ts.StatementSeqNumber = r["sequence_number"]
	default:
		return &TagError{ErrTagDoesNotApply, t}
	}
	return nil
}

func (tr *Transaction) AddTag(t *Tag, r TagResults) *TagError {
	switch t.id {
	case "20":
		tr.TransactionReferenceNumber = r["transaction_reference"]
	default:
		return &TagError{ErrTagDoesNotApply, t}
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
