package mt940

import (
	"testing"
)

func TestTag_Parse(t *testing.T) {
	for id, tag := range Tags {
		t.Run(tag.name, func(t *testing.T) {
			if id != tag.id {
				t.Errorf("mismatched id '%v' != '%v'", id, tag.id)
			}
			for _, ex := range tag.examples {
				out, err := tag.Parse(ex)
				if err != nil {
					t.Error(err)
				} else if out == nil {
					t.Error("output was nil")
				} else {
					groupNames := tag.re.SubexpNames()
					for _, g := range groupNames {
						_, ok := out[g]
						if !ok {
							t.Errorf("group '%v' not found", g)
						}
					}
				}
			}
		})
	}
}
