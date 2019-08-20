package enumerate

type Encoding int

// The Encoding values
const (
	_ Encoding = iota
	EncodingString
	EncodingInt
)

var encodingStrings = map[Encoding]string{
	EncodingString: "string",
	EncodingInt:    "int",
}

// String returns a string representation of the Encoding.
func (t Encoding) String() string {
	if v, ok := encodingStrings[t]; ok {
		return v
	}
	return ""
}

// EncodingFromString returns the Encoding from the given string.
func EncodingFromString(s string) Encoding {
	for k, v := range encodingStrings {
		if v == s {
			return k
		}
	}
	return 0
}
