package svc

import (
	"errors"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

var (
	noUserRepoErr = errors.New("no user repo")
)

func (s *Service) IndexUsers() (users []model.User, err error) {
	repo := s.UserRepo
	if repo == nil {
		return users, noUserRepoErr
	}

	return repo.GetAll()
}

func (s *Service) CreateUser(user *model.User) (kbs.ValErrorSet, error) {
	// Validation
	v := NewUserValidator(user)

	err := v.ValidateForCreate()
	if err != nil {
		return v.Errors, err
	}

	// Confirmation
	user.GenAutoConfirmationToken()

	// Repo
	repo := s.UserRepo
	if repo == nil {
		return nil, noUserRepoErr
	}

	err = repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Output
	return nil, nil
}

func (s *Service) GetUser(slug string) (user model.User, err error) {
	repo := s.UserRepo
	if err != nil {
		return user, err
	}

	user, err = repo.GetBySlug(slug)
	if err != nil {
		return user, err
	}

	return user, nil
}

//func (s *Service) GetUserByUsernamei(username string) (user model.User, err error) {
//repo, err := s.UserRepo
//if repo == nil {
//return noUserRepoErr
//}

//u, err = repo.GetByUsername(u.Username.String)
//if err != nil {
//res.FromModel(nil, getUserErr, err)
//return err
//}

//err = repo.Commit()
//if err != nil {
//res.FromModel(nil, getUserErr, err)
//return err
//}

//// Output
//res.FromModel(&u, okResultInfo, nil)
//return nil
//}

func (s *Service) UpdateUser(slug string, user *model.User) (kbs.ValErrorSet, error) {
	repo := s.UserRepo
	if repo == nil {
		return nil, noUserRepoErr
	}

	// Get user
	current, err := repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create a model
	// ID shouldn't change.
	user.ID = current.ID
	// Username can change if system enabled.
	// Set envar GRN_APP_USERNAME_UPDATABLE=true
	// to let username be updatable.
	if !(s.Cfg.ValAsBool("kbs.username.updatable", false)) {
		user.Username = current.Username
	}

	// Validation
	v := NewUserValidator(user)

	err = v.ValidateForUpdate()
	if err != nil {
		return v.Errors, err
	}

	// Update
	err = repo.Update(user)
	if err != nil {
		return v.Errors, err
	}

	// Output
	return v.Errors, nil
}

//func (s *Service) DeleteUser(req tp.DeleteUserReq, res *tp.DeleteUserRes) error {
//repo, err := s.UserRepo
//if repo == nil {
//return noUserRepoErr
//}

//err = repo.DeleteBySlug(req.Slug)
//if err != nil {
//res.FromModel(deleteUserErr, err)
//return err
//}

//err = repo.Commit()
//if err != nil {
//res.FromModel(deleteUserErr, err)
//return err
//}

//// Output
//res.FromModel(okResultInfo, nil)
//return nil
//}

//func (s *Service) SignUpUser(req tp.SignUpUserReq, res *tp.SignUpUserRes) error {
//// Model
//u := req.ToModel()

//// Validation
//v := NewUserValidator(u)

//err := v.ValidateForSignUp()
//if err != nil {
//res.FromModel(&u, validationErr, err)
//}

//// Generate confirmation token
//u.GenConfirmationToken()

//repo, err := s.UserRepo
//if repo == nil {
//return noUserRepoErr
//}

//err = repo.Create(&u)
//if err != nil {
//res.FromModel(&u, cannotProcErr, err)
//return err
//}

//err = repo.Commit()
//if err != nil {
//res.FromModel(&u, createUserErr, err)
//return err
//}

//// Mail confirmation
//s.sendConfirmationEmail(&u)

//// Output
//res.FromModel(&u, okResultInfo, nil)
//return nil
//}

//func (s *Service) ConfirmUser(req tp.GetUserReq, res *tp.GetUserRes) error {
//// Model
//u := req.ToModel()

//repo, err := s.UserRepo
//if repo == nil {
//return noUserRepoErr
//}

//s.Log.Debug("Values", "slug", u.Slug.String, "token", u.ConfirmationToken.String)

//u, err = repo.GetBySlugAndToken(u.Slug.String, u.ConfirmationToken.String)
//if err != nil {
//res.FromModel(&u, confirmationErr, err)
//return err
//}

//if u.IsConfirmed.Bool {
//res.FromModel(&u, alreadyConfirmedErr, err)
//return errors.New("already confirmed")
//}

//u, err = repo.ConfirmUser(u.Slug.String, u.ConfirmationToken.String)
//if err != nil {
//res.FromModel(&u, confirmationErr, err)
//return err
//}

//err = repo.Commit()
//if err != nil {
//res.FromModel(&u, confirmationErr, err)
//return err
//}

//// Output
//res.FromModel(&u, okResultInfo, nil)
//return nil
//}

//func (s *Service) SignInUser(req tp.SignInUserReq, res *tp.SignInUserRes) error {
//// Model
//u := req.ToModel()

//repo, err := s.UserRepo
//if repo == nil {
//return noUserRepoErr
//}

//u, err = repo.SignIn(u.Username.String, u.Password)
//if err != nil {
//res.FromModel(&u, signinErr, err)
//return err
//}

//err = repo.Commit()
//if err != nil {
//res.FromModel(&u, signinErr, err)
//return err
//}

//// Output
//res.FromModel(&u, okResultInfo, nil)
//return nil
//}
