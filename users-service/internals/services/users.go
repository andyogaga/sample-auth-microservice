package services

import (
	"fmt"
	"net/http"
	"users-service/internals/datastruct"
	"users-service/internals/dto"
	"users-service/internals/repository"
	"users-service/internals/utils"

	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	InitializeUser(*dto.InitializeUser) (*dto.CleanedUser, error)
	RegisterUser(*dto.RegisterUser) (*dto.CleanedUser, error)
	LoginUser(*dto.LoginUser) (*dto.CleanedUser, error)
	GetUser(*dto.GetUser) (*dto.CleanedUser, error)
}

type UserService struct {
	dao            repository.DAO
	profileService IProfileService
}

func NewUserService(dao repository.DAO, profileService *profileService) *UserService {
	return &UserService{
		dao:            dao,
		profileService: profileService,
	}
}

func (u *UserService) GetUser(requestedUser *dto.GetUser) (*dto.CleanedUser, error) {
	user, err := u.dao.NewUserQuery().GetCleanedUser(requestedUser)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) InitializeUser(user *dto.InitializeUser) (*dto.CleanedUser, error) {
	defer utils.RecoverFromPanic()
	dbUser, err := u.GetUser(&dto.GetUser{Phone: user.Phone})
	if dbUser != nil {
		return nil, fmt.Errorf("user already exists")
	} else if err.Error() == "user not found" {
		createdProfile, err := u.profileService.CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(user.Country)})
		if err != nil {
			utils.LogErrors(err)
			return nil, err
		}

		newUser := dto.RegisterUser{
			Phone:     user.Phone,
			ProfileId: createdProfile.ProfileId,
			Email:     "",
			Role:      datastruct.LEAD,
		}
		createdUser, err := u.dao.NewUserQuery().Create(&newUser)
		if err != nil {
			return nil, err
		}
		return createdUser, nil
	} else {
		return nil, err
	}
}

func (u *UserService) RegisterUser(user *dto.RegisterUser) (*dto.CleanedUser, error) {
	defer utils.RecoverFromPanic()

	dbUser, err := u.GetUser(&dto.GetUser{Email: user.Email, Phone: user.Phone})
	if err != nil {
		return nil, err
	}

	hashPasswordStr, err := utils.EncryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	hashedPINStr, err := utils.EncryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	if dbUser != nil && dbUser.Role == datastruct.LEAD {
		tx := u.dao.BeginTransaction()
		newUser := dto.UpdateUser{
			Role:     datastruct.USER,
			Password: hashPasswordStr,
			Phone:    user.Phone,
			Email:    user.Email,
			UserId:   dbUser.UserId,
			PIN:      hashedPINStr,
		}
		profile, err := u.profileService.GetProfile(&dto.GetProfileQuery{ProfileId: dbUser.ProfileId})
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		dbUser.Profile = *profile
		dbUser, err = u.dao.NewUserQuery().Update(&newUser)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return dbUser, nil
	} else if err != nil && err.Error() == "user not found" {
		tx := u.dao.BeginTransaction()
		createdProfile, e := u.profileService.CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(user.Country)})
		if e != nil {
			tx.Rollback()
			return nil, e
		}
		user.Password = hashPasswordStr
		user.ProfileId = createdProfile.ProfileId
		user.Role = datastruct.USER
		dbUser, e = u.dao.NewUserQuery().Create(user)
		if e != nil {
			tx.Rollback()
			return nil, err
		}
		dbUser.Profile = *createdProfile
		tx.Commit()
		return dbUser, nil
	} else if dbUser != nil && dbUser.Role != datastruct.LEAD {
		return nil, utils.NewError(http.StatusConflict, "user already exists", nil)
	} else {
		return nil, utils.NewError(http.StatusInternalServerError, "error creating user", nil)
	}
}

func (u *UserService) LoginUser(user *dto.LoginUser) (*dto.CleanedUser, error) {
	defer utils.RecoverFromPanic()

	dbUser, err := u.dao.NewUserQuery().GetSensitiveUser(&dto.GetUser{Email: user.Email, Phone: user.Phone})
	if err != nil {
		return nil, err
	}
	if dbUser != nil && (dbUser.Role == datastruct.LEAD) {
		return nil, utils.NewError(http.StatusUnauthorized, "user not fully registered, visit https://kendi.io/register to complete your registration", nil)
	} else {
		errEncrypt := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if errEncrypt != nil {
			return nil, utils.NewError(http.StatusUnauthorized, "incorrect login details", nil)
		}
		profile, err := u.profileService.GetProfile(&dto.GetProfileQuery{ProfileId: dbUser.ProfileId})
		if err != nil {
			return nil, err
		}
		cleanedUser := dto.CleanedUser{
			UserId:    dbUser.UserId,
			Phone:     dbUser.Phone,
			ProfileId: dbUser.ProfileId,
			Email:     dbUser.Email,
			Profile:   *profile,
			Role:      dbUser.Role,
		}
		return &cleanedUser, nil
	}
}
