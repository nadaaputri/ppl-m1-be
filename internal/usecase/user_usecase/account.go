package user_usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) CreateAccount(ctx context.Context, in AccountInput) (Account, error) {
	roleID, divisionID, err := u.resolveRefs(ctx, in)
	if err != nil {
		return Account{}, err
	}
	hash, err := utils.HashPassword(in.Password)
	if err != nil {
		return Account{}, err
	}
	row, err := u.store.CreateAccount(ctx, toAccountParams(in, hash, roleID, divisionID))
	if utils.IsUniqueViolation(err) {
		return Account{}, ErrEmailTaken
	}
	if err != nil {
		return Account{}, err
	}
	return toAccount(row, in.Role, in.Division), nil
}

func (u *Usecase) resolveRefs(ctx context.Context, in AccountInput) (int32, int32, error) {
	roleID, roleMsg, err := u.resolveRole(ctx, in.Role)
	if err != nil {
		return 0, 0, err
	}
	divisionID, divisionMsg, err := u.resolveDivision(ctx, in.Division)
	if err != nil {
		return 0, 0, err
	}
	return roleID, divisionID, refErrors(roleMsg, divisionMsg)
}

func (u *Usecase) resolveRole(ctx context.Context, code string) (int32, string, error) {
	if code == partnerRoleCode {
		return 0, "Peran tidak dapat digunakan di sini", nil
	}
	id, err := u.store.GetGlobalRoleIDByCode(ctx, code)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "Peran tidak terdaftar", nil
	}
	return id, "", err
}

func (u *Usecase) resolveDivision(ctx context.Context, code string) (int32, string, error) {
	id, err := u.store.GetGlobalDivisionIDByCode(ctx, code)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "Divisi tidak terdaftar", nil
	}
	return id, "", err
}

func refErrors(roleMsg, divisionMsg string) error {
	fields := map[string]string{}
	if roleMsg != "" {
		fields["role"] = roleMsg
	}
	if divisionMsg != "" {
		fields["division"] = divisionMsg
	}
	if len(fields) == 0 {
		return nil
	}
	return &ValidationError{Fields: fields}
}

func toAccountParams(in AccountInput, hash string, roleID, divisionID int32) sqlc.CreateAccountParams {
	return sqlc.CreateAccountParams{
		PartnerID:    sql.NullInt32{Int32: in.PartnerID, Valid: true},
		Name:         strings.TrimSpace(in.Name),
		Email:        strings.ToLower(strings.TrimSpace(in.Email)),
		PasswordHash: hash,
		RoleID:       sql.NullInt32{Int32: roleID, Valid: true},
		DivisionID:   sql.NullInt32{Int32: divisionID, Valid: true},
		UnitUsahaID:  utils.NullInt32(in.UnitUsahaID),
		CreatedBy:    sql.NullInt32{Int32: in.CreatedBy, Valid: true},
		IpAddress:    utils.NullText(in.IP),
	}
}

func toAccount(row sqlc.CreateAccountRow, role, division string) Account {
	return Account{
		ID: row.ID, Name: row.Name, Email: row.Email, Role: &role, Division: &division,
		PartnerID: utils.Int32Ptr(row.PartnerID), UnitUsahaID: utils.Int32Ptr(row.UnitUsahaID),
		IsActive: row.IsActive, CreatedAt: row.CreatedAt,
	}
}
