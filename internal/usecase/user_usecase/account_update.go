package user_usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"ppl-m1-be/internal/usecase/repository/sqlc"
	"ppl-m1-be/internal/usecase/utils"
)

func (u *Usecase) UpdateAccount(ctx context.Context, in UpdateAccountInput) (Account, error) {
	if err := validateUpdateInput(in); err != nil {
		return Account{}, err
	}
	params, err := u.updateParams(ctx, in)
	if err != nil {
		return Account{}, err
	}
	row, err := u.store.UpdateAccount(ctx, params)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, ErrUserNotFound
	}
	if err != nil {
		return Account{}, err
	}
	return toUpdatedAccount(row), nil
}

func validateUpdateInput(in UpdateAccountInput) error {
	if in.Name == nil && in.Role == nil &&
		in.Division == nil && in.UnitUsahaID == nil {
		return &ValidationError{Fields: map[string]string{
			"body": "Minimal satu field harus diisi",
		}}
	}
	return nil
}

func (u *Usecase) updateParams(ctx context.Context, in UpdateAccountInput) (sqlc.UpdateAccountParams, error) {
	name, err := updateName(in.Name)
	if err != nil {
		return sqlc.UpdateAccountParams{}, err
	}
	roleID, divisionID, err := u.updateReferenceIDs(ctx, in)
	if err != nil {
		return sqlc.UpdateAccountParams{}, err
	}
	unitID, err := updateUnitUsahaID(in.UnitUsahaID)
	if err != nil {
		return sqlc.UpdateAccountParams{}, err
	}
	return sqlc.UpdateAccountParams{
		ID: in.TargetID, Name: name, RoleID: roleID,
		DivisionID: divisionID, UnitUsahaID: unitID,
		ActorID: in.ActorID, IpAddress: utils.NullText(in.IP),
	}, nil
}

func (u *Usecase) updateReferenceIDs(ctx context.Context, in UpdateAccountInput) (sql.NullInt32, sql.NullInt32, error) {
	roleID, err := u.updateRoleID(ctx, in.Role)
	if err != nil {
		return roleID, sql.NullInt32{}, err
	}
	divisionID, err := u.updateDivisionID(ctx, in.Division)
	return roleID, divisionID, err
}

func (u *Usecase) updateRoleID(ctx context.Context, code *string) (sql.NullInt32, error) {
	if code == nil {
		return sql.NullInt32{}, nil
	}
	if !utils.CodePattern.MatchString(*code) {
		return sql.NullInt32{}, &ValidationError{Fields: map[string]string{"role": "Kode Role tidak valid"}}
	}
	id, message, err := u.resolveRole(ctx, *code)
	if err != nil {
		return sql.NullInt32{}, err
	}
	if message != "" {
		return sql.NullInt32{}, &ValidationError{Fields: map[string]string{"role": message}}
	}
	return sql.NullInt32{Int32: id, Valid: true}, nil
}

func (u *Usecase) updateDivisionID(ctx context.Context, code *string) (sql.NullInt32, error) {
	if code == nil {
		return sql.NullInt32{}, nil
	}
	if !utils.CodePattern.MatchString(*code) {
		return sql.NullInt32{}, &ValidationError{Fields: map[string]string{"division": "Kode Division tidak valid"}}
	}
	id, message, err := u.resolveDivision(ctx, *code)
	if err != nil {
		return sql.NullInt32{}, err
	}
	if message != "" {
		return sql.NullInt32{}, &ValidationError{Fields: map[string]string{"division": message}}
	}
	return sql.NullInt32{Int32: id, Valid: true}, nil
}

func updateName(value *string) (sql.NullString, error) {
	if value == nil {
		return sql.NullString{}, nil
	}
	name := strings.TrimSpace(*value)
	if name == "" || len([]rune(name)) > 150 {
		return sql.NullString{}, &ValidationError{Fields: map[string]string{"name": "Nama wajib diisi dan maksimal 150 karakter"}}
	}
	return sql.NullString{String: name, Valid: true}, nil
}

func updateUnitUsahaID(value *int32) (sql.NullInt32, error) {
	if value == nil {
		return sql.NullInt32{}, nil
	}
	if *value < 1 {
		return sql.NullInt32{}, &ValidationError{Fields: map[string]string{"unit_usaha_id": "Nilai tidak valid"}}
	}
	return sql.NullInt32{Int32: *value, Valid: true}, nil
}

func toUpdatedAccount(r sqlc.UpdateAccountRow) Account {
	return Account{
		ID: r.ID, Name: r.Name, Email: r.Email,
		PartnerID:   utils.Int32Ptr(r.PartnerID),
		UnitUsahaID: utils.Int32Ptr(r.UnitUsahaID),
		Role:        utils.StringPtr(r.RoleCode),
		Division:    utils.StringPtr(r.DivisionCode),
		IsActive:    r.IsActive, CreatedAt: r.CreatedAt,
	}
}
