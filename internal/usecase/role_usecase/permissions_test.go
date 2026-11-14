package role_usecase

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"

	"ppl-m1-be/internal/usecase/repository/sqlc"
)

type permissionLookupStub struct {
	roleStore
	roleErr error
	codes   []string
}

func (s *permissionLookupStub) GetGlobalRoleByID(
	_ context.Context, id int32,
) (sqlc.Role, error) {
	return sqlc.Role{ID: id}, s.roleErr
}

func (s *permissionLookupStub) GetRolePermissionCodes(
	_ context.Context, _ int32,
) ([]string, error) {
	return s.codes, nil
}

func TestGetRolePermissionsReturnsAssignedCodes(t *testing.T) {
	codes := []string{"m1.user.view"}
	uc := New(&permissionLookupStub{codes: codes})
	got, err := uc.GetRolePermissions(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, codes) {
		t.Fatalf("permission tidak sesuai: %v", got)
	}
}

func TestGetRolePermissionsRejectsMissingRole(t *testing.T) {
	uc := New(&permissionLookupStub{roleErr: sql.ErrNoRows})
	_, err := uc.GetRolePermissions(context.Background(), 7)
	if !errors.Is(err, ErrRoleNotFound) {
		t.Fatalf("seharusnya ErrRoleNotFound, didapat %v", err)
	}
}

func TestCheckPermissionCodes(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  error
	}{
		{"kosong", []string{}, nil},
		{"valid", []string{"m1.user.view"}, nil},
		{"duplikat", []string{"m1.user.view", "m1.user.view"}, ErrDuplicatePermission},
		{"koma", []string{"m1.user,view"}, ErrInvalidPermissionCode},
		{"spasi", []string{" m1.user.view"}, ErrInvalidPermissionCode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPermissionCodes(tt.codes); !errors.Is(err, tt.want) {
				t.Fatalf("hasil %v, diharapkan %v", err, tt.want)
			}
		})
	}
}
