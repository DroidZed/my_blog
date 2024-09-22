package article

import (
	"context"
	"os"
	"path/filepath"

	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	client *mongo.Client
	dbName string
}

func NewService(
	client *mongo.Client,
	dbName string,
) *Service {
	return &Service{
		client: client,
		dbName: dbName,
	}
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

	coll := s.client.Database(s.dbName).Collection("articles")

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

func (s *Service) Add(ctx context.Context, entity Article) error {

	coll := s.client.Database(s.dbName).Collection("articles")

	_, insertErr := coll.InsertOne(ctx, entity)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

func (s *Service) GetByTitle(ctx context.Context, title string) (*Article, error) {
	return s.GetOne(ctx, utils.GetInput{
		Filter: bson.M{"title": title},
	})
}

func (s *Service) ReadFileContents(articleId string) ([]byte, error) {
	workDir, _ := os.Getwd()

	b, err := os.ReadFile(string(filepath.Join(workDir, "internal/asset/static/markdown", articleId)))

	if err != nil {
		return nil, err
	}

	parsed := markdown.ToHTML(b, nil, nil)

	html := bluemonday.UGCPolicy().SanitizeBytes(parsed)

	return html, nil
}
