// MJSON provides a way to map keys from a json document.
package mjson

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type (
	Spec struct {
		MappingConfig []*mappingConfig `yaml:"mapping"`
	}

	mappingConfig struct {
		Group *map[string]string `yaml:"group"`
		Pairs *map[string]string `yaml:"pairs"`
	}

	pathResult struct {
		part string
		path string
		more bool
	}

	parseContext struct {
		json       string
		val        string
		groupIndex int
		group      bool
	}
)

// Find the mapped key by path. A path is in dot syntax, such as "name.last" or "age".
// When the key is found, it will be replaced by val.
//
//	{
//  	"name": {"first": "Tom", "last": "Anderson"},
//  	"age":37,
//  	"children": ["Sara","Alex","Jack"],
//  	"friends": [
//    		{"first": "James", "last": "Murphy"},
//    		{"first": "Roger", "last": "Craig"}
//  	]
//	}
//
// path: "name.first",		val: "fname"  >> 	found key "first" 		>> "name.fname"
// path: "children", 		val: "childs" >> 	found key "children" 	>> "childs"
// path: "friends.first", 	val: "fname"  >> 	found key "first" twice >> "friends.fname"
func MappingString(json, path, val string) string {
	c := &parseContext{json: json}
	var i int
	mappingPairs(c, path, val, i)
	return c.json
}

func MappingSpec(json string, spec *Spec) string {
	c := &parseContext{json: json}
	for _, val := range spec.MappingConfig {
		if val.Group != nil {
			c.group = true
			mappingContext(c, *val.Group)
		}
		if val.Pairs != nil {
			if c.group {
				mappingGroup(c, *val.Pairs)
			} else {
				mappingContext(c, *val.Pairs)
			}
		}
	}
	return c.json
}

func MappingYAML(json string, filePath string) string {
	buff, err := ioutil.ReadFile(filePath)
	if err != nil {
		return json
	}
	spec := &Spec{}
	err = yaml.Unmarshal(buff, spec)
	if err != nil {
		return json
	}
	return MappingSpec(json, spec)
}

func mappingContext(c *parseContext, m map[string]string) string {
	for path, val := range m {
		var i int
		mappingPairs(c, path, val, i)
	}
	return c.json
}

func mappingGroup(c *parseContext, m map[string]string) string {
	for path, val := range m {
		i := c.groupIndex
		mappingPairs(c, path, val, i)
	}
	c.group = false
	c.groupIndex = 0
	return c.json
}

func mappingPairs(c *parseContext, path, val string, i int) {
	c.val = val
	for ; i < len(c.json); i++ {
		if c.json[i] == '{' {
			i++
			mappingObject(c, path, i)
			break
		}
		if c.json[i] == '[' {
			i++
			mappingArray(c, path, i)
			break
		}
	}
}

func mappingObject(c *parseContext, path string, i int) {
	rp := parsePath(path)
	parseObject(c, rp, i)
}

func mappingArray(c *parseContext, path string, i int) {
	for ; i < len(c.json); i++ {
		if c.json[i] == '{' {
			i++
			mappingObject(c, path, i)
			return
		}
	}
}

func parsePath(path string) (r pathResult) {
	for i := 0; i < len(path); i++ {
		if path[i] == '.' {
			r.part = path[:i]
			r.path = path[i+1:]
			r.more = true
			return
		}
	}
	if len(r.path) == 0 {
		r.part = path
	}
	return
}

func parseObject(c *parseContext, rp pathResult, i int) {
	var key string
	for ; i < len(c.json); i++ {
		if c.json[i] == '"' {
			i++
			key, i = parseJSONKey(c.json, i)
			if rp.part == key {
				if !rp.more {
					if c.group && c.groupIndex == 0 {
						c.groupIndex = i
					}
					if c.val != "_" {
						s := i - len(key) - 1
						c.json = c.json[:s] + c.val + c.json[i-1:]
						i = s + len(c.val) + 1
					}
					break
				}
				for ; i < len(c.json); i++ {
					if c.json[i] == '{' {
						if !rp.more {
							return
						}
						rp = parsePath(rp.path)
						i++
						parseObject(c, rp, i)
						return
					} else if c.json[i] == '"' {
						return
					} else if c.json[i] == '}' {
						return
					} else if c.json[i] == '[' {
						rp = parsePath(rp.path)
						parseArray(c, rp, i)
						return
					}
				}
			} else {
				i = jumpJSONValue(c.json, i)
			}
		} else if c.json[i] == '}' {
			return
		} else if c.json[i] == '{' {
			i = jumpObject(c.json, i)
		}
	}
}

func parseArray(c *parseContext, rp pathResult, i int) {
	for ; i < len(c.json); i++ {
		if c.json[i] == '{' {
			i++
			parseObject(c, rp, i)
			i = jumpObject(c.json, i)
		} else if c.json[i] == ']' {
			break
		}
	}
}

func parseJSONKey(json string, i int) (string, int) {
	var s = i
	var key string
	for ; i < len(json); i++ {
		if json[i] > '\\' {
			continue
		} else if json[i] == '"' {
			i, key = i+1, json[s:i]
			break
		}
	}
	return key, i
}

func jumpJSONValue(json string, i int) int {
	for ; i < len(json); i++ {
		if json[i] == '"' {
			i++
			i = jumpString(json, i)
			break
		} else if json[i] == '{' {
			i++
			i = jumpObject(json, i)
			break
		} else if json[i] == '[' {
			i++
			i = jumpArray(json, i)
			break
		} else if json[i] == ',' {
			break
		}
	}
	return i
}

func jumpString(json string, i int) int {
	_, i = parseJSONKey(json, i)
	return i
}

func jumpObject(json string, i int) int {
	depth := 1
	for ; i < len(json); i++ {
		if json[i] == '{' {
			depth++
		} else if json[i] == '}' {
			depth--
			if depth <= 0 {
				break
			}
		}
	}
	return i
}

func jumpArray(json string, i int) int {
	for ; i < len(json); i++ {
		if json[i] == ']' {
			break
		}
	}
	return i
}