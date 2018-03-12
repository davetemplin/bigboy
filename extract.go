package main

import (
	"fmt"
)

func extract(target *Target, extractChannel chan<- []map[string]interface{}) {
	if target.prefetch == "" {
		extractToChannel(target, target.params, extractChannel)
	} else {
		prefetch(target, extractChannel)
	}
	close(extractChannel)	
}

func extractToChannel(target *Target, params []interface{}, extractChannel chan<- []map[string]interface{}) {
	rows, err := target.connection.db.Query(target.extract, params...)
	check(err)
	defer rows.Close()

	more := true
	for more {
		var data []map[string]interface{}
		data, more, err = getRows(rows, args.nulls, args.page)
		check(err)
		applyTimezoneAll(target.location, data)

		if len(data) > 0 {
			extractChannel <- data
		}
	}
}

func extractIds(target *Target, ids []uint64) ([][]map[string]interface{}, error) {
	query := fmt.Sprintf(target.extract, csv(ids))
	rows, err := target.connection.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([][]map[string]interface{}, 0)
	more := true
	for more {
		var data []map[string]interface{}
		data, more, err = getRows(rows, args.nulls, args.page)
		if err != nil {
			return nil, err
		}

		applyTimezoneAll(target.location, data)
		if len(data) > 0 {
			result = append(result, data)
		}
	}

	return result, nil
}