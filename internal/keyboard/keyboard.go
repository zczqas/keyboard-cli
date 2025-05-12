package keyboard

// KeyDef represents a single key on the keyboard
type KeyDef struct {
	Key    string
	Label  string
	Offset int
}

// GetKeyboardLayout returns the default QWERTY keyboard layout
func GetKeyboardLayout() [][]KeyDef {
	return [][]KeyDef{
		// Top row (no offset)
		{
			{Key: "Q"}, {Key: "W"}, {Key: "E"}, {Key: "R"}, {Key: "T"},
			{Key: "Y"}, {Key: "U"}, {Key: "I"}, {Key: "O"}, {Key: "P"},
			{Key: "[", Label: "["}, {Key: "]", Label: "]"},
		},
		// Home row
		{
			{Key: "A", Offset: 2}, {Key: "S"}, {Key: "D"}, {Key: "F"}, {Key: "G"},
			{Key: "H"}, {Key: "J"}, {Key: "K"}, {Key: "L"},
			{Key: ";", Label: ";"}, {Key: "'", Label: "'"},
		},
		// Bottom row
		{
			{Key: "Z", Offset: 4}, {Key: "X"}, {Key: "C"}, {Key: "V"}, {Key: "B"},
			{Key: "N"}, {Key: "M"},
			{Key: ",", Label: ","}, {Key: ".", Label: "."}, {Key: "/", Label: "/"},
		},
	}
}
