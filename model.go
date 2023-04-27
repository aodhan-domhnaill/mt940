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
}

type Transactions struct {
	transactions []Transaction
}

func (ts *Transactions) AddTag(t *Tag, r *TagResults) error {
	switch t.id {
	default:
		return ErrTagDoesNotApply
	}
}

func (tr *Transaction) AddTag(t *Tag, r *TagResults) error {
	switch t.id {
	default:
		return ErrTagDoesNotApply
	}
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
	tr := Transaction{}
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

		if err := tr.AddTag(&tag, &result); err == ErrTagDoesNotApply {
			if err := t.AddTag(&tag, &result); err != nil {
				return nil, fmt.Errorf("%v for tag %v", err, tag.id)
			}
		} else if err != nil {
			return nil, err
		}
	}

	return t.transactions, nil
}
