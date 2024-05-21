package links

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

const collectionName = "links"

func New(db *mongo.Database, timeout time.Duration) *Repository {
	return &Repository{db: db, timeout: timeout}
}

type Repository struct {
	db      *mongo.Database
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateReq) (database.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	collection := r.db.Collection(collectionName)
	link := database.Link{
		ID:     req.ID,
		URL:    req.URL,
		Title:  req.Title,
		Tags:   req.Tags,
		Images: req.Images,
		UserID: req.UserID,
	}

	_, err := collection.InsertOne(ctx, link)
	if err != nil {
		return database.Link{}, err
	}

	return link, nil
}

func (r *Repository) FindByUserAndURL(ctx context.Context, url, userID string) (database.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	collection := r.db.Collection(collectionName)
	filter := bson.M{"url": url, "userID": userID}

	var result database.Link
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return database.Link{}, nil // Возвращает пустой объект, если запись не найдена
		}
		return database.Link{}, err
	}

	return result, nil
}

func (r *Repository) FindByCriteria(ctx context.Context, criteria Criteria) ([]database.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	collection := r.db.Collection(collectionName)
	filter := bson.M{}
	if criteria.UserID != nil {
		filter["userID"] = *criteria.UserID
	}
	if len(criteria.Tags) > 0 {
		filter["tags"] = bson.M{"$all": criteria.Tags}
	}

	findOptions := options.Find()
	if criteria.Limit != nil {
		findOptions.SetLimit(*criteria.Limit)
	}
	if criteria.Offset != nil {
		findOptions.SetSkip(*criteria.Offset)
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var links []database.Link
	for cursor.Next(ctx) {
		var link database.Link
		if err := cursor.Decode(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
