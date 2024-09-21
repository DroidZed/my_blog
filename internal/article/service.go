package article

import (
	"context"
	"fmt"
	"os"

	"github.com/DroidZed/my_blog/internal/cryptor"
	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
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

func NewService(
	hasher cryptor.CryptoHelper,
	client *mongo.Client,
	dbName string,
) *Service {
	return &Service{
		hasher: hasher,
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

func (s *Service) ReadFileContents(filename string) ([]byte, error) {

	b, err := os.ReadFile(fmt.Sprintf("/home/blog/%s", filename))

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) ConvertMarkdownToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	html := markdown.Render(doc, renderer)

	return bluemonday.UGCPolicy().SanitizeBytes(html)
}

func (s *Service) WriteFile(filename string, contents []byte) error {

	file, err := os.Create(fmt.Sprintf("/home/blog/%s", filename))

	if err != nil {
		return err
	}

	_, err2 := file.Write(contents)

	if err2 != nil {
		return err2
	}

	return nil
}
