package main

import (
	"database/sql"
)

func getRows(rows *sql.Rows, nulls bool, max int) ([]map[string]interface{}, bool, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, false, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result := make([]map[string]interface{}, 0)
	for max != 0 {
		more := rows.Next()
		if !more {
			return result, false, rows.Err()
		}

		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, false, err
		}

		obj := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if nulls || val != nil {
				var v interface{}
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				obj[col] = v
			}
		}
		result = append(result, obj)

		max--
	}

	return result, true, nil
}

func getAllRows(rows *sql.Rows, null bool) ([]map[string]interface{}, error) {
	data, _, err := getRows(rows, null, -1)
	return data, err
}
