package boolexpr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBooleans(t *testing.T) {
	boolExpr := And([]map[string]any{
		{"a": Eq("val_a")},
		Or([]map[string]any{
			{"b": Eq("val_b")},
			Not(map[string]any{"c": Eq("not_val_c")}),
			{"d": NotEq("not_val_d")},
		}),
		{"e": In([]string{"opt_e_1", "opt_e_2"})},
		{"f": NotIn([]string{"not_opt_f_1", "not_opt_f_2"})},
		{"gnull": IsNull(true)},
		{"hnotnull": IsNull(false)},
	})
	boolExprExpanded := map[string]any{
		"_and": []map[string]any{
			{
				"a": map[string]any{"_eq": "val_a"},
			},
			{
				"_or": []map[string]any{
					{
						"b": map[string]any{"_eq": "val_b"},
					},
					{
						"_not": map[string]any{
							"c": map[string]any{
								"_eq": "not_val_c",
							},
						},
					},
					{
						"d": map[string]any{"_neq": "not_val_d"},
					},
				},
			},
			{
				"e": map[string]any{
					"_in": []string{"opt_e_1", "opt_e_2"},
				},
			},
			{
				"f": map[string]any{
					"_nin": []string{"not_opt_f_1", "not_opt_f_2"},
				},
			},
			{
				"gnull": map[string]any{"_is_null": true},
			},
			{
				"hnotnull": map[string]any{"_is_null": false},
			},
		},
	}
	assert.Equal(t, boolExprExpanded, boolExpr, "Failed")
}
