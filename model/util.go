package model

import "strings"

const (
	PublicSchema   = "public"
	DefaultPackage = "model"
)

func Split(input string) (string, string) {
	d := strings.Split(input, ".")
	if len(d) < 2 {
		return PublicSchema, input
	}

	return d[0], d[1]
}

func Join(schema, table string) string {
	return schema + "." + table
}

// Schemas get schemas from table names
func Schemas(tables []string) (schemas []string) {
	index := map[string]struct{}{}
	for _, t := range tables {
		schema, _ := Split(t)
		if _, ok := index[schema]; !ok {
			index[schema] = struct{}{}
			schemas = append(schemas, schema)
		}
	}

	return
}

// DiscloseSchemas discloses "*" in schemas
func DiscloseSchemas(tables []Table, userInput []string) []string {
	index := map[string]struct{}{}

	for _, t := range userInput {
		schema, table := Split(t)

		if table == "*" {
			for _, m := range tables {
				if m.Schema == schema {
					index[Join(m.Schema, m.Name)] = struct{}{}
				}
			}
		} else {
			for _, m := range tables {
				if m.Schema == schema && m.Name == table {
					index[Join(m.Schema, m.Name)] = struct{}{}
				}
			}
		}
	}

	result := make([]string, 0, len(index))
	for key := range index {
		result = append(result, key)
	}

	return result
}

// FollowFKs resolves foreign keys & adds tables in to generate list
func FollowFKs(tables []Table, disclosed []string) []string {
	iTables := map[string]struct{}{}
	for _, d := range disclosed {
		iTables[d] = struct{}{}
	}

	iModels := map[string]int{}
	for i, t := range tables {
		iModels[Join(t.Schema, t.Name)] = i
	}

	count := len(disclosed)
	for i := 0; i < count; i++ {
		table := disclosed[i]
		if m, ok := iModels[table]; ok {
			for _, r := range tables[m].Relations {
				relTable := Join(r.TargetSchema, r.TargetTable)
				_, hasModel := iModels[relTable]
				_, exists := iTables[relTable]

				if hasModel && !exists {
					iTables[relTable] = struct{}{}
					disclosed = append(disclosed, relTable)
					count++
				}
			}
		}
	}

	return disclosed
}

// FilterFKs filters fks that not presented
func FilterFKs(tables []Table, disclosed []string) []Table {
	iTables := map[string]struct{}{}
	for _, d := range disclosed {
		iTables[d] = struct{}{}
	}

	for i := range tables {
		filtered := make([]Relation, 0)

		for _, relation := range tables[i].Relations {
			relTable := Join(relation.TargetSchema, relation.TargetTable)
			if _, ok := iTables[relTable]; ok {
				filtered = append(filtered, relation)
			}
		}

		tables[i].Relations = filtered
	}

	return tables
}

// UniqStrings filter non-unique values
func UniqStrings(input []string) (output []string) {
	index := map[string]struct{}{}

	for _, v := range input {
		if _, ok := index[v]; !ok {
			output = append(output, v)
			index[v] = struct{}{}
		}
	}

	return output
}
