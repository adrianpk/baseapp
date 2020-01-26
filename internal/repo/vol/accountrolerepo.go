package pg

import (
	"errors"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	accountRoleRow struct {
		mutable bool
		model   model.AccountRole
	}
)

var (
	accountRole1 = model.AccountRole{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("9e7f5355-a8a5-46b7-a3b8-4ddc26c9386b"),
			Slug: db.ToNullString("accountRole1-bbc4116229c6"),
		},
		AccountID: kbs.ToUUID("e8b43223-17fe-4e36-bd0f-a7d96e867d95"), // userRes
		RoleID:    kbs.ToUUID("288bb973-2196-4007-808a-d7844ecf4dd9"), // permission1
	}

	accountRole2 = model.AccountRole{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("de90dce3-1c33-4d79-9dfa-a06fbb7d7c00"),
			Slug: db.ToNullString("accountRole2-fd3e9d6aa641"),
		},
		AccountID: kbs.ToUUID("fc86c00c-2d4f-400b-ae57-d9d5c87d13c8"), // accountRes
		RoleID:    kbs.ToUUID("d0d6bc3a-38b0-4a00-83c0-516d2514d7b5"), // permission2
	}

	accountRolesTable = map[uuid.UUID]accountRoleRow{
		accountRole1.ID: accountRoleRow{mutable: false, model: accountRole1},
		accountRole2.ID: accountRoleRow{mutable: false, model: accountRole2},
	}
)

func (ar *AuthRepo) CreateAccountRole(accountRole *model.AccountRole, tx ...*sqlx.Tx) error {
	_, ok := accountRolesTable[accountRole.ID]
	if ok {
		errors.New("duplicate key violates unique constraint")
	}

	if accountRole.ID == uuid.Nil {
		errors.New("Non valid primary key")
	}

	accountRolesTable[accountRole.ID] = accountRoleRow{
		mutable: true,
		model:   *accountRole,
	}

	return nil
}

func (ar *AuthRepo) GetAllAccountRoles() (accountRoles []model.AccountRole, err error) {
	size := len(accountRolesTable)
	out := make([]model.AccountRole, size)
	for _, row := range accountRolesTable {
		out = append(out, row.model)
	}
	return out, nil
}

func (ar *AuthRepo) GetAccountRole(id uuid.UUID) (accountRole model.AccountRole, err error) {
	for _, row := range accountRolesTable {
		if id == row.model.ID {
			return row.model, nil
		}
	}
	return model.AccountRole{}, nil
}

func (ar *AuthRepo) GetAccountRoleBySlug(slug string) (accountRole model.AccountRole, err error) {
	for _, row := range accountRolesTable {
		if slug == row.model.Slug.String {
			return row.model, nil
		}
	}
	return model.AccountRole{}, nil
}

func (ar *AuthRepo) GetAccountRolesByAccountID(accountID uuid.UUID) (accountRoles []model.AccountRole, err error) {
	size := len(accountRolesTable)
	accountRoles = make([]model.AccountRole, size)
	for _, row := range accountRolesTable {
		if accountID == row.model.AccountID {
			accountRoles = append(accountRoles, row.model)
		}
	}
	return accountRoles, nil
}

func (ar *AuthRepo) GetAccountRolesByRoleID(roleID uuid.UUID) (accountRoles []model.AccountRole, err error) {
	size := len(accountRolesTable)
	accountRoles = make([]model.AccountRole, size)
	for _, row := range accountRolesTable {
		if roleID == row.model.RoleID {
			accountRoles = append(accountRoles, row.model)
		}
	}
	return accountRoles, nil
}

func (ar *AuthRepo) UpdateAccountRole(accountRole *model.AccountRole, tx ...*sqlx.Tx) error {
	for _, row := range accountRolesTable {
		if accountRole.ID == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			accountRolesTable[accountRole.ID] = accountRoleRow{
				mutable: true,
				model:   *accountRole,
			}
			return nil
		}
	}
	return errors.New("no records updated")
}

func (ar *AuthRepo) DeleteAccountRole(id uuid.UUID, tx ...*sqlx.Tx) error {
	for _, row := range accountRolesTable {
		if id == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(accountRolesTable, id)
			return nil
		}
	}
	return errors.New("no records deleted")
}

func (ar *AuthRepo) DeleteAccountRoleBySlug(slug string, tx ...*sqlx.Tx) error {
	for _, row := range accountRolesTable {
		if slug == row.model.Slug.String {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(accountRolesTable, row.model.ID)
			return nil
		}
	}
	return errors.New("no records deleted")
}

func (ar *AuthRepo) DeleteAccountRoleBySlugs(accountSlug, roleSlug string, tx ...*sqlx.Tx) error {
	panic("not implemented")
}
