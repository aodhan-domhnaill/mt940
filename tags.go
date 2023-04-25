package mt940

import (
	"errors"
	"regexp"
	"strings"
)

type Tag struct {
	id       string
	re       *regexp.Regexp
	subre    *regexp.Regexp
	name     string
	status   string
	examples []string
}

var (
	ErrTagIdMismatch  = errors.New("mismatched tag id")
	ErrNotImplemented = errors.New("not implemented")
)

var (
	balanceRegexp = regexp.MustCompile(`^(?P<status>[DC])(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?P<currency>.{3})(?P<amount>[0-9,]{0,16})`)
)

var Tags = []Tag{
	Tag{
		name: "TransactionReferenceNumber",
		id:   "20",
		re:   regexp.MustCompile(`(?P<transaction_reference>.{0,16})`),
		examples: []string{
			":20:0000000030210056",
		},
	},
	Tag{
		name: "DateTimeIndication",
		id:   "13",
		re:   regexp.MustCompile(`^(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?P<hour>[0-9]{2})(?P<minute>[0-9]{2})(\+(?P<offset>[0-9]{4})|)$`),
	},
	Tag{
		name: "AccountIdentification",
		id:   "25",
		re:   regexp.MustCompile(`(?P<account_identification>.{0,35})`),
		examples: []string{
			":25:0123456789",
			":25:NL08DEUT0319809633EUR",
			":25:DK0230003617012345",
			":25:FI0281199710012345",
			":25:GB02DABA30128122012345",
			":25:IE02DABA95182390012345",
			":25:NO0281013312345",
			":25:PL02236000050000004550212345",
			":25:SE031200000001220012345",
			":25:FI0734499400012345",
			":25:81199710012345",
		},
	},
	Tag{
		name: "StatementNumber",
		id:   "28C",
		re:   regexp.MustCompile(`(?P<statement_number>[0-9]{1,5})(?:/?(?P<sequence_number>[0-9]{1,5}))?$`),
		examples: []string{
			":28C:3/00001",
			":28C:355/00001",
			":28C:5/1",
			":28C:00532/001",
		},
	},
	Tag{
		name: "BaseBalance",
		re:   balanceRegexp,
	},
	Tag{
		name: "OpeningBalance",
		re:   balanceRegexp,
		id:   "60",
		examples: []string{
			":60F:C111111EUR960",
			":60F:C111118EUR5480,16",
			":60F:C230306DKK985623,04",
		},
	},
	Tag{
		name: "FinalOpeningBalance",
		re:   balanceRegexp,
		id:   "60F",
		examples: []string{
			":60F:C180220GBP16,00",
		},
	},
	Tag{
		name: "StatementLine",
		id:   "61",
		examples: []string{
			":61:1112021202D43,6N477NONREF",
			":61:2303010228CK366336,2NTRFArbi/deposit//1323333800",
		},
	},
	Tag{
		name: "InformationToAccountOwner",
		id:   "86",
		re:   regexp.MustCompile(`(?P<transaction_details>(([\s\S]{0,65}\r?\n?){0,8}[\s\S]{0,65}))`),
		examples: []string{
			":86:/RREF/3825-0031367289 /EREF/1309101116-0000001 /ORDP//NAME/AB AG/REMI/Inv. 1000217666 - 22.724,00, Inv. 1000217693 - 68.130,00,inv. 1000217801 - 16.470,00 /RCMT/EUR 100.000,00 /CHRG/DKK 4,00",
		},
	},
	Tag{
		name: "ClosingBalance",
		re:   balanceRegexp,
		id:   "62",
	},
	Tag{
		name: "IntermediateClosingBalance",
		re:   balanceRegexp,
		id:   "62M",
		examples: []string{
			":62M:C230228DKK12724930,14",
		},
	},
	Tag{
		name: "FinalClosingBalance",
		re:   balanceRegexp,
		id:   "62F",
		examples: []string{
			":62F:C230228DKK12724930,14",
		},
	},

	Tag{
		name: "AvailableBalance",
		re:   balanceRegexp,
		id:   "64",
		examples: []string{
			":64:C230228DKK6698733,27",
			":64:C180220GBP16,00",
		},
	},
	Tag{
		name:     "RelatedReference",
		id:       "21",
		re:       regexp.MustCompile(`(?P<related_reference>.{0,16})`),
		examples: []string{},
	},
	Tag{
		name: "FloorLimitIndicator",
		id:   "34",
		re:   regexp.MustCompile(`(?P<currency>[A-Z]{3})(?P<status>[DC ]?)(?P<amount>[0-9,]{0,16})`),
	},
	Tag{
		name:  "NonSwift",
		id:    "NS",
		re:    regexp.MustCompile(`(?P<non_swift>([0-9]{2}.{0,}\n[0-9]{2}.{0,})*|[^\n]*)$`),
		subre: regexp.MustCompile(`(?P<ns_id>[0-9]{2})(?P<ns_data>.{0,})`),
	},
	Tag{
		name: "IntermediateOpeningBalance",
		re:   balanceRegexp,
		id:   "60M",
	},
	Tag{
		name: "Statement",
		re:   regexp.MustCompile(`^(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?:(?P<entry_month>[0-9]{2})(?P<entry_day>[0-9]{2}))?(?P<status>R?[DC])(?:(?P<funds_code>[A-Z])[\n ]?)?(?P<amount>[[0-9],]{1,15})(?:(?P<id>[A-Z][A-Z0-9 ]{3}))?((?P<customer_reference>(?:(?!//)[^\n]){0,16}))(?://(?P<bank_reference>.{0,23}))?(?:\n?(?P<extra_details>.{0,34}))?$`),
	},
	Tag{
		name: "StatementASNB",
		re:   regexp.MustCompile("(?m)^ (?P<year>[0-9]{2}) # 6!n Value Date (YYMMDD) (?P<month>[0-9]{2}) (?P<day>[0-9]{2}) (?P<entry_month>[0-9]{2})? # [4!n] Entry Date (MMDD) (?P<entry_day>[0-9]{2})? (?P<status>[A-Z]?[DC]) # 2a Debit/Credit Mark (?P<funds_code>[A-Z])? # [1!a] Funds Code (3rd character of the currency code, if needed) \n? # apparently some banks (sparkassen) incorporate newlines here (?P<amount>[[0-9],]{1,15}) # 15d Amount (?P<id>[A-Z][A-Z0-9 ]{3})? # 1!a3!c Transaction Type Identification Code (?P<customer_reference>.{0,34}) # 34x Customer Reference (//(?P<bank_reference>.{0,16}))? # [//16x] Bank Reference (\n?(?P<extra_details>.{0,34}))? # [34x] Supplementary Details $"),
	},

	Tag{
		name: "ForwardAvailableBalance",
		re:   balanceRegexp,
		id:   "65",
	},
	Tag{
		name: "SumEntries",
		id:   "90",
		re:   regexp.MustCompile(`^(?P<number>\d*)(?P<currency>.{3})(?P<amount>[\d,]{1,15})`),
	},
	Tag{
		name:   "SumDebitEntries",
		id:     "90D",
		status: "D",
	},
	Tag{
		name:   "SumCreditEntries",
		id:     "90C",
		status: "C",
	},
}

func (t *Tag) toSlug(name string) string {
	words := regexp.MustCompile(`[A-Z][a-z]+`).FindAllString(name, -1)
	return strings.ToLower(strings.Join(words, "_"))
}

func (t *Tag) Parse(value string) (map[string]string, error) {
	match := t.re.FindStringSubmatch(value)
	if match == nil {
		return nil, errors.New("match failed")
	}

	result := make(map[string]string)
	for i, name := range t.re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result, nil
}

func (t *Tag) Call(transactions *Transactions, value string) string {
	return value
}

func (t *Tag) ID() string {
	return t.id
}

func (t *Tag) Name() string {
	return t.name
}
