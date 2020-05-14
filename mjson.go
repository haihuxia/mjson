// MJSON provides a way to map keys from a json document.
package mjson

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"unsafe"
)

type (
	spec struct {
		Mappings []*MappingConfig `yaml:"mapping"`
	}

	MappingConfig struct {
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
// If the key contains a dot, it can be escaped by "\", such as "fav\.movie".
// When the key is found, it will be replaced by val.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "fav.movie": "Deer Hunter"
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//
//	path: "name.first",    val: "fname"       >> found key "first"       >> replace with "name.fname"
//	path: "children",      val: "my_children" >> found key "children"    >> replace with "my_children"
//	path: "fav.movie"      val: "fav_movie"   >> found key "fav.movie"   >> replace with "fav_movie"
//	path: "friends.first", val: "fname"       >> found key "first" twice >> replace all with "friends.fname"
func MappingString(json, path, val string) string {
	c := &parseContext{json: json}
	var i int
	mappingPairs(c, path, val, i)
	return c.json
}

func Mapping(json string, m []*MappingConfig) string {
	c := &parseContext{json: json}
	for _, val := range m {
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

// The YAML file can be used for group configuration, and the group mapping speed is faster.
//
// mapping.yaml
//
//	mapping:
//	  - group:
//	      name: nick
//	    pairs:
//	      first: fname
//	      last: lname
//	  - pairs:
//	      age: my_age
//	      children: my_children
//	      fav\.movie: fav_movie
//	      friends.first: fname
//
// The grouping is based on the key of the object or array.
// If the group name does not need to be mapped, you can use the underscore as the value.
//
//	mapping:
//	  - group:
//	      name: _
//	    pairs:
//	      first: fname
func MappingYAML(json string, filePath string) string {
	buff, err := ioutil.ReadFile(filePath)
	if err != nil {
		return json
	}
	s := &spec{}
	err = yaml.Unmarshal(buff, s)
	if err != nil {
		return json
	}
	return Mapping(json, s.Mappings)
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data,
		Len:  len(s),
		Cap:  len(s),
	}))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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
			rp := parsePath(path)
			parseArray(c, rp, i)
			break
		}
	}
}

func mappingObject(c *parseContext, path string, i int) {
	rp := parsePath(path)
	parseObject(c, rp, i)
}

func parsePath(path string) (r pathResult) {
	var s int
	for i := 0; i < len(path); i++ {
		if path[i] > '\\' {
			continue
		}
		if path[i] == '\\' {
			s = i
			i++
			continue
		} else if path[i] == '.' {
			if s > 0 {
				r.part = path[:s] + path[s+1:i]
			} else {
				r.part = path[:i]
			}
			r.path = path[i+1:]
			r.more = true
			return
		}
	}
	if s > 0 {
		r.part = path[:s] + path[s+1:]
		return
	}
	r.part = path
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
					if c.val != "_" {
						s := i - len(key) - 1
						c.json = c.json[:s] + c.val + c.json[i-1:]
						i = s + len(c.val) + 1
					}
					if c.group && c.groupIndex == 0 {
						c.groupIndex = i
					}
					break
				}
				for ; i < len(c.json); i++ {
					switch c.json[i] {
					case '"':
						return
					case '{':
						rp = parsePath(rp.path)
						i++
						parseObject(c, rp, i)
						return
					case '[':
						rp = parsePath(rp.path)
						parseArray(c, rp, i)
						return
					case ']', '}':
						return
					}
				}
			} else {
				i = jumpJSONValue(c.json, i)
			}
		} else if c.json[i] == '{' {
			i = jumpObject(c.json, i)
		} else if c.json[i] == '}' {
			return
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
		if json[i] == '"' {
			i, key = i+1, json[s:i]
			return key, i
		}
	}
	return key, i
}

func jumpJSONValue(json string, i int) int {
	for ; i < len(json); i++ {
		switch json[i] {
		default:
			continue
		case '"':
			i++
			return jumpString(json, i)
		case '{':
			i++
			return jumpObject(json, i)
		case '[':
			i++
			return jumpArray(json, i)
		case ',':
			return i
		}
	}
	return i
}

func jumpString(json string, i int) int {
	for ; i < len(json); i++ {
		if json[i] == '"' {
			return i + 1
		}
	}
	return i
}

func jumpObject(json string, i int) int {
	depth := 1
	for ; i < len(json); i++ {
		if json[i] < '{' {
			continue
		}
		if json[i] == '{' {
			depth++
		} else if json[i] == '}' {
			depth--
			if depth <= 0 {
				return i
			}
		}
	}
	return i
}

func jumpArray(json string, i int) int {
	for ; i < len(json); i++ {
		if json[i] == ']' {
			return i
		}
	}
	return i
}