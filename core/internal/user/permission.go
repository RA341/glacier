package user

type Permission uint32

const (
	// PermRead Basic permissions (use iota with bit shifting)
	PermRead   Permission = 1 << iota // 1
	PermWrite                         // 2
	PermDelete                        // 4
	PermCreate                        // 8
)

// Has Checked if a permission is set
func (p Permission) Has(flag Permission) bool {
	return p&flag != 0
}

// Add a permission
func (p Permission) Add(flag Permission) Permission {
	return p | flag
}

// Remove a permission
func (p Permission) Remove(flag Permission) Permission {
	return p &^ flag // &^ is Go's AND NOT (bit clear) operator
}

func (p Permission) String() string {
	flags := []struct {
		bit  Permission
		char string
	}{
		{PermCreate, "c"},
		{PermDelete, "d"},
		{PermWrite, "w"},
		{PermRead, "r"},
	}

	result := make([]byte, 0, len(flags))
	for _, f := range flags {
		if p.Has(f.bit) {
			result = append(result, f.char[0])
		} else {
			result = append(result, '-')
		}
	}

	return string(result)
}
