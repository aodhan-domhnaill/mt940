package mt940

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	id      int
	re      *regexp.Regexp
	logger  *log.Logger
	pattern string
	name    string
	slug    string
}

func NewTag() *Tag {
	t := &Tag{}
	t.re = regexp.MustCompile(t.pattern)
	t.name = strings.ReplaceAll(fmt.Sprintf("%T", t), "*", "")
	t.slug = t.toSlug(t.name)
	t.logger = log.New(os.Stderr, "tag"+t.slug, 0)
	return t
}

func (t *Tag) toSlug(name string) string {
	words := regexp.MustCompile(`[A-Z][a-z]+`).FindAllString(name, -1)
	return strings.ToLower(strings.Join(words, "_"))
}

func (t *Tag) Parse(transactions *Transactions, value string) map[string]string {
	match := t.re.FindStringSubmatch(value)
	if match != nil {
		t.logger.Debugf("matched (%d) %q against %q, got: %v", len(value), value, t.pattern(), match[1:])
	} else {
		t.logger.Errorf("no match for %q against %q", t.pattern(), value)
		panic(fmt.Sprintf("Unable to parse %v from %q", t, value))
	}
	result := make(map[string]string)
	for i, name := range t.re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

func (t *Tag) Call(transactions *Transactions, value string) string {
	return value
}

func (t *Tag) ID() int {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) Slug() string {
	return t.slug
}

type DateTimeIndication struct {
	Tag
}

func (d *DateTimeIndication) Parse(transactions Transactions, value string) interface{} {
	data := d.Tag.Parse(&transactions, value)
	year, _ := strconv.Atoi("20" + data["year"])
	month, _ := strconv.Atoi(data["month"])
	day, _ := strconv.Atoi(data["day"])
	hour, _ := strconv.Atoi(data["hour"])
	minute, _ := strconv.Atoi(data["minute"])
	offset, _ := strconv.Atoi(data["offset"])
	dateTime := DateTime{
		Year:      year,
		Month:     time.Month(month),
		Day:       day,
		Hour:      hour,
		Minute:    minute,
		OffsetUTC: time.Duration(offset) * time.Minute,
	}
	return map[string]DateTime{"date": dateTime}
}

func (s *StatementNumber) Parse(transactions Transactions, value string) interface{} {
	data := s.Tag.Parse(&transactions, value)
	statementNumber, _ := strconv.Atoi(data["statement_number"])
	sequenceNumber, _ := strconv.Atoi(data["sequence_number"])
	return map[string][]int{"numbers": []int{statementNumber, sequenceNumber}}
}

func (d *DateTimeIndication) Call(transactions *Transactions, value string) map[string]interface{} {
	data := d.Tag.Call(transactions, value)
	return map[string]interface{}{
		"date": &DateTime{
			Date: data["year"].(string) + "-" + data["month"].(string) + "-" + data["day"].(string),
			Time: data["hour"].(string) + ":" + data["minute"].(string) + ":" + "00",
			Zone: "+0000",
		},
	}
}

func NewDateTimeIndication() *DateTimeIndication {
	return &DateTimeIndication{
		Tag{
			id:      13,
			pattern: regexp.MustCompile(`^(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?P<hour>[0-9]{2})(?P<minute>[0-9]{2})(\+(?P<offset>[0-9]{4})|)$`),
		},
	}
}

type TransactionReferenceNumber struct {
	Tag
}

func NewTransactionReferenceNumber() *TransactionReferenceNumber {
	return &TransactionReferenceNumber{
		Tag{
			id:      20,
			pattern: regexp.MustCompile(`(?P<transaction_reference>.{0,16})`),
		},
	}
}

type RelatedReference struct {
	Tag
}

func NewRelatedReference() *RelatedReference {
	return &RelatedReference{
		Tag{
			id:      21,
			pattern: regexp.MustCompile(`(?P<related_reference>.{0,16})`),
		},
	}
}

type AccountIdentification struct {
	Tag
}

func NewAccountIdentification() *AccountIdentification {
	return &AccountIdentification{
		Tag{
			id:      25,
			pattern: regexp.MustCompile(`(?P<account_identification>.{0,35})`),
		},
	}
}

type StatementNumber struct {
	Tag
}

func NewStatementNumber() *StatementNumber {
	return &StatementNumber{
		Tag{
			id:      28,
			pattern: regexp.MustCompile(`(?P<statement_number>[0-9]{1,5})(?:/?(?P<sequence_number>[0-9]{1,5}))?$`),
		},
	}
}

type FloorLimitIndicator struct{}

func (fli FloorLimitIndicator) Id() string {
	return "34"
}

func (fli FloorLimitIndicator) Pattern() string {
	return `(?P<currency>[A-Z]{3})(?P<status>[DC ]?)(?P<amount>[0-9,]{0,16})`
}

func (fli FloorLimitIndicator) Call(transactions Transactions, value map[string]string) map[string]interface{} {
	data := super(FloorLimitIndicator, fli).Call(transactions, value)
	if data["status"] != "" {
		return map[string]interface{}{
			strings.ToLower(data["status"]) + "_floor_limit": NewAmount(data),
		}
	}

	dataD := data
	dataC := data
	dataD["status"] = "D"
	dataC["status"] = "C"
	return map[string]interface{}{
		"d_floor_limit": NewAmount(dataD),
		"c_floor_limit": NewAmount(dataC),
	}
}

type NonSwift struct{}

func (ns *NonSwift) Scope() []Scope {
	return []Scope{Transaction{}, Transactions{}}
}

func (ns *NonSwift) Id() string {
	return "NS"
}

func (ns *NonSwift) Pattern() string {
	return `(?P<non_swift>([0-9]{2}.{0,}\n[0-9]{2}.{0,})*|[^\n]*)$`
}

func (ns *NonSwift) SubPattern() string {
	return `(?P<ns_id>[0-9]{2})(?P<ns_data>.{0,})`
}

func (ns *NonSwift) Call(transactions []Transaction, value map[string]string) map[string]interface{} {
	text := []string{}
	data := value["non_swift"]
	lines := strings.Split(data, "\n")
	subPatternM := regexp.MustCompile(ns.SubPattern())
	for _, line := range lines {
		frag := subPatternM.FindStringSubmatch(line)
		if len(frag) == 3 && frag[2] != "" {
			nsId := frag[1]
			nsData := frag[2]
			value["non_swift_"+nsId] = nsData
			text = append(text, nsData)
		} else if len(text) > 0 && text[len(text)-1] != "" {
			text = append(text, "")
		} else if line = strings.TrimSpace(line); line != "" {
			text = append(text, line)
		}
	}
	value["non_swift_text"] = strings.Join(text, "\n")
	value["non_swift"] = data
	return value
}

type BalanceBase struct {
	Tag
}

func (bb *BalanceBase) Call(transactions, value interface{}) interface{} {
	data := bb.Tag.Call(transactions, value)
	data["amount"] = Amount(data)
	data["date"] = Date(data)
	return map[string]interface{}{
		bb.Slug(): Balance(data),
	}
}

func NewBalanceBase() *BalanceBase {
	return &BalanceBase{
		Tag: Tag{
			Pattern: regexp.MustCompile(`^(?P<status>[DC])(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?P<currency>.{3})(?P<amount>[0-9,]{0,16})`),
		},
	}
}

type OpeningBalance struct {
	BalanceBase
}

func NewOpeningBalance() *OpeningBalance {
	return &OpeningBalance{
		BalanceBase: *NewBalanceBase(),
		ID:          60,
	}
}

type FinalOpeningBalance struct {
	BalanceBase
}

func NewFinalOpeningBalance() *FinalOpeningBalance {
	return &FinalOpeningBalance{
		BalanceBase: *NewBalanceBase(),
		ID:          "60F",
	}
}

type IntermediateOpeningBalance struct {
	BalanceBase
}

func NewIntermediateOpeningBalance() *IntermediateOpeningBalance {
	return &IntermediateOpeningBalance{
		BalanceBase: *NewBalanceBase(),
		ID:          "60M",
	}
}

type Statement struct {
	Tag
}

func NewStatement() *Statement {
	return &Statement{
		Tag{
			pattern: `^(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?:(?P<entry_month>[0-9]{2})(?P<entry_day>[0-9]{2}))?(?P<status>R?[DC])(?:(?P<funds_code>[A-Z])[\n ]?)?(?P<amount>[[0-9],]{1,15})(?:(?P<id>[A-Z][A-Z0-9 ]{3}))?((?P<customer_reference>(?:(?!//)[^\n]){0,16}))(?://(?P<bank_reference>.{0,23}))?(?:\n?(?P<extra_details>.{0,34}))?$`,
		},
	}
}

func (s *Statement) Parse(transactions *Transactions, value string) map[string]interface{} {
	data := s.Tag.Parse(transactions, value)
	if _, ok := data["currency"]; !ok {
		data["currency"] = transactions.Currency
	}

	data["amount"] = NewAmount(data)
	date := data["date"].(Date)

	if data["entry_day"] != nil && data["entry_month"] != nil {
		entry_date := NewDate(data["entry_day"].(string), data["entry_month"].(string), date.Year)
		entry_date := time.Date()

		year := 0
		if date.After(entry_date) && date.Sub(entry_date).Hours()/24 >= 330 {
			year = 1
		} else if entry_date.After(date) && entry_date.Sub(date).Hours()/24 >= 330 {
			year = -1
		}

		data["entry_date"] = entry_date
		data["guessed_entry_date"] = NewDate(entry_date.Day, entry_date.Month, entry_date.Year+year)
	}

	return data
}

type StatementASNB struct {
	Statement
}

func NewStatementASNB() *StatementASNB {
	return &StatementASNB{}
}

func (s *StatementASNB) Pattern() string {
	return "(?m)^ (?P<year>[0-9]{2}) # 6!n Value Date (YYMMDD) (?P<month>[0-9]{2}) (?P<day>[0-9]{2}) (?P<entry_month>[0-9]{2})? # [4!n] Entry Date (MMDD) (?P<entry_day>[0-9]{2})? (?P<status>[A-Z]?[DC]) # 2a Debit/Credit Mark (?P<funds_code>[A-Z])? # [1!a] Funds Code (3rd character of the currency code, if needed) \n? # apparently some banks (sparkassen) incorporate newlines here (?P<amount>[[0-9],]{1,15}) # 15d Amount (?P<id>[A-Z][A-Z0-9 ]{3})? # 1!a3!c Transaction Type Identification Code (?P<customer_reference>.{0,34}) # 34x Customer Reference (//(?P<bank_reference>.{0,16}))? # [//16x] Bank Reference (\n?(?P<extra_details>.{0,34}))? # [34x] Supplementary Details $"
}

func (s *StatementASNB) Parse(line string) (Transaction, error) {
	re := regexp.MustCompile(s.Pattern())
	match := re.FindStringSubmatch(line)
	if len(match) == 0 {
		return Transaction{}, ErrNoMatch
	}

	transaction := Transaction{
		Date:              match[1] + "-" + match[2] + "-" + match[3],
		Amount:            match[9],
		CustomerReference: match[10],
		BankReference:     match[12],
		ExtraDetails:      match[14],
	}

	return transaction, nil

}

func (s *StatementASNB) Call(transactions []Transaction, value string) []Transaction {
	return s.Statement.Call(transactions, value)
}

type ClosingBalance struct {
	BalanceBase
}

func (cb ClosingBalance) ID() string {
	return "62"
}

type IntermediateClosingBalance struct {
	ClosingBalance
}

func (icb IntermediateClosingBalance) ID() string {
	return "62M"
}

type FinalClosingBalance struct {
	ClosingBalance
}

func (fcb FinalClosingBalance) ID() string {
	return "62F"
}

type AvailableBalance struct {
	BalanceBase
}

func (ab AvailableBalance) ID() string {
	return "64"
}

type ForwardAvailableBalance struct {
	BalanceBase
}

func (fab ForwardAvailableBalance) ID() string {
	return "65"
}

type TransactionDetails struct {
	Tag
}

func (td TransactionDetails) ID() string {
	return "86"
}

func (td TransactionDetails) Scope() Transaction {
	return Transaction{}
}

func (td TransactionDetails) Pattern() string {
	return `(?P<transaction_details>(([\s\S]{0,65}\r?\n?){0,8}[\s\S]{0,65}))`
}

type SumEntries struct {
	Tag
	status string
}

func (se SumEntries) ID() string {
	return "90"
}

func (se SumEntries) Pattern() string {
	return `^(?P<number>\d*)(?P<currency>.{3})(?P<amount>[\d,]{1,15})`
}

func (se SumEntries) Call(transactions []Transaction, value string) map[string]interface{} {
	data := se.Tag.Call(transactions, value)
	data["status"] = se.status

	return map[string]interface{}{se.Slug(): SumAmount{
		Number:   data["number"].(int),
		Currency: data["currency"].(string),
		Amount:   data["amount"].(string),
		Status:   data["status"].(string),
	}}
}

type SumDebitEntries struct {
	SumEntries
}

func (sde SumDebitEntries) ID() string {
	return "90D"
}

func (sde SumDebitEntries) Status() string {
	return "D"
}

type SumCreditEntries struct {
	SumEntries
}

func (sce SumCreditEntries) ID() string {
	return "90C"
}

func (sce SumCreditEntries) Status() string {
	return "C"
}
