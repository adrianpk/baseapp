package svc

import (
	"fmt"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

const (
	AccountConfirmationErrMsg = "account_confirmation_err_msg"
	AccountSignInErrMsg       = "account_signin_err_msg"
)

func (s *Service) IndexAccounts() (accounts []model.Account, err error) {
	repo := s.AccountRepo
	if repo == nil {
		return accounts, NoRepoErr
	}

	return repo.GetAll()
}

func (s *Service) CreateAccount(account *model.Account) (kbs.ValErrorSet, error) {
	// Validation
	v := NewAccountValidator(account)

	err := v.ValidateForCreate()
	if err != nil {
		return v.Errors, err
	}

	// Repo
	repo := s.AccountRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	err = repo.Create(account)
	if err != nil {
		return nil, err
	}

	// Output
	return nil, nil
}

func (s *Service) GetAccount(slug string) (account model.Account, err error) {
	repo := s.AccountRepo
	if err != nil {
		return account, err
	}

	account, err = repo.GetBySlug(slug)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (s *Service) GetAccountByName(name string) (account model.Account, err error) {
	repo := s.AccountRepo
	if err != nil {
		return account, err
	}

	account, err = repo.GetByName(name)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (s *Service) UpdateAccount(slug string, account *model.Account) (kbs.ValErrorSet, error) {
	repo := s.AccountRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	// Get account
	current, err := repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create a model
	// ID shouldn't change.
	account.ID = current.ID

	// Validation
	v := NewAccountValidator(account)

	err = v.ValidateForUpdate()
	if err != nil {
		return v.Errors, err
	}

	// Update
	err = repo.Update(account)
	if err != nil {
		return v.Errors, err
	}

	// Output
	return v.Errors, nil
}

func (s *Service) DeleteAccount(slug string) error {
	repo := s.AccountRepo
	if repo == nil {
		return NoRepoErr
	}

	err := repo.DeleteBySlug(slug)
	if err != nil {
		return err
	}

	// Output
	return nil
}

// Custom

// GetAccountRoles
func (s *Service) GetAccountRoles(accountSlug string) (roles []model.Role, err error) {
	repo := s.AuthRepo
	if err != nil {
		return roles, err
	}

	roles, err = repo.GetAccountRoles(accountSlug)
	if err != nil {
		return roles, err
	}

	return roles, nil
}

// GetNotAccountRoles
func (s *Service) GetNotAccountRoles(accountSlug string) (roles []model.Role, err error) {
	repo := s.AuthRepo
	if err != nil {
		return roles, err
	}

	roles, err = repo.GetNotAccountRoles(accountSlug)
	if err != nil {
		return roles, err
	}

	return roles, nil
}

// AppendAccountRole
func (s *Service) AppendAccountRole(accountSlug, roleSlug string) (err error) {
	accountRepo := s.AccountRepo
	if err != nil {
		return err
	}

	account, err := accountRepo.GetBySlug(accountSlug)
	if err != nil {
		return err
	}

	authRepo := s.AuthRepo
	if err != nil {
		return err
	}

	role, err := authRepo.GetRoleBySlug(roleSlug)
	if err != nil {
		return err
	}

	// FIX: work in progress; quick and dirty AccountRole generation.
	// Move to a model struct method.
	name := fmt.Sprintf("%s-%s", account.Username.String, role.Name.String)

	accountRole := model.AccountRole{
		Name:      db.ToNullString(name),
		AccountID: account.ID,
		RoleID:    role.ID,
	}

	err = authRepo.CreateAccountRole(&accountRole)
	if err != nil {
		return err
	}

	return nil
}

// AppendAccountRole
func (s *Service) RemoveAccountRole(accountSlug, roleSlug string) (err error) {
	authRepo := s.AuthRepo
	if err != nil {
		return err
	}

	err = authRepo.DeleteAccountRoleBySlugs(accountSlug, roleSlug)
	if err != nil {
		return err
	}

	return nil
}
