package svc

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

func (s *Service) IndexPermissions() (permissions []model.Permission, err error) {
	repo := s.AuthRepo
	if repo == nil {
		return permissions, NoRepoErr
	}

	return repo.GetAllPermissions()
}

func (s *Service) CreatePermission(permission *model.Permission) (kbs.ValErrorSet, error) {
	// Validation
	v := NewPermissionValidator(permission)

	err := v.ValidateForCreate()
	if err != nil {
		return v.Errors, err
	}

	// Repo
	repo := s.AuthRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	err = repo.CreatePermission(permission)
	if err != nil {
		return nil, err
	}

	// Output
	return nil, nil
}

func (s *Service) GetPermission(slug string) (permission model.Permission, err error) {
	repo := s.AuthRepo
	if err != nil {
		return permission, err
	}

	permission, err = repo.GetPermissionBySlug(slug)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func (s *Service) GetPermissionByName(name string) (permission model.Permission, err error) {
	repo := s.AuthRepo
	if err != nil {
		return permission, err
	}

	permission, err = repo.GetPermissionByName(name)
	if err != nil {
		return permission, err
	}

	return permission, nil
}

func (s *Service) UpdatePermission(slug string, permission *model.Permission) (kbs.ValErrorSet, error) {
	repo := s.AuthRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	// Get permission
	current, err := repo.GetPermissionBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create a model
	// ID shouldn't change.
	permission.ID = current.ID

	// Validation
	v := NewPermissionValidator(permission)

	err = v.ValidateForUpdate()
	if err != nil {
		return v.Errors, err
	}

	// Update
	err = repo.UpdatePermission(permission)
	if err != nil {
		return v.Errors, err
	}

	// Output
	return v.Errors, nil
}

func (s *Service) DeletePermission(slug string) error {
	repo := s.AuthRepo
	if repo == nil {
		return NoRepoErr
	}

	err := repo.DeletePermissionBySlug(slug)
	if err != nil {
		return err
	}

	// Output
	return nil
}
