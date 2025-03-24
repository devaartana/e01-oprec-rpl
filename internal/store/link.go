package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Link struct {
	Slug         string    `bson:"slug" json:"slug"`
	OriginalUrl  string    `bson:"original_url" json:"original_url"`
	Created_at   time.Time `bson:"created_at" json:"created_at"`
	Expired_date time.Time `bson:"expired_date" json:"expired_date"`
}

type UserLinks struct {
	Email string `bson:"email" json:"email"`
	Links []Link `bson:"links" json:"links"`
}

type LinkStore struct {
	db *mongo.Client
}

func (l *LinkStore) Create(ctx context.Context, email string, link *Link) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if condition := l.SlugExist(ctx, link.Slug); condition {
		return ErrDuplicateSlug
	}

	filter := bson.M{"email": email}
	update := bson.M{
        "$push": bson.M{"links": link}, 
    }

    _, err := l.db.Database(DB).Collection(Collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkStore) GetBySlug(ctx context.Context, slug string) (*Link, error) {
    ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
    defer cancel()

    filter := bson.M{
        "links": bson.M{
            "$elemMatch": bson.M{
                "slug": slug, 
			},
		},
	}
    
    projection := bson.M{
        "links.$": 1,
        "_id": 0,
    }
    
    options := options.FindOne().SetProjection(projection)
    
    var result struct {
        Links []Link `bson:"links"` 
    }

    if err := l.db.Database(DB).Collection(Collection).FindOne(ctx, filter, options).Decode(&result); err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, ErrNotFound
        }
        return nil, err
    }

    return &result.Links[0], nil
}

func (l *LinkStore) SlugExist(ctx context.Context, slug string) bool {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{
		"links": bson.M{
			"$elemMatch": bson.M{
				"slug": slug, 
			},
		},
	}
	var link Link
	if err := l.db.Database(DB).Collection(Collection).FindOne(ctx, filter).Decode(&link); err != nil {
		return false
	}

	return true
}

func (l *LinkStore) GetAll(ctx context.Context, email string) ([]Link, error) {
    ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
    defer cancel()

    filter := bson.M{"email": email}
    projection := bson.M{"links": 1, "_id": 0}
    options := options.FindOne().SetProjection(projection)

    var result struct {
        Links []Link `bson:"links"`
    }

    err := l.db.Database(DB).Collection(Collection).FindOne(ctx, filter, options).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, ErrNotFound
        }
        return nil, err
    }

    return result.Links, nil
}

func (l *LinkStore) DeleteBySlug(ctx context.Context, email string, slug string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{
        "email": email,
    }

    update := bson.M{
        "$pull": bson.M{
            "links": bson.M{
                "slug": slug,
            },
        },
    }

    result, err := l.db.Database(DB).Collection(Collection).UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
	
	if result.ModifiedCount == 0 {
        return ErrNotFound
    }

	return nil
}	

func (l *LinkStore) UpdateBySlug(ctx context.Context, email string, link *Link) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{
        "email": email,
        "links.slug": link.Slug,
    }

	update := bson.M{
        "$set": bson.M{
            "links.$.original_url": link.OriginalUrl, 
            "links.$.expired_date": link.Expired_date,
        },
    }

	result, err := l.db.Database(DB).Collection(Collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}
