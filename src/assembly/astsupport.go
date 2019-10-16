package assembly

import "go/ast"

// ExtractFieldName extracts the first name of the field and returns it or ""
func ExtractFieldName(field *ast.Field) string {
	if field != nil && len(field.Names) > 0 {
		return field.Names[0].Name
	}

	return ""
}
