package mjson

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

const json = `{
  "name": {"first": "last", "last": "Anderson"},
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
	spec := &Spec{}
	err = yaml.Unmarshal(buff, spec)
	if err != nil {
		t.Errorf("unmarshal %s to yaml failed: %v", buff, err)
		return
	}
	c := &ParseContext{json: json}
	for _, val := range spec.MappingConfig {
		if val.Group != nil {
			c.group = true
			MappingContext(c, *val.Group)
			for k, v := range *val.Group {
				t.Logf("key: %s, val: %v", k, v)
			}
		}
		if val.Pairs != nil {
			if c.group {
				MappingGroup(c, *val.Pairs)
			} else {
				MappingContext(c, *val.Pairs)
			}
			for k1, v1 := range *val.Pairs {
				t.Logf("k1: %s, v1: %s", k1, v1)
			}
		}
	}
	t.Log(c.json)
}

func TestMappingYAML(t *testing.T) {
	t.Log(MappingYAML(json, "mapping.yaml"))
}
