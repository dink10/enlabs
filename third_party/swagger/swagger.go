package swagger

import (
	"net/http"

	// import generated swagger API
	_ "github.com/dink10/enlabs/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Documentation handles pages with swagger documentation.
// @Summary Swagger Docs
// @Description UI for swagger documentation
// @ID swagger-docs
// @Tags Swagger
// @Produce html
// @Router /swagger/index.html [get]
func Documentation() http.HandlerFunc {
	return httpSwagger.Handler(httpSwagger.URL("doc.json"))
}
