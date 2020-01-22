package svc

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
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
