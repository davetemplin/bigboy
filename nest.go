package main

import (
	"database/sql"
	"fmt"
)

func queryNest(db *sql.DB, nest *Nest, list []map[string]interface{}) error {
	keys := distinct(take(list, nest.ParentKey))
	query := fmt.Sprintf(nest.fetch, csv(keys))

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	children, err := getAllRows(rows, args.nulls)
	if err != nil {
		return err
	}
	applyTimezoneAll(nest.location, children)

	for _, parent := range list {
		slice := make([]interface{}, 0)
		for _, child := range children {
			if val, ok := child["_parent"]; ok {
				if toUint64(val) == toUint64(parent[nest.ParentKey]) {
					delete(child, "_parent")
					if _, ok := child["_value"]; ok {
						slice = append(slice, child["_value"])
					} else {
						slice = append(slice, child)
					}
				}
			}
		}
		parent[nest.ChildKey] = slice
	}

	return nil
}
