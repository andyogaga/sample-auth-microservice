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
	InitializeUser(user *dto.InitializeUser) (*datastruct.User, error)
	RegisterUser(user *dto.RegisterUser) (*datastruct.User, error)
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func (u *userService) GetUser(userID string) (*datastruct.User, error) {
	user, err := u.dao.NewUserQuery().Get(&dto.GetUser{UserId: &userID})
	if err != nil {
		log.Printf("user isn't authorized %v", err)
		return nil, err
	}
	return user, nil
}

func (u *userService) InitializeUser(user *dto.InitializeUser) (*datastruct.User, error) {
	defer utils.RecoverFromPanic()
	dbUser, err := u.dao.NewUserQuery().Get(&dto.GetUser{Phone: &user.Phone})
	if dbUser != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	} else if err.Error() == "user not found" {
		createdProfile, err := u.dao.NewProfileQuery().CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(user.Country)})
		if err != nil {
			return nil, err
		}

		role := datastruct.LEAD
		var email = ""
		newUser := dto.RegisterUser{
			Phone:     user.Phone,
			ProfileId: createdProfile.ProfileId,
			Email:     &email,
			Role:      &role,
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

func (u *userService) RegisterUser(user *dto.RegisterUser) (*datastruct.User, error) {
	defer utils.RecoverFromPanic()

	dbUser, err := u.dao.NewUserQuery().Get(&dto.GetUser{Email: user.Email, Phone: &user.Phone})

	hashedPasswordByte, errEncrypt := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if errEncrypt != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}
	hashPasswordStr := string(hashedPasswordByte)
	if dbUser != nil && dbUser.Role == datastruct.LEAD {

		newUser := dto.UpdateUser{
			Role:     datastruct.USER,
			Password: &hashPasswordStr,
			Phone:    &user.Phone,
			Email:    user.Email,
			UserId:   &dbUser.UserId,
		}
		dbUser, e := u.dao.NewUserQuery().Update(&newUser)
		if e != nil {
			return nil, err
		}
		return dbUser, nil
	} else if err != nil && err.Error() == "user not found" {
		createdProfile, e := u.dao.NewProfileQuery().CreateProfile(&dto.CreateProfile{Country: datastruct.Countries(*user.Country)})
		if e != nil {
			return nil, err
		}
		role := datastruct.USER
		uuid := utils.GenerateUUID()
		user.UserId = &uuid
		user.Password = &hashPasswordStr
		user.ProfileId = createdProfile.ProfileId
		user.Role = &role
		dbUser, e = u.dao.NewUserQuery().Create(user)
		if e != nil {
			return nil, err
		}
		return dbUser, nil
	} else if dbUser != nil && dbUser.Role != datastruct.LEAD {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	} else {
		return nil, status.Error(codes.AlreadyExists, "error creating user")
	}
}
