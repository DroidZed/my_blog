package article

import (
	"context"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	hasher cryptor.CryptoHelper
	client *mongo.Client
	dbName string
}

func (s *Service) GetOneByIDProj(ctx context.Context, id string, in utils.GetInput) (*Article, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.GetOne(
		ctx,
		utils.GetInput{
			Filter: bson.M{"_id": objectID},
			Projection: &options.FindOneOptions{
				Projection: in.Projection,
			},
		})
}

func (s *Service) GetByID(ctx context.Context, id string) (*Article, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return s.GetOne(ctx, utils.GetInput{
		Filter: bson.M{"_id": objectID},
	})
}

func (s *Service) GetOne(ctx context.Context, input utils.GetInput) (*Article, error) {

	coll := s.client.Database(s.dbName).Collection("users")

	result := &Article{}

	if err := coll.FindOne(ctx,
		input.Filter,
		&options.FindOneOptions{
			Projection: input.Projection,
		},
	).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
