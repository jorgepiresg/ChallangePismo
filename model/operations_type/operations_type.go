package modelOperaTionsType

type OperationType struct {
	OperationTypeID int    `db:"operation_type_id" json:"operation_type_id"`
	Description     string `db:"description" json:"description"`
	Operation       int    `db:"operation" json:"operation"`
}
