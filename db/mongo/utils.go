package mongo

import (
	types "RuntimeError/types/mongo"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbName = "RuntimeError"

type ModelLabel string

const (
	UserLabel     ModelLabel = "user"
	QuestionLabel ModelLabel = "question"
)

func decodeSingleResult(result *mongo.SingleResult, label ModelLabel) (any, error) {
	switch label {
	case UserLabel:
		user := &types.User{}
		if err := result.Decode(user); err != nil {
			return nil, err
		}
		return user, nil
	case QuestionLabel:
		question := &types.Question{}
		if err := result.Decode(question); err != nil {
			return nil, err
		}
		return question, nil
	}

	return nil, errors.New("unknown label")
}

func decodeCursor(ctx context.Context, cursor *mongo.Cursor, label ModelLabel) (any, error) {
	defer func() {
		if err := cursor.Err(); err != nil {
			//TODO
		}
		if err := cursor.Close(ctx); err != nil {
			//TODO
		}
	}()

	switch label {
	case UserLabel:
		var users []*types.User
		for cursor.Next(ctx) {
			user := &types.User{}
			if err := cursor.Decode(user); err != nil {
				return nil, err
			}

			users = append(users, user)
		}

		return users, nil
	case QuestionLabel:
		var questions []*types.Question
		for cursor.Next(ctx) {
			question := &types.Question{}
			if err := cursor.Decode(question); err != nil {
				return nil, err
			}

			questions = append(questions, question)
		}
		return questions, nil
	}

	return nil, errors.New("unknown label")
}

func GetOne(ctx context.Context, collName string, filter any, label ModelLabel) (any, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, err
	}

	col := client.Database(dbName).Collection(collName)
	return decodeSingleResult(col.FindOne(ctx, filter), label)
}

func GetAll(ctx context.Context, collName string, filter any, label ModelLabel) (any, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, err
	}

	col := client.Database(dbName).Collection(collName)
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return decodeCursor(ctx, cursor, label)
}

// Insert will return hex id of inserted object
func Insert(ctx context.Context, collName string, obj any, label ModelLabel) (string, error) {
	client, err := getMongoClient()
	if err != nil {
		return "", err
	}

	col := client.Database(dbName).Collection(collName)
	res, err := col.InsertOne(ctx, obj)
	if err != nil {
		return "", err
	}

	id, ok := res.InsertedID.(string)
	if !ok {
		return "", errors.New("failed to retrieve id")
	}

	return id, nil
}

func Update(ctx context.Context, collName string, filter any, obj any) (string, error) {
	client, err := getMongoClient()
	if err != nil {
		return "", err
	}

	col := client.Database(dbName).Collection(collName)
	res, err := col.ReplaceOne(ctx, filter, obj)
	if err != nil {
		return "", err
	}
	return res.UpsertedID.(string), nil
}

// Delete expects a hex id
func Delete(ctx context.Context, collName string, filter any, label ModelLabel) (int64, error) {
	client, err := getMongoClient()
	if err != nil {
		return 0, err
	}

	col := client.Database(dbName).Collection(collName)
	count, err := col.DeleteMany(ctx, filter)

	return count.DeletedCount, err
}

// Delete expects a hex id
func DeleteById(ctx context.Context, collName string, id string, label ModelLabel) error {
	_, err := Delete(ctx, collName, bson.M{"_id": id}, label)
	return err
}
