package utils

// func UnmarshalCustomValue(data []byte, typeJsonField, valueJsonField string, customTypes map[string]reflect.Type) (interface{}, error) {
// 	m := map[string]interface{}{}
// 	if err := json.Unmarshal(data, &m); err != nil {
// 		return nil, err
// 	}

// 	typeName := m[typeJsonField].(string)
// 	var value Something
// 	if ty, found := customTypes[typeName]; found {
// 		value = reflect.New(ty).Interface().(Something)
// 	}

// 	valueBytes, err := json.Marshal(m[valueJsonField])
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = json.Unmarshal(valueBytes, &value); err != nil {
// 		return nil, err
// 	}

// 	return value, nil
// }
