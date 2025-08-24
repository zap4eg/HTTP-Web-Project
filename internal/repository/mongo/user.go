package mongo

import (
	"WebProject/internal/core"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (repository *UserRepository) GetAll(ctx context.Context) ([]*core.User, error) {
	cursor, err := repository.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	users := make([]*core.User, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (repository *UserRepository) GetById(ctx context.Context, id string) (*core.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userChannel := make(chan *core.User)

	var err error
	go func() {
		err = repository.retrieveUser(ctx, id, userChannel)
	}()
	if err != nil {
		return nil, err
	}

	select {
	case user := <-userChannel:
		return user, nil
	case <-ctxTimeout.Done():
		{
			return nil, ctxTimeout.Err()
		}
	}
}

func (repository *UserRepository) retrieveUser(ctx context.Context, id string, channel chan<- *core.User) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	user := &core.User{}

	filter := bson.M{"_id": objectId}
	err = repository.collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return err
	}

	channel <- user
	return nil
}

func (repository *UserRepository) Save(ctx context.Context, user *core.User) (*core.User, error) {
	result, err := repository.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}
