package exodus

import "github.com/gostalt/exodus/column"

// Schema is a slice of items that satisfy the Columnable interface.
type Schema []column.Columnable
