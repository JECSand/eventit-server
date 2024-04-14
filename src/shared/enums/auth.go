package enums

// SessionType enumerates the potential values for Token.Type
type SessionType int

const (
	USER SessionType = iota + 1
	INTEGRATION
)

// Stringify converts SessionType enum into a string value
func (w SessionType) Stringify() string {
	return [...]string{"USER", "INTEGRATION"}[w-1]
}

// EnumIndex returns the current index of the SessionType enum value
func (w SessionType) EnumIndex() int {
	return int(w)
}

func SessionTypeFromString(inStr string) SessionType {
	switch inStr {
	case "INTEGRATION":
		return INTEGRATION
	default:
		return USER
	}
}

// Role enumerates the potential values for User.Role
type Role int

const (
	MEMBER Role = iota + 1
	ADMIN
	ROOT
)

// Stringify converts Stringify enum into a string value
func (r Role) Stringify() string {
	return [...]string{"MEMBER", "ADMIN", "ROOT"}[r-1]
}

// EnumIndex returns the current index of the Role enum value
func (r Role) EnumIndex() int {
	return int(r)
}
