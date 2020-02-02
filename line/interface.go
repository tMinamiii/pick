package line

import "github.com/google/btree"

type Line interface {
	btree.Item

	ID() uint64
	//Buffer returns the raw buffer
	Buffer() string

	// DisplayString returns the string to be displayed. This means if you have
	// a null spearator, the contents after the separator are not included
	// in this string
	DisplayString() string

	// Output returns the string to be display as peco finishes up doing its
	// thing. This means if you have null separator, the contents before the
	// separtor are not included in this string
	Output() string

	// IsDirty returns true if this line should be forcefully redrawn
	IsDirty() bool

	// SetDirty sets the dirty flag on or off
	SetDirty(bool)
}
