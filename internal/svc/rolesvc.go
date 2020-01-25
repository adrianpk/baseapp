package svc

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

func (s *Service) IndexRoles() (roles []model.Role, err error) {
	repo := s.AuthRepo
	if repo == nil {
		return roles, NoRepoErr
	}

	return repo.GetAllRoles()
}

func (s *Service) CreateRole(role *model.Role) (kbs.ValErrorSet, error) {
	// Validation
	v := NewRoleValidator(role)

	err := v.ValidateForCreate()
	if err != nil {
		return v.Errors, err
	}

	// Repo
	repo := s.AuthRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	err = repo.CreateRole(role)
	if err != nil {
		return nil, err
	}

	// Output
	return nil, nil
}

func (s *Service) GetRole(slug string) (role model.Role, err error) {
	repo := s.AuthRepo
	if err != nil {
		return role, err
	}

	role, err = repo.GetRoleBySlug(slug)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (s *Service) GetRoleByName(name string) (role model.Role, err error) {
	repo := s.AuthRepo
	if err != nil {
		return role, err
	}

	role, err = repo.GetRoleByName(name)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (s *Service) UpdateRole(slug string, role *model.Role) (kbs.ValErrorSet, error) {
	repo := s.AuthRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	// Get role
	current, err := repo.GetRoleBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create a model
	// ID shouldn't change.
	role.ID = current.ID

	// Validation
	v := NewRoleValidator(role)

	err = v.ValidateForUpdate()
	if err != nil {
		return v.Errors, err
	}

	// Update
	err = repo.UpdateRole(role)
	if err != nil {
		return v.Errors, err
	}

	// Output
	return v.Errors, nil
}

func (s *Service) DeleteRole(slug string) error {
	repo := s.AuthRepo
	if repo == nil {
		return NoRepoErr
	}

	err := repo.DeleteRoleBySlug(slug)
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
