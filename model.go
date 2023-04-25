package mt940

import (
	"bufio"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type Model struct{}

type Amount struct {
	amount   *big.Float
	currency string
}

func NewAmount(amount string, status string, currency string) *Amount {
	amt, _ := new(big.Float).SetString(strings.ReplaceAll(amount, ",", "."))
	if status == "D" {
		amt = new(big.Float).Neg(amt)
	}
	return &Amount{
		amount:   amt,
		currency: currency,
	}
}

func (a *Amount) String() string {
	return fmt.Sprintf("%s %s", a.amount.String(), a.currency)
}

func (a *Amount) Equals(other *Amount) bool {
	return a.amount.Cmp(other.amount) == 0 && a.currency == other.currency
}

type Balance struct {
	status string
	amount *Amount
	date   string
}

func NewBalance(status string, amount *Amount, date string) *Balance {
	return &Balance{
		status: status,
		amount: amount,
		date:   date,
	}
}

func (b *Balance) String() string {
	return fmt.Sprintf("%s @ %s", b.amount.String(), b.date)
}

func (b *Balance) Equals(other *Balance) bool {
	return b.status == other.status && b.amount.Equals(other.amount)
}

func (t *Transactions) Add(transaction *Transaction) {
	t.transactions = append(t.transactions, transaction)
}

func (t *Transactions) SetData(key string, value interface{}) {
	t.data[key] = value
}

func (t *Transactions) GetCurrency() string {
	var currency string
	balanceKeys := []string{
		"final_opening_balance",
		"opening_balance",
		"intermediate_opening_balance",
		"available_balance",
		"forward_available_balance",
		"final_closing_balance",
		"closing_balance",
		"intermediate_closing_balance",
		"c_floor_limit",
		"d_floor_limit",
	}
	for _, key := range balanceKeys {
		if val, ok := t.data[key]; ok {
			if bal, ok := val.(*Balance); ok {
				currency = bal.amount.currency
				break
			}
		}
	}
	return currency
}

type Transaction struct {
	transactionReference string
	relatedReference     string
	amount               *Amount
	date                 string
	entryDate            string
	valueDate            string
	bookCode             string
	details              string
	bankCode             string
	customerReference    string
	bankReference        string
	supplementaryDetails string
}

func NewTransaction(transactionReference string, relatedReference string, amount *Amount, date string,
	entryDate string, valueDate string, bookCode string, details string, bankCode string,
	customerReference string, bankReference string, supplementaryDetails string) *Transaction {
	return &Transaction{
		transactionReference: transactionReference,
		relatedReference:     relatedReference,
		amount:               amount,
		date:                 date,
		entryDate:            entryDate,
		valueDate:            valueDate,
		bookCode:             bookCode,
		details:              details,
		bankCode:             bankCode,
		customerReference:    customerReference,
		bankReference:        bankReference,
		supplementaryDetails: supplementaryDetails,
	}
}

type Transactions struct {
	processors   map[string][]interface{}
	tags         map[int]Tag
	transactions []*Transaction
	data         map[string]interface{}
}

func NewTransactions(processors map[string][]interface{}, tags map[int]Tag) *Transactions {
	t := &Transactions{
		processors:   defaultProcessors(),
		tags:         defaultTags(),
		transactions: make([]*Transaction, 0),
		data:         make(map[string]interface{}),
	}

	if processors != nil {
		t.processors = processors
	}

	if tags != nil {
		t.tags = tags
	}

	return t
}

func (t *Transactions) Currency() string {
	balances := []interface{}{
		t.data["final_opening_balance"],
		t.data["opening_balance"],
		t.data["intermediate_opening_balance"],
		t.data["available_balance"],
		t.data["forward_available_balance"],
		t.data["final_closing_balance"],
		t.data["closing_balance"],
		t.data["intermediate_closing_balance"],
		t.data["c_floor_limit"],
		t.data["d_floor_limit"],
	}
	var balance interface{}
	for _, b := range balances {
		if b != nil {
			balance = b
			break
		}
	}

	if balance != nil {
		if a, ok := balance.(Amount); ok {
			return a.currency
		}

		if a, ok := balance.(*Amount); ok {
			return a.currency
		}
	}

	return ""
}

func (t *Transactions) strip(lines []string) []string {
	var result []string

	for _, line := range lines {
		// We don't like carriage returns in case of Windows files so let's
		// just replace them with nothing
		line = strings.ReplaceAll(line, "\r", "")

		// Strip trailing whitespace from lines since they cause incorrect files
		line = strings.TrimSpace(line)

		// Skip separators
		if line == "-" {
			continue
		}

		// Return actual lines
		if line != "" {
			result = append(result, line)
		}
	}

	return result
}

func (t *Transactions) Parse(data string) ([]Transaction, error) {
	// Remove extraneous whitespace and such
	data = strings.Join(t.strip(strings.Split(data, "\n")), "\n")

	var transactions []Transaction
	var currentTransaction *Transaction
	var currentTag *Tag

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a transaction start
		if line == t.tags[Tags.StatementLineStart].Start {
			currentTransaction = &Transaction{data: make(map[int]*Tag)}
			transactions = append(transactions, *currentTransaction)
			continue
		}

		// Check if the line is a transaction end
		if line == t.tags[Tags.StatementLineEnd].Start {
			currentTransaction = nil
			continue
		}

		tagRegex := regexp.MustCompile("^:\n?(?P<full_tag>(?P<tag>[0-9]{2}|NS)(?P<sub_tag>[A-Z])?):")

		// Check if the line contains a tag
		if matches := tagRegex.FindAllStringSubmatch(line, -1); matches != nil {
			validMatches := t.sanitizeTagIDMatches(matches)

			for _, match := range validMatches {
				tagId := t.normalizeTagID(match[1])
				tagValue := match[2]

				// Create new tag and add it to the current transaction
				currentTag = &Tag{ID: tagId, Value: tagValue}
				currentTransaction.Data[tagId] = currentTag

				// Process the tag with pre-processors
				if processors, ok := t.processors[fmt.Sprintf("pre_%d", tagId)]; ok {
					for _, processor := range processors {
						if err := processor(currentTransaction, currentTag); err != nil {
							return nil, fmt.Errorf("error processing tag %d with pre-processor: %w", tagId, err)
						}
					}
				}
			}
		} else if currentTag != nil {
			// If there is no tag, the line is the continuation of the previous tag value
			currentTag.Value += "\n" + line
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning data: %w", err)
	}

	return transactions, nil
}

func (t *Transactions) sanitizeTagIDMatches(matches [][]string) [][]string {
	iNext := 0
	var sanitizedMatches [][]string
	for i, match := range matches {
		// match was rejected
		if i < iNext {
			continue
		}

		// next match would be
		iNext = i + 1

		// normalize tag id
		tagID := t.normalizeTagID(match[1])

		// tag should be known
		if _, ok := t.tags[tagID]; !ok {
			panic(fmt.Sprintf("Unknown tag %v in line: %v", tagID, match))
		}

		// special treatment for long tag content with possible
		// bad line wrap which produces tag_id like line beginnings
		// seen with :86: tag
		if tagID == tags.TRANSACTION_DETAILS.ID {
			// search subsequent tags for unknown tag ids
			// these lines likely belong to the previous tag
			for j := iNext; j < len(matches); j++ {
				nextTagID := t.normalizeTagID(match[1])

				if _, ok := t.tags[nextTagID]; ok {
					// this one is the next valid match
					iNext = j
					break
				}
				// else reject match
			}
		}

		// a valid match
		sanitizedMatches = append(sanitizedMatches, match)
	}

	return sanitizedMatches
}

func (t *Transactions) normalizeTagID(tagID string) int {
	// Since non-digit tags exist, make the conversion optional
	id, err := strconv.Atoi(tagID)
	if err != nil {
		panic(err)
	}

	return id
}
