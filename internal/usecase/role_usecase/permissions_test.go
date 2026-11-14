package role_usecase

import (
	"errors"
	"testing"
)

func TestCheckPermissionCodes(t *testing.T) {
    tests := []struct {
        name string
        codes []string
        want error
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