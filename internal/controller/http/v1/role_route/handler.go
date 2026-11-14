package role_route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ppl-m1-be/internal/usecase/role_usecase"
	"ppl-m1-be/pkg/httpserver"
)

type Handler struct {
	uc *role_usecase.Usecase
}

func (h *Handler) List(c *gin.Context) {
	in, ok := listInput(c)
	if !ok {
		return
	}
	result, err := h.uc.List(c.Request.Context(), in)
	if err != nil {
		h.writeError(c, err)
		return
	}
	httpserver.SuccessList(c, toRoleList(result.Items), in.Page, in.PerPage, result.Total)
}

func (h *Handler) Create(c *gin.Context) {
    req, ok := bindCreate(c)
    if !ok {
        return
    }
    role, err := h.uc.Create(c.Request.Context(), role_usecase.CreateInput{
        Code: req.Code, Name: req.Name, Description: req.Description,
        PermissionCodes: req.PermissionCodes,
    })
    if err != nil {
        h.writeError(c, err)
        return
    }
    httpserver.Success(c, http.StatusCreated, toRoleResponse(role))
}

func listInput(c *gin.Context) (role_usecase.ListInput, bool) {
    page, perPage, ok := httpserver.ParsePage(c)
    if !ok {
        return role_usecase.ListInput{}, false
    }
    active, ok := httpserver.QueryBool(c, "is_active")
    return role_usecase.ListInput{
        IsActive: active, Page: page, PerPage: perPage,
    }, ok
}

func bindCreate(c *gin.Context) (createRequest, bool) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpserver.BindError(c, err)
		return req, false
	}
	if fields := extraFieldErrors(req); len(fields) > 0 {
		httpserver.FailFields(c, httpserver.MsgValidation, fields)
		return req, false
	}
	return req, true
}

func (h *Handler) writeError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, role_usecase.ErrCodeReserved):
        httpserver.FailFields(c, httpserver.MsgValidation, map[string]string{"code": "Kode ini dicadangkan sistem"})
    case errors.Is(err, role_usecase.ErrCodeTaken):
        httpserver.Fail(c, http.StatusConflict, httpserver.CodeRoleCodeDuplicate, "Kode peran sudah dipakai")
    case errors.Is(err, role_usecase.ErrRoleNotFound):
        httpserver.Fail(c, http.StatusNotFound, httpserver.CodeNotFound, "Peran tidak ditemukan")
    case errors.Is(err, role_usecase.ErrPermissionNotFound),
        errors.Is(err, role_usecase.ErrDuplicatePermission),
        errors.Is(err, role_usecase.ErrInvalidPermissionCode):
        h.writePermissionError(c, err)
    default:
        slog.Error("operasi peran gagal", "error", err)
        httpserver.Fail(c, http.StatusInternalServerError, httpserver.CodeInternal, "Terjadi kesalahan pada server")
    }
}

func (h *Handler) writePermissionError(c *gin.Context, err error) {
    message := "Permission tidak terdaftar"
    if errors.Is(err, role_usecase.ErrDuplicatePermission) {
        message = "Terdapat kode permission duplikat"
    }
    if errors.Is(err, role_usecase.ErrInvalidPermissionCode) {
        message = "Format kode permission tidak valid"
    }
    httpserver.FailFields(c, httpserver.MsgValidation,
        map[string]string{"permission_codes": message})
}

func (h *Handler) ListPermissions(c *gin.Context) {
    page, perPage, ok := httpserver.ParsePage(c)
    if !ok {
        return
    }
    result, err := h.uc.ListPermissions(c.Request.Context(),
        role_usecase.PermissionListInput{Page: page, PerPage: perPage})
    if err != nil {
        h.writeError(c, err)
        return
    }
    httpserver.SuccessList(c, toPermissionList(result.Items), page, perPage, result.Total)
}

func (h *Handler) GetRolePermissions(c *gin.Context) {
    id, ok := httpserver.ParamInt32(c, "id")
    if !ok {
        return
    }
    codes, err := h.uc.GetRolePermissions(c.Request.Context(), id)
    if err != nil {
        h.writeError(c, err)
        return
    }
    httpserver.Success(c, http.StatusOK, toRolePermissionsResponse(id, codes))
}

func (h *Handler) SetRolePermissions(c *gin.Context) {
    id, ok := httpserver.ParamInt32(c, "id")
    if !ok {
        return
    }
    codes, ok := bindPermissionCodes(c)
    if !ok {
        return
    }
    if err := h.uc.SetRolePermissions(c.Request.Context(), id, codes); err != nil {
        h.writeError(c, err)
        return
    }
    httpserver.Success(c, http.StatusOK, toRolePermissionsResponse(id, codes))
}

func bindPermissionCodes(c *gin.Context) ([]string, bool) {
    var req setPermissionsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpserver.BindError(c, err)
        return nil, false
    }
    if req.PermissionCodes == nil {
        httpserver.FailFields(c, httpserver.MsgValidation,
            map[string]string{"permission_codes": "Isian ini wajib diisi"})
        return nil, false
    }
    return *req.PermissionCodes, true
}