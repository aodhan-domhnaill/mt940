package mt940

import (
	"fmt"
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

type TagError struct {
	ParseError
	*Tag
	Value string
}

func (te *TagError) Error() string {
	return fmt.Sprintf(
		"tag parsing error: %v on tag %v value %v",
		te.ParseError.Error(), te.Tag, te.Value)
}

type TagResults map[string]string

var (
	ErrTagIdMismatch  = NewParseError("mismatched tag id")
	ErrNotImplemented = NewParseError("not implemented")
	ErrNotExist       = NewParseError("tag does not exist")
	ErrMisformatedTag = NewParseError("tag is misformatted")
	ErrTagDidNotParse = NewParseError("tag parsing failed")
)

var (
	balanceRegexp = regexp.MustCompile(`(?P<status>[DC])(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?P<currency>.{3})(?P<amount>[0-9,]{0,16})`)
	tagRegex      = regexp.MustCompile(`(?m)^:\n?(?P<full_tag>(?P<tag>[0-9]{2}|NS)(?P<sub_tag>[A-Z])?):`)
)

var Tags = map[string]Tag{
	"20": Tag{
		name: "TransactionReferenceNumber",
		id:   "20",
		re:   regexp.MustCompile(`(?P<transaction_reference>.{0,16})`),
		examples: []string{
			":20:0000000030210056",
		},
	},
	"13": Tag{
		name: "DateTimeIndication",
		id:   "13",
		re:   regexp.MustCompile(`^(?P<year>[0-9]{2})(?P<month>[0-9]{2})(?P<day>[0-9]{2})(?P<hour>[0-9]{2})(?P<minute>[0-9]{2})(\+(?P<offset>[0-9]{4})|)$`),
	},
	"25": Tag{
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
	"28C": Tag{
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
	"60": Tag{
		name: "OpeningBalance",
		id:   "60",
		re:   balanceRegexp,
		examples: []string{
			":60F:C111111EUR960",
			":60F:C111118EUR5480,16",
			":60F:C230306DKK985623,04",
		},
	},
	"60F": Tag{
		name: "FinalOpeningBalance",
		id:   "60F",
		re:   balanceRegexp,
		examples: []string{
			":60F:C180220GBP16,00",
		},
	},
	"61": Tag{
		name: "StatementLine",
		id:   "61",
		re: regexp.MustCompile(
			`(?P<year>[0-9]{2})` + // 6!n Value Date (YYMMDD)
				`(?P<month>[0-9]{2})` +
				`(?P<day>[0-9]{2})` +
				`(?P<entry_month>[0-9]{2})?` + // [4!n] Entry Date (MMDD)
				`(?P<entry_day>[0-9]{2})?` +
				`(?P<status>R?[DC])` + // 2a Debit/Credit Mark
				`(?P<funds_code>[A-Z])?` + // [1!a] Funds Code (3rd character of the currency
				// code, if needed)
				`[\n ]?` + // apparently some banks (sparkassen) incorporate newlines here
				// cuscal can also send a space here as well
				`(?P<amount>[0-9,]{1,15})` + // 15d Amount
				`(?P<id>[A-Z][A-Z0-9 ]{3})?` + // 1!a3!c Transaction Type Identification Code
				// We need the (slow) repeating negative lookahead to search for // so we
				// don't acciddntly include the bank reference in the customer reference.
				`(?P<customer_reference>((?:(?:[^/]|/[^/]))[^\n]){0,16})` + // 16x Customer Reference
				`(//(?P<bank_reference>.{0,23}))?` + // [//23x] Bank Reference
				`(\n?(?P<extra_details>.{0,34}))?`, // [34x] Supplementary Details
		),
		examples: []string{
			":61:1112021202D43,6N477NONREF",
			":61:2303010228CK366336,2NTRFArbi/deposit//1323333800",
			":61:2001010101D65,00NOVBNL47INGB9999999999\n        hr gjlm paulissen    \n  ",
		},
	},
	"86": Tag{
		name: "InformationToAccountOwner",
		id:   "86",
		re:   regexp.MustCompile(`(?P<transaction_details>(([\s\S]{0,65}\r?\n?){0,8}[\s\S]{0,65}))`),
		examples: []string{
			":86:/RREF/3825-0031367289 /EREF/1309101116-0000001 /ORDP//NAME/AB AG/REMI/Inv. 1000217666 - 22.724,00, Inv. 1000217693 - 68.130,00,inv. 1000217801 - 16.470,00 /RCMT/EUR 100.000,00 /CHRG/DKK 4,00",
		},
	},
	"62": Tag{
		name: "ClosingBalance",
		id:   "62",
		re:   balanceRegexp,
	},
	"62M": Tag{
		name: "IntermediateClosingBalance",
		id:   "62M",
		re:   balanceRegexp,
		examples: []string{
			":62M:C230228DKK12724930,14",
		},
	},
	"62F": Tag{
		name: "FinalClosingBalance",
		id:   "62F",
		re:   balanceRegexp,
		examples: []string{
			":62F:C230228DKK12724930,14",
		},
	},
	"64": Tag{
		name: "AvailableBalance",
		id:   "64",
		re:   balanceRegexp,
		examples: []string{
			":64:C230228DKK6698733,27",
			":64:C180220GBP16,00",
		},
	},
	"21": Tag{
		name:     "RelatedReference",
		id:       "21",
		re:       regexp.MustCompile(`(?P<related_reference>.{0,16})`),
		examples: []string{},
	},
	"34": Tag{
		name: "FloorLimitIndicator",
		id:   "34",
		re:   regexp.MustCompile(`(?P<currency>[A-Z]{3})(?P<status>[DC ]?)(?P<amount>[0-9,]{0,16})`),
	},
	"NS": Tag{
		name:  "NonSwift",
		id:    "NS",
		re:    regexp.MustCompile(`(?P<non_swift>([0-9]{2}.{0,}\n[0-9]{2}.{0,})*|[^\n]*)$`),
		subre: regexp.MustCompile(`(?P<ns_id>[0-9]{2})(?P<ns_data>.{0,})`),
	},
	"60M": Tag{
		name: "IntermediateOpeningBalance",
		id:   "60M",
		re:   balanceRegexp,
	},
	"65": Tag{
		name: "ForwardAvailableBalance",
		id:   "65",
		re:   balanceRegexp,
	},
	"90": Tag{
		name: "SumEntries",
		id:   "90",
		re:   regexp.MustCompile(`^(?P<number>[0-9]*)(?P<currency>.{3})(?P<amount>[[0-9],]{1,15})`),
	},
	"90D": Tag{
		name:   "SumDebitEntries",
		id:     "90D",
		status: "D",
	},
	"90C": Tag{
		name:   "SumCreditEntries",
		id:     "90C",
		status: "C",
	},
}

func (t *Tag) Parse(value string) (TagResults, *TagError) {
	ind := tagRegex.FindStringIndex(value)
	if ind == nil {
		return nil, &TagError{ErrMisformatedTag, t, value}
	}

	if t.re == nil {
		return nil, &TagError{ErrNotImplemented, t, value}
	}
	match := t.re.FindStringSubmatch(strings.TrimSpace(value[ind[1]:]))
	if match == nil {
		return nil, &TagError{ErrTagDidNotParse, t, value}
	}

	result := make(map[string]string)
	for i, name := range t.re.SubexpNames() {
		result[name] = match[i]
	}
	return result, nil
}
