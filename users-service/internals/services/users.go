package services

import (
	"log"

	"users-service/internals/datastruct"
	"users-service/internals/dto"
	"users-service/internals/repository"
	"users-service/internals/utils"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	InitializeUser(user *dto.InitializeUser) (*dto.CleanedUser, error)
	RegisterUser(user *dto.RegisterUser) (*dto.CleanedUser, error)
	LoginUser(user *dto.LoginUser) (*dto.CleanedUser, error)
}

type userService struct {
	dao            repository.DAO
	profileService ProfileService
}

func NewUserService(dao repository.DAO, profileService ProfileService) UserService {
	return &userService{
		dao:            dao,
		profileService: profileService,
	}
}

func (u *userService) GetUser(userID string) (*dto.CleanedUser, error) {
	user, err := u.dao.NewUserQuery().GetCleanedUser(&dto.GetUser{UserId: userID})
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return user, nil
}

func (u *userService) InitializeUser(user *dto.InitializeUser) (*dto.CleanedUser, error) {
	defer utils.RecoverFromPanic()
	dbUser, err := u.dao.NewUserQuery().GetCleanedUser(&dto.GetUser{Phone: user.Phone})
	if dbUser != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	} else if err.Error() == "user not found" {
		createdProfile, err := u.dao.NewProfileQuery().CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(user.Country)})
		if err != nil {
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
		return nil, status.Errorf(codes.AlreadyExists, "error creating user")
	}
}

func (u *userService) RegisterUser(user *dto.RegisterUser) (*dto.CleanedUser, error) {
	defer utils.RecoverFromPanic()

	dbUser, err := u.dao.NewUserQuery().GetCleanedUser(&dto.GetUser{Email: user.Email, Phone: user.Phone})

	hashedPasswordByte, errEncrypt := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errEncrypt != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	hashPasswordStr := string(hashedPasswordByte)
	if dbUser != nil && dbUser.Role == datastruct.LEAD {

		newUser := dto.UpdateUser{
			Role:     datastruct.USER,
			Password: hashPasswordStr,
			Phone:    user.Phone,
			Email:    user.Email,
			UserId:   dbUser.UserId,
		}
		profile, err := u.profileService.GetProfile(&dto.GetProfileQuery{ProfileId: dbUser.ProfileId})
		if err != nil {
			return nil, err
		}
		dbUser.Profile = *profile
		dbUser, e := u.dao.NewUserQuery().Update(&newUser)
		if e != nil {
			return nil, err
		}
		return dbUser, nil
	} else if err != nil && err.Error() == "user not found" {
		createdProfile, e := u.dao.NewProfileQuery().CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(user.Country)})
		if e != nil {
			return nil, err
		}
		user.Password = hashPasswordStr
		user.ProfileId = createdProfile.ProfileId
		user.Role = datastruct.USER
		dbUser, e = u.dao.NewUserQuery().Create(user)
		if e != nil {
			return nil, err
		}
		dbUser.Profile = *createdProfile
		return dbUser, nil
	} else if dbUser != nil && dbUser.Role != datastruct.LEAD {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	} else {
		return nil, status.Error(codes.AlreadyExists, "error creating user")
	}
}

func (u *userService) LoginUser(user *dto.LoginUser) (*dto.CleanedUser, error) {
	defer utils.RecoverFromPanic()

	dbUser, err := u.dao.NewUserQuery().GetSensitiveUser(&dto.GetUser{Email: user.Email, Phone: user.Phone})
	if err != nil {
		return nil, err
	}
	if dbUser != nil && (dbUser.Role == datastruct.LEAD) {
		return nil, status.Error(codes.PermissionDenied, "user not fully registered, visit https://kendi.io/register to complete your registration")
	} else {
		errEncrypt := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if errEncrypt != nil {
			return nil, status.Error(codes.PermissionDenied, "incorrect login details")
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
