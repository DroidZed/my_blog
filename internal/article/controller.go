package article

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/DroidZed/my_blog/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

type ArticleManager interface {
	GetOneByIDProj(ctx context.Context, id string, in utils.GetInput) (*Article, error)
	GetByID(ctx context.Context, id string) (*Article, error)
	GetOne(ctx context.Context, input utils.GetInput) (*Article, error)
}

type Controller struct {
	service ArticleManager
	logger  *slog.Logger
}

func New(service ArticleManager, logger *slog.Logger) *Controller {
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

	fileContents, err := readFileContents(article.FileID)

	if err != nil {
		c.logger.Error("failed to find article", slog.String("err", err.Error()))
		utils.JsonResponse(w, http.StatusNotFound, utils.DtoResponse{Error: fmt.Sprintf("Article with file name %s could not be found.", article.FileID)})
		return
	}

	converted := convertMarkdownToHTML(fileContents)

	ArticleOne(string(converted), *article).Render(r.Context(), w)

}

func readFileContents(filename string) ([]byte, error) {

	b, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func convertMarkdownToHTML(md []byte) []byte {
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
