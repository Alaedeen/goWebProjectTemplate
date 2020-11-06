package repository

import (
	"crypto/sha1"
	"errors"
	"strings"

	models "github.com/Alaedeen/goWebProjectTemplate/models"
	"gorm.io/gorm"
)

// UserRepository ...
type UserRepository interface {
	GetUsers(role string, offset int, limit int) ([]models.User, error, int64)
	GetUsersByName(name string, role string, offset int, limit int) ([]models.User, error, int64)
	GetUser(id uint) (models.User, error)
	GetUserBy(keys []string, values []interface{}) (models.User, error)
	CreateUser(u models.User) (models.User, error)
	DeleteUser(id uint) error
	UpdateUser(m map[string]interface{}, id uint) error
}

// UserRepo ...
type UserRepo struct {
	Db *gorm.DB
}

// GetUsers ...
func (r *UserRepo) GetUsers(role string, offset int, limit int) ([]models.User, error, int64) {
	var Users []models.User
	var User models.User
	var count int64
	var err error
	if role == "user" {
		err = r.Db.Where("admin = ?", false).Offset(offset).Limit(limit).Find(&Users).Error
		r.Db.Model(&User).Where("admin = ?", false).Count(&count)
	} else if role == "admin" {
		err = r.Db.Where("admin = ? AND super_admin = ?", true, false).Offset(offset).Limit(limit).Find(&Users).Error
		r.Db.Model(&User).Where("admin = ? AND super_admin = ?", true, false).Count(&count)
	} else {
		err = r.Db.Where("admin = ? AND super_admin = ?", true, true).Offset(offset).Limit(limit).Find(&Users).Error
		r.Db.Model(&User).Where("admin = ? AND super_admin = ?", true, true).Count(&count)
	}

	return Users, err, count
}

// GetUsersByName ...
func (r *UserRepo) GetUsersByName(name string, role string, offset int, limit int) ([]models.User, error, int64) {
	var Users []models.User
	var User models.User
	var count int64
	var err error
	if role == "user" {
		err = r.Db.Where("admin = ? AND UPPER(name) LIKE ?", false, "%"+strings.ToUpper(name)+"%").Offset(offset).Limit(limit).Find(&Users).Error
		r.Db.Model(&User).Where("admin = ? AND UPPER(name) LIKE ?", false, "%"+strings.ToUpper(name)+"%").Count(&count)
	} else if role == "admin" {
		err = r.Db.Where("admin = ? AND super_admin = ? AND UPPER(name) LIKE ?", true, false, "%"+strings.ToUpper(name)+"%").Offset(offset).Limit(limit).Find(&Users).Error
		r.Db.Model(&User).Where("admin = ? AND super_admin = ? AND UPPER(name) LIKE ?", true, false, "%"+strings.ToUpper(name)+"%").Count(&count)
	} else {
		err = r.Db.Where("admin = ? AND super_admin = ? AND UPPER(name) LIKE ?", true, true, "%"+strings.ToUpper(name)+"%").Offset(offset).Limit(limit).Find(&Users).Error
		r.Db.Model(&User).Where("admin = ? AND super_admin = ? AND UPPER(name) LIKE ?", true, true, "%"+strings.ToUpper(name)+"%").Count(&count)
	}

	return Users, err, count
}

// GetUser ...
func (r *UserRepo) GetUser(id uint) (models.User, error) {
	var User models.User
	err := r.Db.First(&User, id).Error
	return User, err
}

// GetUserBy ...
func (r *UserRepo) GetUserBy(keys []string, values []interface{}) (models.User, error) {
	var User models.User
	var m map[string]interface{}
	var password string
	m = make(map[string]interface{})
	for index, value := range keys {
		if value == "password" {
			crypt := sha1.New()
			password = values[index].(string)
			crypt.Write([]byte(password))
			m[value] = crypt.Sum(nil)
		} else {
			m[value] = values[index]
		}

	}
	err := r.Db.Where(m).Find(&User).Error

	return User, err
}

// CreateUser ...
func (r *UserRepo) CreateUser(u models.User) (models.User, error) {
	User := u
	var user models.User
	//err := r.Db.Where(map[string]interface{}{"name": u.Name}).Find(&user).Error
	count := r.Db.Find(&user, "name = ?", u.Name).RowsAffected
	if count != 0 {
		return user, errors.New("ERROR: name already used")
	}
	//err = r.Db.Where(map[string]interface{}{"email": u.Email}).Find(&user).Error
	count = r.Db.Find(&user, "email = ?", u.Email).RowsAffected
	if count != 0 {
		return user, errors.New("ERROR: mail already used")
	}
	err := r.Db.Create(&User).Error
	return User, err
}

// DeleteUser ...
func (r *UserRepo) DeleteUser(id uint) error {
	user := models.User{}
	err := r.Db.First(&user, id).Error
	if err != nil {
		return err
	}

	user.ID = id
	err = r.Db.Delete(&user).Error
	return err

}

// UpdateUser ...
func (r *UserRepo) UpdateUser(m map[string]interface{}, id uint) error {
	user := models.User{}
	err := r.Db.Where("name = ? AND id != ?", m["name"], id).Find(&user).Error
	if err == nil {
		return errors.New("ERROR: name already used")
	}
	err = r.Db.Where("email = ? AND id != ?", m["email"], id).Find(&user).Error
	if err == nil {
		return errors.New("ERROR: mail already used")
	}
	err = r.Db.First(&user, id).Error
	if err != nil {
		return err
	}
	user.ID = id
	err1 := r.Db.Model(&user).Updates(m).Error
	return err1

}
