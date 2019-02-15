package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Coin tests

func TestCoin(t *testing.T) {
	t.Parallel()
	require.Panics(t, func() { NewInt64Coin("b", 0) })

	require.NotPanics(t, func() { NewInt64Coin("btc", 0) })
	require.Panics(t, func() { NewPositiveInt64Coin("btc", 0) })
	require.Panics(t, func() { NewInt64Coin("atom", -1) })
	require.Panics(t, func() { NewPositiveInt64Coin("atom", -1) })

	require.NotPanics(t, func() { NewCoin("atom", NewInt(0)) })
	require.Panics(t, func() { NewCoin("atom", NewInt(-1)) })
	require.Panics(t, func() { NewPositiveCoin("atom", NewInt(0)) })
	require.Panics(t, func() { NewPositiveCoin("atom", NewInt(-1)) })

	// test denom case
	require.Panics(t, func() { NewInt64Coin("Atom", 10) })
	require.Panics(t, func() { NewPositiveInt64Coin("Atom", 10) })
	require.Panics(t, func() { NewCoin("Atom", NewInt(10)) })
	require.Panics(t, func() { NewPositiveCoin("Atom", NewInt(10)) })

	// test leading/trailing spaces

	require.Panics(t, func() { NewInt64Coin("atom ", 10) })
	require.Panics(t, func() { NewPositiveInt64Coin("atom ", 10) })
	require.Panics(t, func() { NewCoin("atom ", NewInt(10)) })
	require.Panics(t, func() { NewPositiveCoin("atom ", NewInt(10)) })

	require.Panics(t, func() { NewInt64Coin("atom ", 10) })
	require.Panics(t, func() { NewPositiveInt64Coin("atom ", 10) })
	require.Panics(t, func() { NewCoin("atom ", NewInt(10)) })
	require.Panics(t, func() { NewPositiveCoin("atom ", NewInt(10)) })

	require.Equal(t, NewInt(5), NewInt64Coin("btc", 5).Amount)
	require.Equal(t, NewInt(5), NewCoin("btc", NewInt(5)).Amount)
}

func TestSameDenomAsCoin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), true},
		{NewInt64Coin("atom", 1), NewInt64Coin("btc", 1), false},
		{NewInt64Coin("steak", 1), NewInt64Coin("steak", 10), true},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.SameDenomAs(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin denominations didn't match, tc #%d", tcIndex)
	}
}

func TestIsEqualCoin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), true},
		{NewInt64Coin("atom", 1), NewInt64Coin("btc", 1), false},
		{NewInt64Coin("steak", 1), NewInt64Coin("steak", 10), false},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.IsEqual(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin equality relation is incorrect, tc #%d", tcIndex)
	}
}

func TestPlusCoin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne    Coin
		inputTwo    Coin
		expected    Coin
		shouldPanic bool
	}{
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), NewInt64Coin("atom", 2), false},
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 0), NewInt64Coin("atom", 1), false},
		{NewInt64Coin("atom", 1), NewInt64Coin("btc", 1), NewInt64Coin("atom", 1), true},
	}

	for tcIndex, tc := range cases {
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Plus(tc.inputTwo) })
		} else {
			res := tc.inputOne.Plus(tc.inputTwo)
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestMinusCoin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne    Coin
		inputTwo    Coin
		expected    Coin
		shouldPanic bool
	}{
		{NewInt64Coin("atom", 1), NewInt64Coin("btc", 1), NewInt64Coin("atom", 1), true},
		{NewInt64Coin("atom", 10), NewInt64Coin("atom", 1), NewInt64Coin("atom", 9), false},
		{NewInt64Coin("atom", 5), NewInt64Coin("atom", 3), NewInt64Coin("atom", 2), false},
		{NewInt64Coin("atom", 5), NewInt64Coin("atom", 0), NewInt64Coin("atom", 5), false},
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 5), Coin{}, true},
	}

	for tcIndex, tc := range cases {
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Minus(tc.inputTwo) })
		} else {
			res := tc.inputOne.Minus(tc.inputTwo)
			require.Equal(t, tc.expected, res, "difference of coins is incorrect, tc #%d", tcIndex)
		}
	}

	tc := struct {
		inputOne Coin
		inputTwo Coin
		expected int64
	}{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), 0}
	res := tc.inputOne.Minus(tc.inputTwo)
	require.Equal(t, tc.expected, res.Amount.Int64())
}

func TestIsGTECoin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), true},
		{NewInt64Coin("atom", 2), NewInt64Coin("atom", 1), true},
		{NewInt64Coin("atom", 1), NewInt64Coin("btc", 1), false},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.IsGTE(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin GTE relation is incorrect, tc #%d", tcIndex)
	}
}

func TestIsLTCoin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
	}{
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), false},
		{NewInt64Coin("atom", 2), NewInt64Coin("atom", 1), false},
		{NewInt64Coin("atom", 0), NewInt64Coin("btc", 1), false},
		{NewInt64Coin("atom", 1), NewInt64Coin("btc", 1), false},
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 1), false},
		{NewInt64Coin("atom", 1), NewInt64Coin("atom", 2), true},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.IsLT(tc.inputTwo)
		require.Equal(t, tc.expected, res, "coin LT relation is incorrect, tc #%d", tcIndex)
	}
}

func TestCoinIsZero(t *testing.T) {
	t.Parallel()
	coin := NewInt64Coin("atom", 0)
	res := coin.IsZero()
	require.True(t, res)

	coin = NewInt64Coin("atom", 1)
	res = coin.IsZero()
	require.False(t, res)
}

// ----------------------------------------------------------------------------
// Coins tests

func TestIsZeroCoins(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne Coins
		expected bool
	}{
		{Coins{}, true},
		{Coins{NewInt64Coin("atom", 0)}, true},
		{Coins{NewInt64Coin("atom", 0), NewInt64Coin("btc", 0)}, true},
		{Coins{NewInt64Coin("atom", 1)}, false},
		{Coins{NewInt64Coin("atom", 0), NewInt64Coin("btc", 1)}, false},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsZero()
		require.Equal(t, tc.expected, res)
	}
}

func TestEqualCoins(t *testing.T) {
	t.Parallel()
	cases := []struct {
		inputOne Coins
		inputTwo Coins
		expected bool
	}{
		{Coins{}, Coins{}, true},
		{Coins{NewInt64Coin("atom", 0)}, Coins{NewInt64Coin("atom", 0)}, true},
		{Coins{NewInt64Coin("atom", 0), NewInt64Coin("btc", 1)}, Coins{NewInt64Coin("atom", 0), NewInt64Coin("btc", 1)}, true},
		{Coins{NewInt64Coin("atom", 0)}, Coins{NewInt64Coin("btc", 0)}, false},
		{Coins{NewInt64Coin("atom", 0)}, Coins{NewInt64Coin("atom", 1)}, false},
		{Coins{NewInt64Coin("atom", 0)}, Coins{NewInt64Coin("atom", 0), NewInt64Coin("btc", 1)}, false},
		{Coins{NewInt64Coin("atom", 0), NewInt64Coin("btc", 1)}, Coins{NewInt64Coin("btc", 1), NewInt64Coin("atom", 0)}, true},
	}

	for tcnum, tc := range cases {
		res := tc.inputOne.IsEqual(tc.inputTwo)
		require.Equal(t, tc.expected, res, "Equality is differed from expected. tc #%d, expected %b, actual %b.", tcnum, tc.expected, res)
	}
}

func TestPlusCoins(t *testing.T) {
	t.Parallel()
	zero := NewInt(0)
	one := NewInt(1)
	two := NewInt(2)

	cases := []struct {
		inputOne Coins
		inputTwo Coins
		expected Coins
	}{
		{Coins{{"atom", one}, {"btc", one}}, Coins{{"atom", one}, {"btc", one}}, Coins{{"atom", two}, {"btc", two}}},
		{Coins{{"atom", zero}, {"btc", one}}, Coins{{"atom", zero}, {"btc", zero}}, Coins{{"btc", one}}},
		{Coins{{"atom", two}}, Coins{{"btc", zero}}, Coins{{"atom", two}}},
		{Coins{{"atom", one}}, Coins{{"atom", one}, {"btc", two}}, Coins{{"atom", two}, {"btc", two}}},
		{Coins{{"atom", zero}, {"btc", zero}}, Coins{{"atom", zero}, {"btc", zero}}, Coins(nil)},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.Plus(tc.inputTwo)
		assert.True(t, res.IsValid())
		require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
	}
}

func TestMinusCoins(t *testing.T) {
	t.Parallel()
	zero := NewInt(0)
	one := NewInt(1)
	two := NewInt(2)

	testCases := []struct {
		inputOne    Coins
		inputTwo    Coins
		expected    Coins
		shouldPanic bool
	}{
		{Coins{{"atom", two}}, Coins{{"atom", one}, {"btc", two}}, Coins{{"atom", one}, {"btc", two}}, true},
		{Coins{{"atom", two}}, Coins{{"btc", zero}}, Coins{{"atom", two}}, false},
		{Coins{{"atom", one}}, Coins{{"btc", zero}}, Coins{{"atom", one}}, false},
		{Coins{{"atom", one}, {"btc", one}}, Coins{{"atom", one}}, Coins{{"btc", one}}, false},
		{Coins{{"atom", one}, {"btc", one}}, Coins{{"atom", two}}, Coins{}, true},
	}

	for i, tc := range testCases {
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Minus(tc.inputTwo) })
		} else {
			res := tc.inputOne.Minus(tc.inputTwo)
			assert.True(t, res.IsValid())
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", i)
		}
	}
}

func TestCoins(t *testing.T) {
	t.Parallel()
	good := Coins{
		{"gas", NewInt(1)},
		{"mineral", NewInt(1)},
		{"tree", NewInt(1)},
	}
	mixedCase1 := Coins{
		{"gAs", NewInt(1)},
		{"MineraL", NewInt(1)},
		{"TREE", NewInt(1)},
	}
	mixedCase2 := Coins{
		{"gAs", NewInt(1)},
		{"mineral", NewInt(1)},
	}
	mixedCase3 := Coins{
		{"gAs", NewInt(1)},
	}
	empty := Coins{
		{"gold", NewInt(0)},
	}
	null := Coins{}
	badSort1 := Coins{
		{"tree", NewInt(1)},
		{"gas", NewInt(1)},
		{"mineral", NewInt(1)},
	}

	// both are after the first one, but the second and third are in the wrong order
	badSort2 := Coins{
		{"gas", NewInt(1)},
		{"tree", NewInt(1)},
		{"mineral", NewInt(1)},
	}
	badAmt := Coins{
		{"gas", NewInt(1)},
		{"tree", NewInt(0)},
		{"mineral", NewInt(1)},
	}
	dup := Coins{
		{"gas", NewInt(1)},
		{"gas", NewInt(1)},
		{"mineral", NewInt(1)},
	}
	neg := Coins{
		{"gas", NewInt(-1)},
		{"mineral", NewInt(1)},
	}

	assert.True(t, good.IsValid(), "Coins are valid")
	assert.False(t, mixedCase1.IsValid(), "Coins denoms contain upper case characters")
	assert.False(t, mixedCase2.IsValid(), "First Coins denoms contain upper case characters")
	assert.False(t, mixedCase3.IsValid(), "Single denom in Coins contains upper case characters")
	assert.True(t, good.IsPositive(), "Expected coins to be positive: %v", good)
	assert.False(t, null.IsPositive(), "Expected coins to not be positive: %v", null)
	assert.True(t, good.IsAllGTE(empty), "Expected %v to be >= %v", good, empty)
	assert.False(t, good.IsAllLT(empty), "Expected %v to be < %v", good, empty)
	assert.True(t, empty.IsAllLT(good), "Expected %v to be < %v", empty, good)
	assert.False(t, badSort1.IsValid(), "Coins are not sorted")
	assert.False(t, badSort2.IsValid(), "Coins are not sorted")
	assert.False(t, badAmt.IsValid(), "Coins cannot include 0 amounts")
	assert.False(t, dup.IsValid(), "Duplicate coin")
	assert.False(t, neg.IsValid(), "Negative first-denom coin")
}

func TestCoinsGT(t *testing.T) {
	t.Parallel()
	one := NewInt(1)
	two := NewInt(2)

	assert.False(t, Coins{}.IsAllGT(Coins{}))
	assert.True(t, Coins{{"atom", one}}.IsAllGT(Coins{}))
	assert.False(t, Coins{{"atom", one}}.IsAllGT(Coins{{"atom", one}}))
	assert.False(t, Coins{{"atom", one}}.IsAllGT(Coins{{"btc", one}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAllGT(Coins{{"btc", one}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllGT(Coins{{"btc", two}}))
}

func TestCoinsGTE(t *testing.T) {
	t.Parallel()
	one := NewInt(1)
	two := NewInt(2)

	assert.True(t, Coins{}.IsAllGTE(Coins{}))
	assert.True(t, Coins{{"atom", one}}.IsAllGTE(Coins{}))
	assert.True(t, Coins{{"atom", one}}.IsAllGTE(Coins{{"atom", one}}))
	assert.False(t, Coins{{"atom", one}}.IsAllGTE(Coins{{"btc", one}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAllGTE(Coins{{"btc", one}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllGTE(Coins{{"btc", two}}))
}

func TestCoinsLT(t *testing.T) {
	t.Parallel()
	one := NewInt(1)
	two := NewInt(2)

	assert.False(t, Coins{}.IsAllLT(Coins{}))
	assert.False(t, Coins{{"atom", one}}.IsAllLT(Coins{}))
	assert.False(t, Coins{{"atom", one}}.IsAllLT(Coins{{"atom", one}}))
	assert.False(t, Coins{{"atom", one}}.IsAllLT(Coins{{"btc", one}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllLT(Coins{{"btc", one}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllLT(Coins{{"btc", two}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllLT(Coins{{"atom", one}, {"btc", one}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAllLT(Coins{{"atom", one}, {"btc", two}}))
	assert.True(t, Coins{}.IsAllLT(Coins{{"atom", one}}))
}

func TestCoinsLTE(t *testing.T) {
	t.Parallel()
	one := NewInt(1)
	two := NewInt(2)

	assert.True(t, Coins{}.IsAllLTE(Coins{}))
	assert.False(t, Coins{{"atom", one}}.IsAllLTE(Coins{}))
	assert.True(t, Coins{{"atom", one}}.IsAllLTE(Coins{{"atom", one}}))
	assert.False(t, Coins{{"atom", one}}.IsAllLTE(Coins{{"btc", one}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllLTE(Coins{{"btc", one}}))
	assert.False(t, Coins{{"atom", one}, {"btc", one}}.IsAllLTE(Coins{{"btc", two}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAllLTE(Coins{{"atom", one}, {"btc", one}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAllLTE(Coins{{"atom", one}, {"btc", two}}))
	assert.True(t, Coins{}.IsAllLTE(Coins{{"atom", one}}))
}

func TestParse(t *testing.T) {
	t.Parallel()
	one := NewInt(1)

	cases := []struct {
		input         string
		valid         bool // if false, we expect an error on parse
		validPositive bool
		expected      Coins // if valid is true, make sure this is returned
	}{
		{"", true, false, nil},
		{"0", false, false, nil}, // // empty denom
		{"1foo", true, true, Coins{{"foo", one}}},
		{"10bar", true, true, Coins{{"bar", NewInt(10)}}},
		{"99bar,1foo", true, true, Coins{{"bar", NewInt(99)}, {"foo", one}}},
		{"98 bar , 1 foo  ", true, true, Coins{{"bar", NewInt(98)}, {"foo", one}}},
		{"  55\t \t bling\n", true, true, Coins{{"bling", NewInt(55)}}},
		{"2foo, 97 bar", true, true, Coins{{"bar", NewInt(97)}, {"foo", NewInt(2)}}},
		{"5 mycoin,", false, false, nil},             // no empty coins in a list
		{"2 3foo, 97 bar", false, false, nil},        // 3foo is invalid coin name
		{"11me coin, 12you coin", false, false, nil}, // no spaces in coin names
		{"1.2btc", false, false, nil},                // amount must be integer
		{"5foo-bar", false, false, nil},              // once more, only letters in coin name
		{"5foo,-3bar", false, false, nil},            // all coins must pass validation
		{"5.2foo", false, false, nil},                // decimal coin
		{"-5foo", false, false, nil},                 // negative coin
		{"0foo", false, false, nil},                  // invalid zero
		{"-0foo", false, false, nil},                 // negative zero
	}

	for tcIndex, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			res, err := ParseCoins(tc.input)
			if !tc.valid {
				require.NotNil(t, err, "%s: %#v. tc #%d", tc.input, res, tcIndex)
			} else if assert.Nil(t, err, "%s: %+v", tc.input, err) {
				require.Equal(t, tc.expected, res, "coin parsing was incorrect, tc #%d", tcIndex)
			}

			res, err = ParsePositiveCoins(tc.input)
			if !tc.validPositive {
				require.NotNil(t, err, "%s: %#v. tc #%d", tc.input, res, tcIndex)
			} else if assert.Nil(t, err, "%s: %+v", tc.input, err) {
				require.Equal(t, tc.expected, res, "coin parsing was incorrect, tc #%d", tcIndex)
			}
		})
	}
}

func TestSortCoins(t *testing.T) {
	t.Parallel()
	good := Coins{
		NewInt64Coin("gas", 1),
		NewInt64Coin("mineral", 1),
		NewInt64Coin("tree", 1),
	}
	empty := Coins{
		NewInt64Coin("gold", 0),
	}
	badSort1 := Coins{
		NewInt64Coin("tree", 1),
		NewInt64Coin("gas", 1),
		NewInt64Coin("mineral", 1),
	}
	badSort2 := Coins{ // both are after the first one, but the second and third are in the wrong order
		NewInt64Coin("gas", 1),
		NewInt64Coin("tree", 1),
		NewInt64Coin("mineral", 1),
	}
	badAmt := Coins{
		NewInt64Coin("gas", 1),
		NewInt64Coin("tree", 0),
		NewInt64Coin("mineral", 1),
	}
	dup := Coins{
		NewInt64Coin("gas", 1),
		NewInt64Coin("gas", 1),
		NewInt64Coin("mineral", 1),
	}

	cases := []struct {
		coins         Coins
		before, after bool // valid before/after sort
	}{
		{good, true, true},
		{empty, false, false},
		{badSort1, false, true},
		{badSort2, false, true},
		{badAmt, false, false},
		{dup, false, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.before, tc.coins.IsValid(), "coin validity is incorrect before sorting, tc #%d", tcIndex)
		tc.coins.Sort()
		require.Equal(t, tc.after, tc.coins.IsValid(), "coin validity is incorrect after sorting, tc #%d", tcIndex)
	}
}

func TestAmountOf(t *testing.T) {
	t.Parallel()
	case0 := Coins{}
	case1 := Coins{
		NewInt64Coin("gold", 0),
	}
	case2 := Coins{
		NewInt64Coin("gas", 1),
		NewInt64Coin("mineral", 1),
		NewInt64Coin("tree", 1),
	}
	case3 := Coins{
		NewInt64Coin("mineral", 1),
		NewInt64Coin("tree", 1),
	}
	case4 := Coins{
		NewInt64Coin("gas", 8),
	}

	cases := []struct {
		coins           Coins
		amountOfGAS     int64
		amountOfMINERAL int64
		amountOfTREE    int64
	}{
		{case0, 0, 0, 0},
		{case1, 0, 0, 0},
		{case2, 1, 1, 1},
		{case3, 0, 1, 1},
		{case4, 8, 0, 0},
	}

	for _, tc := range cases {
		assert.Equal(t, NewInt(tc.amountOfGAS), tc.coins.AmountOf("gas"))
		assert.Equal(t, NewInt(tc.amountOfMINERAL), tc.coins.AmountOf("mineral"))
		assert.Equal(t, NewInt(tc.amountOfTREE), tc.coins.AmountOf("tree"))
	}

	assert.Panics(t, func() { cases[0].coins.AmountOf("Invalid") })
}

func TestCoinsIsAnyGTE(t *testing.T) {
	t.Parallel()
	one := NewInt(1)
	two := NewInt(2)

	assert.False(t, Coins{}.IsAnyGTE(Coins{}))
	assert.False(t, Coins{{"atom", one}}.IsAnyGTE(Coins{}))
	assert.False(t, Coins{}.IsAnyGTE(Coins{{"atom", one}}))
	assert.False(t, Coins{{"atom", one}}.IsAnyGTE(Coins{{"atom", two}}))
	assert.False(t, Coins{{"atom", one}}.IsAnyGTE(Coins{{"btc", one}}))
	assert.True(t, Coins{{"atom", one}, {"btc", two}}.IsAnyGTE(Coins{{"atom", two}, {"btc", one}}))
	assert.True(t, Coins{{"atom", one}}.IsAnyGTE(Coins{{"atom", one}}))
	assert.True(t, Coins{{"atom", two}}.IsAnyGTE(Coins{{"atom", one}}))
	assert.True(t, Coins{{"atom", one}}.IsAnyGTE(Coins{{"atom", one}, {"btc", two}}))
	assert.True(t, Coins{{"btc", two}}.IsAnyGTE(Coins{{"atom", one}, {"btc", two}}))
	assert.False(t, Coins{{"btc", one}}.IsAnyGTE(Coins{{"atom", one}, {"btc", two}}))
	assert.True(t, Coins{{"atom", one}, {"btc", two}}.IsAnyGTE(Coins{{"atom", one}, {"btc", one}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAnyGTE(Coins{{"atom", one}, {"btc", two}}))
	assert.True(t, Coins{{"atom", one}, {"btc", one}}.IsAnyGTE(Coins{{"abc", one}, {"cde", one}, {"btc", one}, {"xyz", one}}))
}

func TestCoinString(t *testing.T) {
	type fields struct {
		Denom  string
		Amount Int
	}
	tests := []struct {
		name string
		coin Coin
		want string
	}{
		{"zero", NewInt64Coin("atom", 0), ""},
		{"value", NewInt64Coin("atom", 10), "10atom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.coin.String())
		})
	}
}

func TestCoinsString(t *testing.T) {
	zero := NewInt64Coin("atom", 0)
	tests := []struct {
		name  string
		coins Coins
		want  string
	}{
		{"zero", Coins{zero}, ""},
		{"value", Coins{NewInt64Coin("atom", 10)}, "10atom"},
		{"zero,positive", Coins{zero, NewInt64Coin("atom", 10)}, "10atom"},
		{"order does not matter", Coins{zero, NewInt64Coin("atom", 10)}, Coins{NewInt64Coin("atom", 10), zero}.String()},
		{"sort", Coins{NewInt64Coin("btc", 5), NewInt64Coin("atom", 10)}, "10atom,5btc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.coins.String())
		})
	}
}
