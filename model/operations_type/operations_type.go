package modelOperaTionsType

type OperationType struct {
	OperationTypeID int    `db:"operation_type_id"`
	Description     string `db:"description"`
	Operation       int    `db:"operation"`
}
