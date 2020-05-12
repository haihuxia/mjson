package mjson

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

const json = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age1": 44, "abc": {"nets": ["ig", "fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"first": "Roger", "last": "Craig", "age1": 68, "abc": {"nets": ["fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"first": "Jane", "last": "Murphy", "age1": 47, "abc": {"nets": ["ig", "tw"], "abc": {"bbb": "ccc"}}}
  ]
}`

func TestMappingFromYAML(t *testing.T) {
	buff, err := ioutil.ReadFile("mapping.yaml")
	if err != nil {
		t.Errorf("read file failed: %v", err)
		return
	}
	spec := &spec{}
	err = yaml.Unmarshal(buff, spec)
	if err != nil {
		t.Errorf("unmarshal %s to yaml failed: %v", buff, err)
		return
	}
	c := &parseContext{json: json}
	for _, val := range spec.MappingConfig {
		if val.Group != nil {
			c.group = true
			mappingContext(c, *val.Group)
			for k, v := range *val.Group {
				t.Logf("key: %s, val: %v", k, v)
			}
		}
		if val.Pairs != nil {
			if c.group {
				mappingGroup(c, *val.Pairs)
			} else {
				mappingContext(c, *val.Pairs)
			}
			for k1, v1 := range *val.Pairs {
				t.Logf("k1: %s, v1: %s", k1, v1)
			}
		}
	}
	t.Log(c.json)
}

func TestMapping(t *testing.T) {
	jsonStr := MappingString(json, "name", "my_name")

	jsonStr = MappingString(jsonStr, "name.last", "my_name")

	jsonStr = MappingString(jsonStr, "friends.first", "first_name")

	jsonStr = MappingString(jsonStr, "friends.abc.nets", "abc_nets")

	t.Log(jsonStr)

	t.Log(MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "name", "my_name"))

	t.Log(MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "name.first", "fname"))

	t.Log(MappingString(
		`{"friends": [{"first": "Tom", "last": "Anderson"}}, {"first": "Dale", "last": "Murphy"}}]`,
		"friends.first", "fname"))
}

func TestMappingYAML(t *testing.T) {
	t.Log(MappingYAML(json, "mapping.yaml"))
}
