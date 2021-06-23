package types

import (
	"fmt"
)

const (
	ModuleName    = "trees"
	RouterKey     = ModuleName
	QuerierRoute  = ModuleName
	StoreKey      = ModuleName
	TreeKey       = "tree-value-"
	TreeCountKey  = "tree-count-"
	OrderKey      = "order-value-"
	OrderCountKey = "order-count-"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetLimitedOrderKey(treeId string, orderId string) []byte {
	return []byte(OrderKey + fmt.Sprintf("%s-%s", treeId, orderId))
}

func GetLimitOrderCountKey(treeId string) []byte {
	return []byte(OrderCountKey + fmt.Sprintf("%s", treeId))
}
