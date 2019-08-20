package enumerate

import (
	"encoding/json"
)

type UserType int

// The UserType values
const (
	_ UserType = iota
	UserTypeAdmin
	UserTypeSupport
)

var userTypeStrings = map[UserType]string{
	UserTypeAdmin:   "admin",
	UserTypeSupport: "support",
}

// String returns a string representation of the UserType.
func (t UserType) String() string {
	if v, ok := userTypeStrings[t]; ok {
		return v
	}
	return ""
}

// UserTypeFromString returns the UserType from the given string.
func UserTypeFromString(s string) UserType {
	for k, v := range userTypeStrings {
		if v == s {
			return k
		}
	}
	return 0
}

// MarshalJSON marshals the UserType to JSON.
func (t UserType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON unmarshals the UserType from JSON.
func (t *UserType) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*t = UserTypeFromString(v)
	return nil
}
