package pola_test

import (
	"testing"

	"github.com/ipsusila/pola"
	"github.com/stretchr/testify/assert"
)

func TestConvertion(t *testing.T) {
	type item struct {
		val   any
		Bool  bool
		Int   int64
		Float float64
		Time  bool
	}
	type aliasString string
	testItems := []item{
		{val: "2022-10-11", Bool: false, Int: 0, Float: 0, Time: true},
		{val: "2022/12/22", Bool: false, Int: 0, Float: 0, Time: true},
		{val: "10", Bool: false, Int: 10, Float: 10, Time: false},
		{val: nil, Bool: false, Int: 0, Float: 0, Time: false},
		{val: 11, Bool: true, Int: 11, Float: 11, Time: false},
		{val: -11, Bool: false, Int: -11, Float: -11, Time: false},
		{val: 123.1111, Bool: true, Int: 123, Float: 123.1111, Time: false},
		{val: true, Bool: true, Int: 1, Float: 1, Time: false},
		{val: false, Bool: false, Int: 0, Float: 0, Time: false},
		{val: "", Bool: false, Int: 0, Float: 0, Time: false},
		{val: 0, Bool: false, Int: 0, Float: 0, Time: false},
		{val: aliasString("2022-10-11 13:00:00"), Bool: false, Int: 0, Float: 0, Time: true},
		{val: aliasString("2022/10/11 13:00:00"), Bool: false, Int: 0, Float: 0, Time: true},
		{val: aliasString("123"), Bool: false, Int: 123, Float: 123, Time: false},
	}

	for i, iv := range testItems {
		vb := pola.ToBool(iv.val)
		assert.Equal(t, vb, iv.Bool, "ToBool>%d: %#v", i, iv.val)

		vi, _ := pola.ToInt(iv.val)
		assert.Equal(t, iv.Int, vi, "ToInt>%d: %#v", i, iv.val)

		vf, _ := pola.ToFloat(iv.val)
		assert.Equal(t, iv.Float, vf, "ToFloat>%d: %#v", i, iv.val)

		_, ok := pola.ToTime(iv.val)
		assert.Equal(t, iv.Time, ok, "ToTime>%d: %#v", i, iv.val)
	}
}
