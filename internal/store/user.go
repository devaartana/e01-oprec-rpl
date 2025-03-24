package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username   string    `bson:"username" json:"username"`
	Password   []byte    `bson:"password" json:"password"`
	Email      string    `bson:"email" json:"email"`
	Created_at time.Time `bson:"created_at" json:"created_at"`
	Links      []Link    `bson:"links" json:"links"`
}

type UserStore struct {
	db *mongo.Client
}

func (u *User) SetPassword(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

func (u *User) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(text))
}

func (s *UserStore) Create(ctx context.Context, user *User) error {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	if _, err := s.GetByEmail(ctx, user.Email); err == nil {
		return ErrDuplicateEmail
	}

	_, err := s.db.Database(DB).Collection(Collection).InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{"email": user.Email}
	updateData := bson.M{
		"$set": bson.M{
			"username": user.Username,
		},
	}

	result, err := s.db.Database(DB).Collection(Collection).UpdateOne(ctx, filter, updateData)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{"email": email}
	projection := bson.M{
		"username":   1,
		"email":      1,
		"password":   1,
		"created_at": 1,
	}
	options := options.FindOne().SetProjection(projection)

	var result User
	err := s.db.Database(DB).Collection(Collection).FindOne(ctx, filter, options).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *UserStore) DeleteByEmail(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{"email": email}

	result, err := s.db.Database(DB).Collection(Collection).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UserStore) GetAllUsers(ctx context.Context) ([]User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	cursor, err := s.db.Database(DB).Collection(Collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
