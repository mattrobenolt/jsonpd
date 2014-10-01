package main

import "testing"

var cases = []struct {
	in  string
	out error
}{
	{"foo", nil},
	{"FOO", nil},
	{"$foo", nil},
	{"$.foo", nil},
	{"foo2", nil},
	{"_foo", nil},
	{"foo[1]", nil},
	{"jQuery181021338515185342904_1405452779651", nil},
	{"", E_EMPTY},
	{"this", E_RESERVED},
	{"2foo", E_INVALID},
	{".foo", E_INVALID},
	{"/*", E_INVALID},
	{"foo-bar", E_INVALID},
	{"foo bar", E_INVALID},
	{"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", E_TOOLONG},
}

func TestValidity(t *testing.T) {
	for _, tt := range cases {
		if ok := isValid(tt.in); ok != tt.out {
			t.Errorf("isValid(%q) => %q, want %q", tt.in, ok, tt.out)
		}
	}
}

func benchmarkValid(cb string, b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValid(cb)
	}
}

func BenchmarkEmpty(b *testing.B) { benchmarkValid("", b) }
func BenchmarkTooLong(b *testing.B) {
	benchmarkValid("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", b)
}
func BenchmarkReserved(b *testing.B)        { benchmarkValid("this", b) }
func BenchmarkBadPatternShort(b *testing.B) { benchmarkValid("foo-bar", b) }
func BenchmarkBadPatternLong(b *testing.B) {
	benchmarkValid("jQuery181021338515185342904_1405452779651-", b)
}
func BenchmarkGoodShort(b *testing.B) { benchmarkValid("foo", b) }
func BenchmarkGoodLong(b *testing.B)  { benchmarkValid("jQuery181021338515185342904_1405452779651", b) }
