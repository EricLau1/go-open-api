package handlers

import (
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-playground/validator/v10"
	"go-open-api/services"
	"go-open-api/types"
	"go-open-api/utils"
	"net/http"
	"strings"
)

type usersHandlers struct {
	usersService services.UsersService
	validate     *validator.Validate
}

func RegisterUsersHandlers(usersService services.UsersService, ws *restful.WebService) {

	h := usersHandlers{usersService: usersService, validate: validator.New()}

	tags := []string{"users"}

	ws.Route(
		ws.POST("/v1/users").
			To(h.CreateUser).
			Doc("create an user").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Reads(types.CreateUserInput{}).
			Writes(types.User{}).
			Produces(restful.MIME_JSON).
			ReturnsWithHeaders(http.StatusCreated, "", types.User{}, map[string]restful.Header{
				"Location": {
					Items: &restful.Items{
						Type:   "string",
						Format: "uri",
					},
					Description: "path to get user by id",
				},
			}).
			Returns(http.StatusBadRequest, "", types.Error{}),
	)
	ws.Route(
		ws.PUT("/v1/users/{id}/email").
			To(h.ChangeEmail).
			Doc("update email").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(restful.PathParameter("id", "user id")).
			Reads(types.UpdateEmailInput{}).
			Writes(types.User{}).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, "", types.User{}).
			Returns(http.StatusBadRequest, "", types.Error{}),
	)
	ws.Route(
		ws.PUT("/v1/users/{id}/password").
			To(h.ChangePassword).
			Doc("update password").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(restful.PathParameter("id", "user id")).
			Reads(types.UpdatePasswordInput{}).
			Writes(types.User{}).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, "", types.User{}).
			Returns(http.StatusBadRequest, "", types.Error{}),
	)
	ws.Route(
		ws.GET("/v1/users").
			To(h.GetUsers).
			Doc("list users").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Writes([]types.User{}).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, "", []types.User{}).
			Returns(http.StatusInternalServerError, "", types.Error{}),
	)
	ws.Route(
		ws.GET("/v1/users/{id}").
			To(h.GetUser).
			Doc("get user by id").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(restful.PathParameter("id", "user id")).
			Writes(types.User{}).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, "", types.User{}).
			Returns(http.StatusBadRequest, "", types.Error{}),
	)
	ws.Route(
		ws.DELETE("/v1/users/{id}").
			To(h.DeleteUser).
			Doc("delete user by id").
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Param(restful.PathParameter("id", "user id")).
			Produces(restful.MIME_JSON).
			ReturnsWithHeaders(http.StatusNoContent, "", nil, map[string]restful.Header{
				"Entity": {
					Items: &restful.Items{
						Type:   "string",
						Format: "uuid",
					},
					Description: "deleted user id",
				},
			}).
			Returns(http.StatusBadRequest, "", types.Error{}),
	)
}

func (h *usersHandlers) CreateUser(request *restful.Request, response *restful.Response) {
	input := new(types.CreateUserInput)

	err := request.ReadEntity(input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	input.Email = utils.NormalizeEmail(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	err = h.validate.Struct(input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	output, err := h.usersService.CreateUser(request.Request.Context(), input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	response.Header().Set("Location", fmt.Sprintf("%s/%s", request.Request.URL.Path, output.ID))

	_ = response.WriteHeaderAndJson(http.StatusCreated, output, restful.MIME_JSON)
}

func (h *usersHandlers) ChangeEmail(request *restful.Request, response *restful.Response) {
	input := new(types.UpdateEmailInput)

	err := request.ReadEntity(input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	input.ID = request.PathParameter("id")
	input.Email = utils.NormalizeEmail(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	err = h.validate.Struct(input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	output, err := h.usersService.UpdateEmail(request.Request.Context(), input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	_ = response.WriteAsJson(output)
}

func (h *usersHandlers) ChangePassword(request *restful.Request, response *restful.Response) {
	input := new(types.UpdatePasswordInput)

	err := request.ReadEntity(input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	input.ID = request.PathParameter("id")
	input.CurrentPassword = strings.TrimSpace(input.CurrentPassword)
	input.NewPassword = strings.TrimSpace(input.NewPassword)

	err = h.validate.Struct(input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	output, err := h.usersService.UpdatePassword(request.Request.Context(), input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	_ = response.WriteAsJson(output)
}

func (h *usersHandlers) GetUser(request *restful.Request, response *restful.Response) {
	input := &types.GetUserInput{ID: strings.TrimSpace(request.PathParameter("id"))}

	output, err := h.usersService.GetUser(request.Request.Context(), input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	_ = response.WriteAsJson(output)
}

func (h *usersHandlers) GetUsers(request *restful.Request, response *restful.Response) {
	output, err := h.usersService.GetUsers(request.Request.Context())
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusInternalServerError, types.NewError(err), restful.MIME_JSON)
		return
	}

	_ = response.WriteAsJson(output)
}

func (h *usersHandlers) DeleteUser(request *restful.Request, response *restful.Response) {
	input := &types.DeleteUserInput{ID: strings.TrimSpace(request.PathParameter("id"))}

	err := h.usersService.DeleteUser(request.Request.Context(), input)
	if err != nil {
		_ = response.WriteHeaderAndJson(http.StatusBadRequest, types.NewError(err), restful.MIME_JSON)
		return
	}

	response.Header().Set("Entity", input.ID)
	response.WriteHeader(http.StatusNoContent)
}
