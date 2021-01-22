// Copyright 2021 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package pjson

// Bit flags passed to the "info" parameter of the iter function which
// provides additional information about the
const (
	_       = 1 << iota
	String  // the data is a JSON String
	Number  // the data is a JSON Number
	True    // the data is a JSON True
	False   // the data is a JSON False
	Null    // the data is a JSON NUll
	Object  // the data is a JSON Object (open or close character)
	Array   // the data is a JSON Array (open or close character)
	Comma   // the data is a JSON comma character ','
	Colon   // the data is a JSON colon character ':'
	Start   // the data is the start of the JSON document
	End     // the data is the end of the JSON document
	Open    // the data is an open character (Object or Array, '{' or '[')
	Close   // the data is an close character (Object or Array, '}' or ']')
	Key     // the data is a JSON Object key
	Value   // the data is a JSON Object or Array value
	Escaped // the data is a String with at least one escape character ('\')
	Sign    // the data is a signed Number (has a '-' prefix)
	Dot     // the data is a Number has a dot (radix point)
	E       // the data is a Number in scientific notation (has 'E' or 'e')
)

// Parse JSON.
// The iter function is a callback that fires for every element in the JSON
// document. Elements include all values and tokens.
// The 'start' and 'end' params are the start and end indexes of their
// respective element, such that json[start:end] will equal the complete
// element data.
// The 'info' param provides extra information about the element data.
// Returning 0 from 'iter' will stop the parsing.
// Returning 1 from 'iter' will continue the parsing.
// Returning -1 from 'iter' will skip all children elements in the current
// Object or Array, which only applies when the 'info' for current element
// has the Open bit set, otherwise it effectively works like returning 1.
// This operation returns zero or a negative value if an error occured. This
// value represents the position that the parser was at when it discovered the
// error. To get the true offset multiple this value by -1.
//   e := Parse(json, iter)
//   if e < 0 {
//       pos := e * -1
//       return fmt.Errorf("parsing error at position %d", pos)
//   }
// This operation returns a positive value when successful. If the 'iter'
// stopped early then this value will be the position the parser was at when it
// stopped, otherwise the value will be equal the length of the original json
// document.
func Parse(json []byte, iter func(start, end, info int) int) int {
	i, ok, _ := vdoc(json, 0, iter)
	if !ok {
		i *= -1
	}
	return i
}

var ws = [256]byte{' ': 1, '\t': 1, '\n': 1, '\r': 1}

func isws(ch byte) bool {
	return ws[ch] == 1
}

type vfn func(start, end, info int) int

func vdoc(json []byte, i int, f vfn) (oi int, ok, stop bool) {
	i, ok, stop = vany(json, i, Start, f)
	if stop {
		return i, ok, stop
	}
	for ; i < len(json); i++ {
		if isws(json[i]) {
			continue
		}
		return i, false, true
	}
	return i, true, false
}

var strtoks = [256]byte{
	0x00: 1, 0x01: 1, 0x02: 1, 0x03: 1, 0x04: 1, 0x05: 1, 0x06: 1, 0x07: 1,
	0x08: 1, 0x09: 1, 0x0A: 1, 0x0B: 1, 0x0C: 1, 0x0D: 1, 0x0E: 1, 0x0F: 1,
	0x10: 1, 0x11: 1, 0x12: 1, 0x13: 1, 0x14: 1, 0x15: 1, 0x16: 1, 0x17: 1,
	0x18: 1, 0x19: 1, 0x1A: 1, 0x1B: 1, 0x1C: 1, 0x1D: 1, 0x1E: 1, 0x1F: 1,
	'"': 1, '\\': 1,
}

func isstrtok(ch byte) bool {
	return strtoks[ch] == 1
}

const unroll = true

// validstring - the prefix '"' character has already been processed
func vstring(json []byte, i int) (outi, info int, ok, stop bool) {
	for {
		if unroll {
			for i < len(json)-7 {
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
				if isstrtok(json[i]) {
					if json[i] == '"' {
						return i + 1, info, true, false
					}
					goto tok
				}
				i++
			}
		}
		for ; i < len(json); i++ {
			if isstrtok(json[i]) {
				if json[i] == '"' {
					return i + 1, info, true, false
				}
				goto tok
			}
		}
		break
	tok:
		if json[i] == '"' {
			return i + 1, info, true, false
		}
		if json[i] < ' ' {
			return i, info, false, true
		}
		if json[i] == '\\' {
			info |= Escaped
			i++
			if i == len(json) {
				return i, info, false, true
			}
			switch json[i] {
			default:
				return i, info, false, true
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
			case 'u':
				for j := 0; j < 4; j++ {
					i++
					if i >= len(json) {
						return i, info, false, true
					}
					if !((json[i] >= '0' && json[i] <= '9') ||
						(json[i] >= 'a' && json[i] <= 'f') ||
						(json[i] >= 'A' && json[i] <= 'F')) {
						return i, info, false, true
					}
				}
			}
		}
		i++
	}
	return i, info, false, true
}

func vany(json []byte, i int, dinfo int, f vfn) (oi int, ok, stop bool) {
	for ; i < len(json); i++ {
		if isws(json[i]) {
			continue
		}
		mark := i
		var info int
		if json[i] == '"' {
			i, info, ok, stop = vstring(json, i+1)
			info |= String
		} else if json[i] == '{' {
			f2 := f
			if f != nil {
				r := f(i, i+1, Object|Open|dinfo)
				if r == 0 {
					return i, true, true
				}
				if r == -1 {
					f2 = nil
				}
			}
			i, ok, stop = vobject(json, i+1, f2)
			if stop {
				return i, ok, stop
			}
			if f != nil {
				if dinfo&Start == Start {
					dinfo &= ^Start
					dinfo |= End
				}
				if f(i-1, i, Object|Close|dinfo) == 0 {
					return i, true, true
				}
			}
			return i, true, false
		} else if json[i] == '[' {
			f2 := f
			if f != nil {
				r := f(i, i+1, Array|Open|dinfo)
				if r == 0 {
					return i, true, true
				}
				if r == -1 {
					f2 = nil
				}
			}
			i, ok, stop = varray(json, i+1, f2)
			if stop {
				return i, ok, stop
			}
			if f != nil {
				if dinfo&Start == Start {
					dinfo &= ^Start
					dinfo |= End
				}
				if f(i-1, i, Array|Close|dinfo) == 0 {
					return i, true, true
				}
			}
			return i, true, false
		} else if json[i] == '-' || isnum(json[i]) {
			i, info, ok, stop = vnumber(json, i+1)
			info |= Number
		} else if json[i] == 't' {
			i, ok, stop = vtrue(json, i+1)
			info |= True
		} else if json[i] == 'n' {
			i, ok, stop = vnull(json, i+1)
			info |= Null
		} else if json[i] == 'f' {
			i, ok, stop = vfalse(json, i+1)
			info |= False
		} else {
			return i, false, true
		}
		if stop {
			return i, ok, stop
		}
		if f != nil {
			if dinfo&Start == Start {
				dinfo |= End
			}
			if f(mark, i, info|dinfo) == 0 {
				return i, true, true
			}
		}
		return i, ok, stop
	}
	return i, false, true
}

func vobject(json []byte, i int, f vfn) (oi int, ok, stop bool) {
	for ; i < len(json); i++ {
		if isws(json[i]) {
			continue
		}
		if json[i] == '}' {
			return i + 1, true, false
		}
		if json[i] == '"' {
		key:
			mark := i
			var info int
			i, info, ok, stop = vstring(json, i+1)
			if stop {
				return i, ok, stop
			}
			if f != nil {
				if f(mark, i, info|Key|String) == 0 {
					return i, true, true
				}
			}
			if i, ok, stop = vcolon(json, i); stop {
				return i, ok, stop
			}
			if f != nil {
				if f(i-1, i, Colon) == 0 {
					return i, true, true
				}
			}
			if i, ok, stop = vany(json, i, Value, f); stop {
				return i, ok, stop
			}
			if i, ok, stop = vcomma(json, i, '}'); stop {
				return i, ok, stop
			}
			if json[i] == '}' {
				return i + 1, true, false
			}
			if f != nil {
				if f(i, i+1, Comma) == 0 {
					return i, true, true
				}
			}
			i++
			for ; i < len(json); i++ {
				if isws(json[i]) {
					continue
				}
				if json[i] == '"' {
					goto key
				}
				break
			}
			break
		}
		break
	}
	return i, false, true
}

func varray(json []byte, i int, f vfn) (oi int, ok, stop bool) {
	for ; i < len(json); i++ {
		if isws(json[i]) {
			continue
		}
		if json[i] == ']' {
			return i + 1, true, false
		}
		for ; i < len(json); i++ {
			if isws(json[i]) {
				continue
			}
			if i, ok, stop = vany(json, i, Value, f); stop {
				return i, ok, stop
			}
			if i, ok, stop = vcomma(json, i, ']'); stop {
				return i, ok, stop
			}
			if json[i] == ']' {
				return i + 1, true, false
			}
			if f != nil {
				if f(i, i+1, Comma) == 0 {
					return i, true, true
				}
			}
		}
	}
	return i, false, true
}

func vcolon(json []byte, i int) (outi int, ok, stop bool) {
loop:
	if i < len(json) {
		if json[i] == ':' {
			return i + 1, true, false
		}
		if isws(json[i]) {
			i++
			goto loop
		}
	}
	return i, false, true
}

func vcomma(json []byte, i int, end byte) (outi int, ok, stop bool) {
loop:
	if i < len(json) {
		if json[i] == ',' {
			return i, true, false
		}
		if json[i] == end {
			return i, true, false
		}
		if isws(json[i]) {
			i++
			goto loop
		}
	}
	return i, false, true
}

func isnum(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func vnumber(json []byte, i int) (outi, info int, ok, stop bool) {
	i-- // go back one byte

	if json[i] == '-' {
		info |= Sign
		i++
		if i == len(json) || !isnum(json[i]) {
			return i, info, false, true
		}
	}
	if json[i] == '0' {
		i++
	} else {
		for ; i < len(json); i++ {
			if !isnum(json[i]) {
				goto frac
			}
		}
	}
	if i == len(json) {
		return i, info, true, false
	}

frac:
	if json[i] == '.' {
		info |= Dot
		i++
		if i == len(json) {
			return i, info, false, true
		}
		if !isnum(json[i]) {
			return i, info, false, true
		}
		i++
		for ; i < len(json); i++ {
			if !isnum(json[i]) {
				goto exp
			}
		}
	}
	if i == len(json) {
		return i, info, true, false
	}

exp:
	if json[i] == 'e' || json[i] == 'E' {
		info |= E
		i++
		if i == len(json) {
			return i, info, false, true
		}
		if json[i] == '+' || json[i] == '-' {
			i++
		}
		if i == len(json) {
			return i, info, false, true
		}
		if !isnum(json[i]) {
			return i, info, false, true
		}
		i++
		for ; i < len(json); i++ {
			if !isnum(json[i]) {
				break
			}
		}
	}
	return i, info, true, false
}

func vtrue(json []byte, i int) (outi int, ok, stop bool) {
	if i+3 <= len(json) && (uint32(json[i])|
		uint32(json[i+1])<<8|
		uint32(json[i+2])<<16) == 6649202 {
		return i + 3, true, false
	}
	return i, false, true
}

func vfalse(json []byte, i int) (outi int, ok, stop bool) {
	if i+4 <= len(json) && uint32(json[i])|
		uint32(json[i+1])<<8|
		uint32(json[i+2])<<16|
		uint32(json[i+3])<<24 == 1702063201 {
		return i + 4, true, false
	}
	return i, false, true
}

func vnull(json []byte, i int) (outi int, ok, stop bool) {
	if i+3 <= len(json) && uint32(json[i])|
		uint32(json[i+1])<<8|
		uint32(json[i+2])<<16 == 7105653 {
		return i + 3, true, false
	}
	return i, false, true
}
