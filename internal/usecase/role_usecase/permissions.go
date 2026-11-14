package role_usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

func (u *Usecase) ListPermissions(ctx context.Context, in PermissionListInput) (PermissionListResult, error) {
    total, err := u.store.CountPermissions(ctx)
    if err != nil {
        return PermissionListResult{}, err
    }
    rows, err := u.store.ListPermissions(ctx, sqlc.ListPermissionsParams{
        PageLimit: int32(in.PerPage), PageOffset: int32((in.Page - 1) * in.PerPage),
    })
    if err != nil {
        return PermissionListResult{}, err
    }
    items := make([]Permission, 0, len(rows))
    for _, row := range rows {
        items = append(items, Permission{Code: row.Code, Module: row.Module, Description: row.Description})
    }
    return PermissionListResult{Items: items, Total: int(total)}, nil
}

func (u *Usecase) GetRolePermissions(ctx context.Context, roleID int32) ([]string, error) {
    _, err := u.store.GetGlobalRoleByID(ctx, roleID)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrRoleNotFound
    }
    if err != nil {
        return nil, err
    }
    return u.store.GetRolePermissionCodes(ctx, roleID)
}

func (u *Usecase) SetRolePermissions(ctx context.Context, roleID int32, codes []string) error {
    if _, err := u.GetRolePermissions(ctx, roleID); err != nil {
        return err
    }
    if err := u.validatePermissionCodes(ctx, codes); err != nil {
        return err
    }
    return u.store.SetRolePermissions(ctx, sqlc.SetRolePermissionsParams{
        RoleID: roleID, PermissionCodes: strings.Join(codes, ","),
    })
}

func (u *Usecase) validatePermissionCodes(ctx context.Context, codes []string) error {
    if err := checkPermissionCodes(codes); err != nil {
        return err
    }
    count, err := u.store.CountPermissionsByCodes(ctx, strings.Join(codes, ","))
    if err != nil {
        return err
    }
    if int(count) != len(codes) {
        return ErrPermissionNotFound
    }
    return nil
}

func checkPermissionCodes(codes []string) error {
    seen := make(map[string]bool, len(codes))
    for _, code := range codes {
        if code == "" || strings.TrimSpace(code) != code || strings.Contains(code, ",") {
            return ErrInvalidPermissionCode
        }
        if seen[code] {
            return ErrDuplicatePermission
        }
        seen[code] = true
    }
    return nil
}