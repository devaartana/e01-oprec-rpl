package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(addr string, maxOpenConnection int, maxIdleConnection int, maxIdleTime string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	idleTimeDuration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	clientOption := options.Client().ApplyURI(addr).
		SetMaxPoolSize(uint64(maxOpenConnection)). 
		SetMinPoolSize(uint64(maxIdleConnection)). 
		SetMaxConnIdleTime(idleTimeDuration)

	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
