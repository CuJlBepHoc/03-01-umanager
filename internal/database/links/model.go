package links

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateReq struct {
	ID     primitive.ObjectID
	URL    string
	Title  string
	Tags   []string
	Images []string
}

type Criteria struct {
	Tags   []string
	Limit  *int64
	Offset *int64
}
