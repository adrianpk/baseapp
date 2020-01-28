package svc

func (s *Service) ResourcePermissionTagsByPath(path string) (tags []string, err error) {
	repo := s.AuthRepo
	if repo == nil {
		return tags, NoRepoErr
	}

	tags, err = repo.GetResourcePermissionTagsByPath(path)
	if err != nil {
		return tags, err
	}

	// Output
	return tags, nil
}
