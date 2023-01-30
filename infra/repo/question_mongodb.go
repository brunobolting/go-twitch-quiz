package repo

import (
	"context"
	"os"

	"github.com/brunobolting/go-twitch-chat/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionMongodb struct {
	client *mongo.Client
	collection string
}

func NewQuestionMongodb(client *mongo.Client) *QuestionMongodb {
	return &QuestionMongodb{
		client: client,
		collection: "questions",
	}
}

func (repo *QuestionMongodb) Create(e *entity.Question) (entity.ID, error) {
	collection := repo.getCollection()

	_, err := collection.InsertOne(context.TODO(), e)
	if err != nil {
		return e.ID, err
	}

	return e.ID, nil
}

func (repo *QuestionMongodb) Get(id entity.ID) (*entity.Question, error) {
	var question entity.Question

	collection := repo.getCollection()
	filter := bson.M{"_id": id}

	err := collection.FindOne(context.TODO(), filter).Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entity.ErrNothingFound
		}
		return nil, err
	}

	return &question, nil
}

func (repo *QuestionMongodb) GetRandom(whereNotIn []string) (*entity.Question, error) {
	var question entity.Question
	where := bson.A{}

	for _, w := range whereNotIn {
		where = append(where, w)
	}

	collection := repo.getCollection()
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"_id", bson.M{"$nin": where}}}}},
		bson.D{{"$sample", bson.D{{"size", 1}}}},
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	cursor.Next(context.TODO())

	err = cursor.Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entity.ErrNothingFound
		}
		return nil, err
	}

	return &question, nil
}

func (repo *QuestionMongodb) Update(e *entity.Question) error {
	collection := repo.getCollection()
	filter := bson.M{"_id": e.ID}

	_, err := collection.ReplaceOne(context.TODO(), filter, e)
	if err != nil {
		return err
	}

	return nil
}

func (repo *QuestionMongodb) Delete(id entity.ID) error {
	collection := repo.getCollection()
	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (repo *QuestionMongodb) getCollection() *mongo.Collection {
	return repo.client.Database(os.Getenv("DB_DATABASE")).Collection(repo.collection)
}
