package svc

import (
	"fmt"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

func (s *Service) IndexResources() (resources []model.Resource, err error) {
	repo := s.AuthRepo
	if repo == nil {
		return resources, NoRepoErr
	}

	return repo.GetAllResources()
}

func (s *Service) CreateResource(resource *model.Resource) (kbs.ValErrorSet, error) {
	// Validation
	v := NewResourceValidator(resource)

	err := v.ValidateForCreate()
	if err != nil {
		return v.Errors, err
	}

	// Repo
	repo := s.AuthRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	// Set tag if needed
	resource.GenTagIfEmpty()

	err = repo.CreateResource(resource)
	if err != nil {
		return nil, err
	}

	// Output
	return nil, nil
}

func (s *Service) GetResource(slug string) (resource model.Resource, err error) {
	repo := s.AuthRepo
	if err != nil {
		return resource, err
	}

	resource, err = repo.GetResourceBySlug(slug)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func (s *Service) GetResourceByName(name string) (resource model.Resource, err error) {
	repo := s.AuthRepo
	if err != nil {
		return resource, err
	}

	resource, err = repo.GetResourceByName(name)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func (s *Service) GetResourceByTag(tag string) (resource model.Resource, err error) {
	repo := s.AuthRepo
	if err != nil {
		return resource, err
	}

	resource, err = repo.GetResourceByTag(tag)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func (s *Service) GetResourceByPath(path string) (resource model.Resource, err error) {
	repo := s.AuthRepo
	if err != nil {
		return resource, err
	}

	resource, err = repo.GetResourceByPath(path)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func (s *Service) UpdateResource(slug string, resource *model.Resource) (kbs.ValErrorSet, error) {
	repo := s.AuthRepo
	if repo == nil {
		return nil, NoRepoErr
	}

	// Get resource
	current, err := repo.GetResourceBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create a model
	// ID shouldn't change.
	resource.ID = current.ID

	// Validation
	v := NewResourceValidator(resource)

	err = v.ValidateForUpdate()
	if err != nil {
		return v.Errors, err
	}

	// Update
	err = repo.UpdateResource(resource)
	if err != nil {
		return v.Errors, err
	}

	// Output
	return v.Errors, nil
}

func (s *Service) DeleteResource(slug string) error {
	repo := s.AuthRepo
	if repo == nil {
		return NoRepoErr
	}

	err := repo.DeleteResourceBySlug(slug)
	if err != nil {
		return err
	}

	// Output
	return nil
}

// Custom

// GetResourcePermissions
func (s *Service) GetResourcePermissions(resourceSlug string) (permissions []model.Permission, err error) {
	repo := s.AuthRepo
	if err != nil {
		return permissions, err
	}

	permissions, err = repo.GetResourcePermissions(resourceSlug)
	if err != nil {
		return permissions, err
	}

	return permissions, nil
}

// GetNotResourcePermissions
func (s *Service) GetNotResourcePermissions(resourceSlug string) (permissions []model.Permission, err error) {
	repo := s.AuthRepo
	if err != nil {
		return permissions, err
	}

	permissions, err = repo.GetNotResourcePermissions(resourceSlug)
	if err != nil {
		return permissions, err
	}

	return permissions, nil
}

// AppendResourcePermission
func (s *Service) AppendResourcePermission(resourceSlug, permissionSlug string) (err error) {
	authRepo := s.AuthRepo
	if err != nil {
		return err
	}

	resource, err := authRepo.GetResourceBySlug(resourceSlug)
	if err != nil {
		return err
	}

	permission, err := authRepo.GetPermissionBySlug(permissionSlug)
	if err != nil {
		return err
	}

	// FIX: work in progress; quick and dirty ResourcePermission generation.
	// Move to a model struct method.
	name := fmt.Sprintf("%s-%s", resource.Name.String, permission.Name.String)

	resourcePermission := model.ResourcePermission{
		Name:         db.ToNullString(name),
		ResourceID:   resource.ID,
		PermissionID: permission.ID,
	}

	err = authRepo.CreateResourcePermission(&resourcePermission)
	if err != nil {
		return err
	}

	return nil
}

// AppendResourcePermission
func (s *Service) RemoveResourcePermission(resourceSlug, permissionSlug string) (err error) {
	authRepo := s.AuthRepo
	if err != nil {
		return err
	}

	err = authRepo.DeleteResourcePermissionsBySlugs(resourceSlug, permissionSlug)
	if err != nil {
		return err
	}

	return nil
}
