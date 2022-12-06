package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	cb func(ctx mongo.SessionContext) (interface{}, error)
}

// will pass mongo session context to generic callback
func cbAdapter(cb func(ctx context.Context) (interface{}, error)) func(ctx mongo.SessionContext) (interface{}, error) {
	return func(ctx mongo.SessionContext) (interface{}, error) {
		return cb(ctx)
	}
}

func NewTransaction(cb func(ctx context.Context) (interface{}, error)) *Transaction {
	return &Transaction{
		cb: cbAdapter(cb),
	}
}

func (t *Transaction) Execute(ctx context.Context) (interface{}, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, err
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	return session.WithTransaction(ctx, t.cb)
}