package mjson

import (
	"testing"
)

const json = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "abc": {"nets": ["ig", "fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"first": "Roger", "last": "Craig", "age": 68, "abc": {"nets": ["fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"first": "Jane", "last": "Murphy", "age": 47, "abc": {"nets": ["ig", "tw"], "abc": {"bbb": "ccc"}}}
  ]
}`

func TestMapping(t *testing.T) {
	t.Log(MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "name", "my_name"))
	t.Log(MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "name.first", "fname"))
	t.Log(MappingString(
		`{"friends": [{"first": "Tom", "last": "Anderson"}}, {"first": "Dale", "last": "Murphy"}}]`,
		"friends.first", "first_name"))
}

func TestMappingYAML(t *testing.T) {
	t.Log(MappingYAML(json, "mapping.yaml"))
}

func TestMappingStringErrorKey(t *testing.T) {
	// Can't find "abc" path map nothing
	t.Log(MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "abc", "my_name"))

	// Can't find "friends.abc.nets.abc" path map nothing
	t.Log(MappingString(json, "friends.abc.nets.abc", "aaa"))

	// Can't find "friends.nets" path map nothing
	t.Log(MappingString(json, "friends.nets", "abc_nets"))
}

func TestParsePath(t *testing.T) {
	path := "abc\\.bbb"
	t.Logf("%#v", parsePath(path))

	path = "abc\\.bbb.ccc"
	t.Logf("%#v", parsePath(path))

	path = "abc.bbb.ccc"
	t.Logf("%#v", parsePath(path))
}