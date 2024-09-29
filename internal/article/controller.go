package article

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/DroidZed/my_blog/cmd/web/pages"
	"github.com/DroidZed/my_blog/internal/jwtverify"
	"github.com/DroidZed/my_blog/internal/user"
	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleService interface {
	GetOneByIDProj(ctx context.Context, id string, in utils.GetInput) (*Article, error)
	GetByID(ctx context.Context, id string) (*Article, error)
	GetByTitle(ctx context.Context, title string) (*Article, error)
	GetOne(ctx context.Context, input utils.GetInput) (*Article, error)
	ReadFileContents(articleId string) (*bytes.Buffer, error)
	Add(ctx context.Context, entity Article) error
}

type userProvider interface {
	Add(ctx context.Context, u user.User) error
	GetByIDProj(ctx context.Context, id string, in utils.GetInput) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetOne(ctx context.Context, in utils.GetInput) (*user.User, error)
	Validate(ctx context.Context, email, password string) (*user.User, error)
}

type Controller struct {
	service     ArticleService
	userService userProvider
	logger      *slog.Logger
}

func NewController(
	service ArticleService,
	logger *slog.Logger,
	userService userProvider,
) *Controller {
	return &Controller{
		service:     service,
		logger:      logger,
		userService: userService,
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

	buf, err := c.service.ReadFileContents(article.FileID)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with file name %s could not be found.", article.FileID)})
		return
	}

	// Create an unsafe component containing raw HTML.
	content := Unsafe(buf.String())

	pages.ArticleOne(content, article.Title, article.Photo).Render(r.Context(), w)

}

func (c *Controller) AddArticle(w http.ResponseWriter, r *http.Request) {

	authorID := r.Context().Value(jwtverify.AuthCtxKey{}).(string)

	var article *Article = &Article{}
	if err := utils.DecodeBody(r, article); err != nil {
		c.logger.Error("failed to decode body", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Bad body"})
		return
	}

	article.ID = primitive.NewObjectID()
	article.AuthorID = authorID
	article.FileID = fmt.Sprintf("%s.md", utils.GenUUID())

	c.service.Add(r.Context(), *article)

	utils.JsonResponse(
		w,
		http.StatusCreated,
		utils.DtoResponse{Message: "Article saved!"},
	)

}

func (c *Controller) GetArticleInfo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	article, err := c.service.GetByID(r.Context(), id)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with id %s could not be found.", id)})
		return
	}

	author, err := c.userService.GetByID(r.Context(), article.AuthorID)
	if err != nil {
		c.logger.Error("invalid user", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Bad body"})
		return
	}

	data := ArticleWithUser{
		Article: *article,
		Author:  *author,
	}

	utils.JsonResponse(w, http.StatusOK, data)
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
