package article

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/DroidZed/my_blog/internal/views/pages"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleService interface {
	GetOneByIDProj(ctx context.Context, id string, in utils.GetInput) (*Article, error)
	GetByID(ctx context.Context, id string) (*Article, error)
	GetByTitle(ctx context.Context, title string) (*Article, error)
	GetOne(ctx context.Context, input utils.GetInput) (*Article, error)
	ReadFileContents(articleId string) ([]byte, error)
	Add(ctx context.Context, entity Article) error
}

type Controller struct {
	service ArticleService
	logger  *slog.Logger
}

func NewController(service ArticleService, logger *slog.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  logger,
	}
}

func (c Controller) GetArticle(w http.ResponseWriter, r *http.Request) {

	title := chi.URLParam(r, "title")

	article, err := c.service.GetByTitle(r.Context(), title)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with title %s could not be found.", title)})
		return
	}

	fileContents, err := c.service.ReadFileContents(article.FileID)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with file name %s could not be found.", article.FileID)})
		return
	}

	pages.ArticleOne(string(fileContents), article.Title, article.Photo).Render(r.Context(), w)

}

func (c *Controller) AddArticle(w http.ResponseWriter, r *http.Request) {

	authorId := r.Context().Value(jwtverify.AuthCtxKey{}).(string)

	var article *Article = &Article{}

	if err := utils.DecodeBody(r, article); err != nil {
		c.logger.Error("failed to decode body", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Bad body"})
		return
	}

	article.ID = primitive.NewObjectID()
	article.AuthorId = authorId
	article.FileID = fmt.Sprintf("%s.md", utils.GenUUID())

	c.service.Add(r.Context(), *article)

	utils.JsonResponse(
		w,
		http.StatusCreated,
		utils.DtoResponse{Message: "Article saved!"},
	)

}
