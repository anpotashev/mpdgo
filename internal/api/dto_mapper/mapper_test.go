package dto_mapper

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testCase struct {
	testName      string
	from          interface{}
	to            interface{}
	expectedError error
	expected      interface{}
}

func TestMap(t *testing.T) {
	tests := []testCase{
		{
			testName: "from and to contain similar fields",
			from: struct {
				Name string
				Age  int
			}{
				Name: "Alex",
				Age:  18,
			},
			to: struct {
				Name string
				Age  int
			}{},
			expectedError: nil,
			expected: struct {
				Name string
				Age  int
			}{
				Name: "Alex",
				Age:  18,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.testName+" [from is a struct]", func(t *testing.T) {
			valFrom := reflect.ValueOf(test.from)
			valTo := reflect.ValueOf(test.to)
			valExpected := reflect.ValueOf(test.expected)
			if valFrom.Kind() != reflect.Struct || valTo.Kind() != reflect.Struct || valExpected.Kind() != reflect.Struct {
				t.Errorf("incorrect test data. From, to, and expected should be stucts")
			}
			source := test.from
			target := reflect.New(valTo.Type())
			target.Elem().Set(valTo)
			err := Map(source, target.Interface())
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(test.expected, target.Elem().Interface()))
		})
		t.Run(test.testName+" [from is a pointer to struct]", func(t *testing.T) {
			valFrom := reflect.ValueOf(test.from)
			valTo := reflect.ValueOf(test.to)
			valExpected := reflect.ValueOf(test.expected)
			if valFrom.Kind() != reflect.Struct || valTo.Kind() != reflect.Struct || valExpected.Kind() != reflect.Struct {
				t.Errorf("incorrect test data. From, to, and expected should be stucts")
			}
			source := reflect.New(valFrom.Type())
			source.Elem().Set(valFrom)
			target := reflect.New(valTo.Type())
			target.Elem().Set(valTo)
			err := Map(source.Interface(), target.Interface())
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(test.expected, target.Elem().Interface()))
		})
	}
	//t.Run("from and to contain similar fields", func(t *testing.T) {
	//	type source struct {
	//		Name string
	//		Age  int
	//	}
	//	type target struct {
	//		Name string
	//		Age  int
	//	}
	//	s := source{
	//		Name: "from",
	//		Age:  18,
	//	}
	//	var targetVar target
	//	err := Map(s, &targetVar)
	//	assert.NoError(t, err)
	//	assert.Equal(t, s.Name, targetVar.Name)
	//	assert.Equal(t, s.Age, targetVar.Age)
	//})
	//t.Run("from has more fields than to", func(t *testing.T) {
	//	type source struct {
	//		Name string
	//		Age  int
	//		Sex  string
	//	}
	//	type target struct {
	//		Name string
	//		Age  int
	//	}
	//	s := source{
	//		Name: "from",
	//		Age:  18,
	//		Sex:  "MALE",
	//	}
	//	var targetVar target
	//	err := Map(s, &targetVar)
	//	assert.NoError(t, err)
	//	assert.Equal(t, s.Name, targetVar.Name)
	//	assert.Equal(t, s.Age, targetVar.Age)
	//})
	//t.Run("from has more fields than to", func(t *testing.T) {
	//	type source struct {
	//		Name string
	//		Age  int
	//	}
	//	type target struct {
	//		Name string
	//		Age  int
	//		Sex  string
	//	}
	//	s := source{
	//		Name: "from",
	//		Age:  18,
	//	}
	//	var targetVar target
	//	err := Map(s, &targetVar)
	//	assert.NoError(t, err)
	//	assert.Equal(t, s.Name, targetVar.Name)
	//	assert.Equal(t, s.Age, targetVar.Age)
	//	assert.Equal(t, "", targetVar.Sex)
	//})
}
