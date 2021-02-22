# pjson

[![GoDoc](https://godoc.org/github.com/tidwall/pjson?status.svg)](https://godoc.org/github.com/tidwall/pjson)

A JSON stream parser for Go and ([Rust](https://github.com/tidwall/pjson.rs))



## Example

The example below prints all string values from a JSON document.

```go
package main

import "github.com/tidwall/pjson"

func main() {
	var json = `
	{
	  "name": {"first": "Tom", "last": "Anderson"},
	  "age":37,
	  "children": ["Sara","Alex","Jack"],
	  "fav.movie": "Deer Hunter",
	  "friends": [
		{"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
		{"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
		{"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
	  ]
	}
	`
	pjson.Parse([]byte(json), 0, func(start, end, info int) int {
		if info&(pjson.String|pjson.Value) == pjson.String|pjson.Value {
			println(json[start:end])
		}
		return 1
	})
}

// output:
// "Tom"
// "Anderson"
// "Sara"
// "Alex"
// "Jack"
// "Deer Hunter"
// "Dale"
// "Murphy"
// "ig"
// "fb"
// "tw"
// "Roger"
// "Craig"
// "fb"
// "tw"
// "Jane"
// "Murphy"
// "ig"
// "tw"
```

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

pjson source code is available under the MIT [License](/LICENSE).

