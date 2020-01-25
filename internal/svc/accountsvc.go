package svc

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
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
