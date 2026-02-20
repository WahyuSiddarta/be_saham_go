package helper

import "time"

func ParseRFC3339Pointer(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}

	parsedTime, err := time.Parse(time.RFC3339, *raw)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

func StringPointerOrNil(value string) *string {
	if value == "" {
		return nil
	}
	valueCopy := value
	return &valueCopy
}
