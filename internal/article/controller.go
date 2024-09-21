package article

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleService interface {
	GetOneByIDProj(ctx context.Context, id string, in utils.GetInput) (*Article, error)
	GetByID(ctx context.Context, id string) (*Article, error)
	GetOne(ctx context.Context, input utils.GetInput) (*Article, error)
	ReadFileContents(filename string) ([]byte, error)
	ConvertMarkdownToHTML(md []byte) []byte
	WriteFile(filename string, contents []byte) error
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

	id := chi.URLParam(r, "id")

	article, err := c.service.GetByID(r.Context(), id)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with id %s could not be found.", id)})
		return
	}

	fileContents, err := c.service.ReadFileContents(article.FileID)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with file name %s could not be found.", article.FileID)})
		return
	}

	converted := c.service.ConvertMarkdownToHTML(fileContents)

	ArticleOne(string(converted), *article).Render(r.Context(), w)

}

func (c *Controller) AddArticle(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		c.logger.Error("send article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Invalid body"})
		return
	}

	title := r.PostForm.Get("title")
	tags := strings.Split(r.PostForm.Get("tags"), ",")

	fileId := fmt.Sprintf("%s.html", utils.GenUUID())
	photoId := fmt.Sprintf("%s.webp", utils.GenUUID())

	authorId := r.PostForm.Get("authorId")

	article := &Article{
		ID:       primitive.NewObjectID(),
		Title:    title,
		Photo:    photoId,
		Tags:     tags,
		FileID:   fileId,
		AuthorId: authorId,
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.logger.Error("send article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Invalid multipart"})
		return
	}

	articleFile, _, err1 := r.FormFile("file")
	if err1 != nil {
		c.logger.Error("getting article from file", slog.String("err", err1.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Invalid file"})
		return
	}
	posterFile, _, err2 := r.FormFile("poster")
	if err2 != nil {
		c.logger.Error("getting image from file", slog.String("err", err2.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Invalid image"})
		return
	}

	var (
		fileContents   []byte
		posterContents []byte
	)

	// reading and saving article markdown
	_, err3 := articleFile.Read(fileContents)
	if err3 != nil {
		c.logger.Error("reading article", slog.String("err", err3.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Unable to read file contents"})
		return
	}
	converted := c.service.ConvertMarkdownToHTML(fileContents)
	err4 := c.service.WriteFile(fileId, converted)
	if err4 != nil {
		c.logger.Error("writing article", slog.String("err", err4.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Unable to save image"})
		return
	}


	// reading and saving article img
	_, err5 := posterFile.Read(posterContents)
	if err5 != nil {
		c.logger.Error("reading image", slog.String("err", err5.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Unable to read file contents"})
		return
	}
	err6 := c.service.WriteFile(photoId, posterContents)
	if err6 != nil {
		c.logger.Error("writing image", slog.String("err", err6.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: "Unable to save image"})
		return
	}

	c.service.Add(r.Context(), *article)

	utils.JsonResponse(
		w,
		http.StatusCreated,
		utils.DtoResponse{Message: "Article saved!"},
	)

}
