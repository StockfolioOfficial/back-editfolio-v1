package handler

import (
	"net/http"

	"github.com/stockfolioofficial/back-editfolio/core/debug"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stockfolioofficial/back-editfolio/domain"
)

const (
	tag = "[USER] "
)

func NewUserHttpHandler(useCase domain.UserUseCase) *HttpHandler {
	return &HttpHandler{useCase: useCase}
}

type HttpHandler struct {
	useCase domain.UserUseCase
}

type CreateCustomerRequest struct {
	// Name, 길이 2~60 제한
	Name string `json:"name" validate:"required,min=2,max=60" example:"ljs"`

	// Email, 이메일 주소
	Email string `json:"email" validate:"required,email" example:"example@example.com"`

	// Mobile, 형식 : 01012345678
	Mobile string `json:"mobile" validate:"required,sf_mobile" example:"01012345678"`
} // @name CreateCustomerUserRequest

type CreatedCustomerResp struct {
	Id uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
} // @name CreatedCustomerResponse

type UpdatePasswordRequest struct {
	UserId string `json:"-" header:"User-Id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Mobile, 형식 : 01012345678
	OldPassword string `json:"oldPassword" validate:"required" example:"01012345678"`
	NewPassword string `json:"newPassword" validate:"required" example:"01087654321"`
} // @name UpdatePasswordRequest

// @Summary 고객 유저 생성
// @Description 고객 유저를 생성하는 기능
// @Accept json
// @Produce json
// @Param customerUserBody body CreateCustomerRequest true "Customer User Body"
// @Success 201 {object} CreatedCustomerResp
// @Router /user/customer [post]
func (h *HttpHandler) createCustomer(ctx echo.Context) error {
	var req CreateCustomerRequest

	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "create customer, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	newId, err := h.useCase.CreateCustomerUser(ctx.Request().Context(), domain.CreateCustomerUser{
		Name:   req.Name,
		Email:  req.Email,
		Mobile: req.Mobile,
	})

	switch err {
	case nil:
		return ctx.JSON(http.StatusCreated, CreatedCustomerResp{Id: newId})
	default:
		log.WithError(err).Error(tag, "create customer, unhandled error useCase.CreateCustomerUser")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}
}

// @Security Auth-Jwt-Bearer
// @Summary 어드민 비밀번호 수정
// @Description 어드민 유저의 비밀번호를 수정하는 API
// @Accept json
// @Produce json
// @Param updateAdminPassword body UpdatePasswordRequest true "Update Admin Password"
// @Success 204 "비밀번호 변경 성공"
// @Router /user/admin/pw [patch]
func (h *HttpHandler) updateAdminPassword(ctx echo.Context) error {
	var req UpdatePasswordRequest

	req.UserId = ctx.Request().Header.Get("User-Id")
	err := ctx.Bind(&req)
	if err != nil {
		log.WithError(err).Trace(tag, "update password, request body bind error")
		return ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: err.Error(),
		})
	}

	err = h.useCase.UpdateAdminPassword(ctx.Request().Context(), domain.UpdateAdminPassword{
		UserId:      uuid.MustParse(req.UserId),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	switch err {
	case nil:
		return ctx.NoContent(http.StatusNoContent)
	case domain.UserWrongPassword:
		return ctx.JSON(http.StatusUnauthorized, domain.UserWrongPasswordToUpdatePassword)
	case domain.ItemNotFound:
		return ctx.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	default:
		log.WithError(err).Error(tag, "update password, unhandled error useCase.UpdateAdminPassword")
		return ctx.JSON(http.StatusInternalServerError, domain.ServerInternalErrorResponse)
	}

}

func (h *HttpHandler) Bind(e *echo.Echo) {
	//CRUD, customer or admin
	e.POST("/user/customer", h.createCustomer)

	//sign, auth
	e.POST("/user/sign", h.signInUser)

	//Update Admin Password
	e.PATCH("/user/admin/pw", h.updateAdminPassword, debug.JwtBypassOnDebug())
}
