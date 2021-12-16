package middlewares

import (
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

func CORS(container *restful.Container) {
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"Location", "Entity"},
		AllowedHeaders: []string{restful.HEADER_ContentType, restful.HEADER_Accept},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		CookiesAllowed: true,
		Container:      restful.DefaultContainer,
	}
	container.Filter(cors.Filter)
}
