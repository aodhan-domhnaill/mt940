package mt940

import (
	"testing"
)

func TestTag_Parse(t *testing.T) {
	for _, tag := range Tags {
		t.Run(tag.name, func(t *testing.T) {
			for _, ex := range tag.examples {
				out, err := tag.Parse(ex)
				if err != nil {
					t.Error(err)
				}
				if out == nil {
					t.Error("output was nil")
				}
			}
		})
	}
}
