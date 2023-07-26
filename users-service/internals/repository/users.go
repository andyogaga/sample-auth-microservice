package repository

import (
	"fmt"
	datastruct "users-service/internals/datastruct"
	dto "users-service/internals/dto"
	utils "users-service/internals/utils"

	gorm "gorm.io/gorm"
)

type UsersQuery interface {
	GetSensitiveUser(*dto.GetUser) (*datastruct.User, error)
	GetCleanedUser(*dto.GetUser) (*dto.CleanedUser, error)
	Create(*dto.RegisterUser) (*dto.CleanedUser, error)
	Update(*dto.UpdateUser) (*dto.CleanedUser, error)
}

const (
	PHONE   = "phone"
	EMAIL   = "email"
	USER_ID = "user_id"
)

type usersQuery struct {
	PostgresDB *gorm.DB
}

func (u *usersQuery) getUserById(userID string) (*datastruct.User, error) {
	var userModel datastruct.User
	user := u.PostgresDB.Raw("SELECT * FROM users WHERE users.user_id = ? LIMIT 1", userID).Scan(&userModel)
	if user.Error != nil || user.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &userModel, nil
}

func (u *usersQuery) getUserByEmail(email string) (*datastruct.User, error) {
	var userModel datastruct.User
	user := u.PostgresDB.Raw("SELECT * FROM users WHERE users.email = ? LIMIT 1", email).Scan(&userModel)
	if user.Error != nil || user.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &userModel, nil
}

func (u *usersQuery) getUserByPhone(phone string) (*datastruct.User, error) {
	var userModel datastruct.User
	user := u.PostgresDB.Raw("SELECT * FROM users WHERE users.phone = ? LIMIT 1", phone).Scan(&userModel)
	if user.Error != nil || user.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &userModel, nil
}

func (u *usersQuery) GetSensitiveUser(user *dto.GetUser) (*datastruct.User, error) {
	if user.Phone != "" {
		return u.getUserByPhone(user.Phone)
	}
	if user.UserId != "" {
		return u.getUserById(user.UserId)
	}
	if user.Email != "" {
		return u.getUserByEmail(user.Email)
	}
	return nil, fmt.Errorf("user not found")
}

func (u *usersQuery) GetCleanedUser(user *dto.GetUser) (*dto.CleanedUser, error) {
	dbUser, err := u.GetSensitiveUser(user)
	if err != nil {
		return nil, err
	}
	cleanUser := dto.CleanedUser{
		UserId:    dbUser.UserId,
		Phone:     dbUser.Phone,
		ProfileId: dbUser.ProfileId,
		Email:     dbUser.Email,
		Profile:   dbUser.Profile,
		Role:      dbUser.Role,
	}
	return &cleanUser, nil
}

func (u *usersQuery) Create(initUser *dto.RegisterUser) (*dto.CleanedUser, error) {
	user := datastruct.User{
		UserId:    utils.GenerateUUID(),
		Phone:     initUser.Phone,
		Verified:  false,
		ProfileId: initUser.ProfileId,
		Email:     initUser.Email,
		Role:      initUser.Role,
		Password:  initUser.Password,
		PIN:       initUser.PIN,
	}
	result := u.PostgresDB.Create(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("error creating user")
	}
	cleanUser := dto.CleanedUser{
		UserId:    user.UserId,
		Phone:     user.Phone,
		Verified:  false,
		ProfileId: user.ProfileId,
		Email:     user.Email,
		Role:      user.Role,
	}
	return &cleanUser, nil
}

func (u *usersQuery) Update(user *dto.UpdateUser) (*dto.CleanedUser, error) {
	dbUser, err := u.GetSensitiveUser(&dto.GetUser{UserId: user.UserId})
	if err != nil {
		utils.LogErrors(err)
		return nil, fmt.Errorf("error updating user")
	}
	if user.Email != "" {
		dbUser.Email = user.Email
	}
	if user.Phone != "" {
		dbUser.Phone = user.Phone
	}
	if user.Role != "" {
		dbUser.Role = user.Role
	}
	if user.Password != "" {
		dbUser.Password = user.Password
	}
	if user.PIN != "" {
		dbUser.PIN = user.PIN
	}

	result := u.PostgresDB.Save(dbUser)
	if result.Error != nil {
		utils.LogErrors(result.Error)
		return nil, fmt.Errorf("error updating user")
	}
	cleanUser := dto.CleanedUser{
		UserId:    dbUser.UserId,
		Phone:     dbUser.Phone,
		ProfileId: dbUser.ProfileId,
		Email:     dbUser.Email,
		Profile:   dbUser.Profile,
		Role:      dbUser.Role,
	}
	return &cleanUser, nil
}
