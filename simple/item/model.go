package item

type Item struct {
	Id *string `json:"id"`
	Description *string`json:"description"`
	Type *string `json:"type"`
}

type ItemID struct {
	ID *string `json:"id"`
}