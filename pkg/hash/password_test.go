package hash

import (
	"testing"
)

const (
	saltOne = "one"
	saltTwo = "two"

	passwordOne = "123456"
	passwordTwo = "654321"
)

type standardTest struct {
	name  string
	salt  string
	input string
}

type multipleHashersTest struct {
	name     string
	test1    standardTest
	test2    standardTest
	errorMsg string
	checker  func(string, string) bool
}

var standardTests = []standardTest{
	{name: "not empty salt - not empty input", salt: saltOne, input: passwordOne},
	{name: "empty salt - not empty input", salt: "", input: passwordOne},
	{name: "not empty salt - empty input", salt: saltOne, input: ""},
	{name: "empty salt - empty input", salt: "", input: ""},
}

var multipleHashersTests = []multipleHashersTest{
	{
		name:     "same salt - same input",
		test1:    standardTest{salt: saltOne, input: passwordOne},
		test2:    standardTest{salt: saltOne, input: passwordOne},
		errorMsg: "hashes must be equal",
		checker:  mustBeEqual,
	},
	{
		name:     "same salt - different input",
		test1:    standardTest{salt: saltOne, input: passwordOne},
		test2:    standardTest{salt: saltOne, input: passwordTwo},
		errorMsg: "hashes must be different",
		checker:  mustBeDifferent,
	},
	{
		name:     "different salt - same input",
		test1:    standardTest{salt: saltOne, input: passwordOne},
		test2:    standardTest{salt: saltTwo, input: passwordOne},
		errorMsg: "hashes must be different",
		checker:  mustBeDifferent,
	},
}

func mustBeEqual(in1, in2 string) bool {
	return in1 == in2
}

func mustBeDifferent(in1, in2 string) bool {
	return in1 != in2
}

func TestStandardHasher(t *testing.T) {
	for _, test := range standardTests {
		t.Run(test.name, func(t *testing.T) {
			hasher := NewSHA1Hasher(test.salt)
			hash, err := hasher.Hash(test.input)

			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if hash == "" {
				t.Error("hash must not be empty")
				return
			}

			if hash == passwordOne {
				t.Error("hash must not be the same as input")
				return
			}
		})
	}
}

func TestMultipleHashers(t *testing.T) {
	for _, test := range multipleHashersTests {
		t.Run(test.name, func(t *testing.T) {
			hasher1 := NewSHA1Hasher(test.test1.salt)
			hasher2 := NewSHA1Hasher(test.test2.salt)

			hash1, err := hasher1.Hash(test.test1.input)
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			hash2, err := hasher2.Hash(test.test2.input)
			if err != nil {
				t.Errorf("unexpecting error: %s", err.Error())
				return
			}

			if !test.checker(hash1, hash2) {
				t.Errorf(test.errorMsg)
			}
		})
	}
}
