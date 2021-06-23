package types

const (
	QueryGetTreeById = "tree_by_id"
)

type QueryReqGetTreeByID struct {
	Id string `json:"id"`
}

func NewQueryReqGetTreeById(id string) QueryReqGetTreeByID {
	return QueryReqGetTreeByID{Id: id}
}
