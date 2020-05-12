## MJSON

[![Build Status](https://travis-ci.org/haihuxia/mjson.svg?branch=master)](https://travis-ci.org/haihuxia/mjson)

MJSON is a Go package that provides a way to map keys from a json document.

### Installing

```
go get -u github.com/haihuxia/mjson
```

### Example

* MappingString

```
mjson.MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "name", "my_name")
// ==> `{"my_name": {"first": "Tom", "last": "Anderson"}}`

mjson.MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "name.first", "fname")
// ==> `{"name": {"fname": "Tom", "last": "Anderson"}}`

mjson.MappingString(`{"friends": [{"first": "Tom", "last": "Anderson"}, {"first": "Dale", "last": "Murphy"}]}`, "friends.first", "fname")
// ==> `{"friends": [{"fname": "Tom", "last": "Anderson"}, {"fname": "Dale", "last": "Murphy"}]}`

mjson.MappingString(`{"name": {"first": "Tom", "last": "Anderson"}}`, "abc", "my_name")
// ==> `{"name": {"first": "Tom", "last": "Anderson"}}`
// Can't find "abc" path map nothing

mjson.MappingString(`{"friends": [{"first": "Dale", "last": "Murphy", "age": 44, "abc": {"nets": ["ig", "fb", "tw"], "abc": {"bbb": "ccc"}}}]}`, "friends.nets", "abc_nets")
// ==> `{"friends": [{"first": "Dale", "last": "Murphy", "age": 44, "abc": {"nets": ["ig", "fb", "tw"], "abc": {"bbb": "ccc"}}}]}`
// Can't find "friends.nets" path map nothing, "friends.abc.nets" is the correct.
```

* MappingYAML

```yaml
# mapping.yaml
mapping:
  - group:
      name: nick
    pairs:
      first: fname
      last: lname
  - pairs:
      children: childs
  - pairs:
      friends.first: fname
```

```go
package main

import "github.com/haihuxia/mjson"

const json = `{
  "name": {"first": "last", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "abc": {"nets": ["ig", "fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"first": "Roger", "last": "Craig", "age": 68, "abc": {"nets": ["fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"first": "Jane", "last": "Murphy", "age": 47, "abc": {"nets": ["ig", "tw"], "abc": {"bbb": "ccc"}}}
  ]
}`

func main() {
    println(mjson.MappingYAML(json, "mapping.yaml"))
}
```

This will print:

```json
{
  "nick": {"fname": "last", "lname": "Anderson"},
  "age":37,
  "childs": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"fname": "Dale", "last": "Murphy", "age": 44, "abc": {"nets": ["ig", "fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"fname": "Roger", "last": "Craig", "age": 68, "abc": {"nets": ["fb", "tw"], "abc": {"bbb": "ccc"}}},
    {"fname": "Jane", "last": "Murphy", "age": 47, "abc": {"nets": ["ig", "tw"], "abc": {"bbb": "ccc"}}}
  ]
}
```

### License

MJSON source code is available under the MIT [License](https://github.com/haihuxia/mjson/blob/master/LICENSE).