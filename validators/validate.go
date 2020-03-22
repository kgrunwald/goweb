package validators

import (
	"fmt"
	"reflect"
)

func RequiredFields(o interface{}, fields ...string) error {
	v := reflect.ValueOf(o).Elem()
	var missing []string
	for _, field := range fields {
		fieldVal := v.FieldByName(field)
		if !fieldVal.IsValid() || (fieldVal.Interface() == reflect.Zero(fieldVal.Type()).Interface()) {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("Missing required fields %+v", missing)
	}

	return nil
}