
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450062818217984>


package abi

import (
	"reflect"
	"testing"
)

type reflectTest struct {
	name  string
	args  []string
	struc interface{}
	want  map[string]string
	err   string
}

var reflectTests = []reflectTest{
	{
		name: "OneToOneCorrespondance",
		args: []string{"fieldA"},
		struc: struct {
			FieldA int `abi:"fieldA"`
		}{},
		want: map[string]string{
			"fieldA": "FieldA",
		},
	},
	{
		name: "MissingFieldsInStruct",
		args: []string{"fieldA", "fieldB"},
		struc: struct {
			FieldA int `abi:"fieldA"`
		}{},
		want: map[string]string{
			"fieldA": "FieldA",
		},
	},
	{
		name: "MoreFieldsInStructThanArgs",
		args: []string{"fieldA"},
		struc: struct {
			FieldA int `abi:"fieldA"`
			FieldB int
		}{},
		want: map[string]string{
			"fieldA": "FieldA",
		},
	},
	{
		name: "MissingFieldInArgs",
		args: []string{"fieldA"},
		struc: struct {
			FieldA int `abi:"fieldA"`
			FieldB int `abi:"fieldB"`
		}{},
		err: "struct: abi tag 'fieldB' defined but not found in abi",
	},
	{
		name: "NoAbiDescriptor",
		args: []string{"fieldA"},
		struc: struct {
			FieldA int
		}{},
		want: map[string]string{
			"fieldA": "FieldA",
		},
	},
	{
		name: "NoArgs",
		args: []string{},
		struc: struct {
			FieldA int `abi:"fieldA"`
		}{},
		err: "struct: abi tag 'fieldA' defined but not found in abi",
	},
	{
		name: "DifferentName",
		args: []string{"fieldB"},
		struc: struct {
			FieldA int `abi:"fieldB"`
		}{},
		want: map[string]string{
			"fieldB": "FieldA",
		},
	},
	{
		name: "DifferentName",
		args: []string{"fieldB"},
		struc: struct {
			FieldA int `abi:"fieldB"`
		}{},
		want: map[string]string{
			"fieldB": "FieldA",
		},
	},
	{
		name: "MultipleFields",
		args: []string{"fieldA", "fieldB"},
		struc: struct {
			FieldA int `abi:"fieldA"`
			FieldB int `abi:"fieldB"`
		}{},
		want: map[string]string{
			"fieldA": "FieldA",
			"fieldB": "FieldB",
		},
	},
	{
		name: "MultipleFieldsABIMissing",
		args: []string{"fieldA", "fieldB"},
		struc: struct {
			FieldA int `abi:"fieldA"`
			FieldB int
		}{},
		want: map[string]string{
			"fieldA": "FieldA",
			"fieldB": "FieldB",
		},
	},
	{
		name: "NameConflict",
		args: []string{"fieldB"},
		struc: struct {
			FieldA int `abi:"fieldB"`
			FieldB int
		}{},
		err: "abi: multiple variables maps to the same abi field 'fieldB'",
	},
	{
		name: "Underscored",
		args: []string{"_"},
		struc: struct {
			FieldA int
		}{},
		err: "abi: purely underscored output cannot unpack to struct",
	},
	{
		name: "DoubleMapping",
		args: []string{"fieldB", "fieldC", "fieldA"},
		struc: struct {
			FieldA int `abi:"fieldC"`
			FieldB int
		}{},
		err: "abi: multiple outputs mapping to the same struct field 'FieldA'",
	},
	{
		name: "AlreadyMapped",
		args: []string{"fieldB", "fieldB"},
		struc: struct {
			FieldB int `abi:"fieldB"`
		}{},
		err: "struct: abi tag in 'FieldB' already mapped",
	},
}

func TestReflectNameToStruct(t *testing.T) {
	for _, test := range reflectTests {
		t.Run(test.name, func(t *testing.T) {
			m, err := mapArgNamesToStructFields(test.args, reflect.ValueOf(test.struc))
			if len(test.err) > 0 {
				if err == nil || err.Error() != test.err {
					t.Fatalf("Invalid error: expected %v, got %v", test.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				for fname := range test.want {
					if m[fname] != test.want[fname] {
						t.Fatalf("Incorrect value for field %s: expected %v, got %v", fname, test.want[fname], m[fname])
					}
				}
			}
		})
	}
}

