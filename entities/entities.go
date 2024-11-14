package entities

func GetAllEntities() []interface{} {
	return []interface{}{
		Article{},
		Author{},
	}
}
