package repository

import (
	"errors"
	"simple-erp-service/internal/data-structure/models"
	"simple-erp-service/internal/utils"

	"gorm.io/gorm"
)

// UserRepository define as operações de acesso a dados para usuários
type UserRepository interface {
	Repository
	FindAll(pagination *models.Pagination) ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByIDWithRole(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsernameExcept(username string, id uint) (bool, error)
	ExistsByEmailExcept(email string, id uint) (bool, error)
	CountByRoleID(roleID uint) (int64, error)
}

// GormUserRepository implementa UserRepository usando GORM
type GormUserRepository struct {
	*BaseRepository
}

// NewUserRepository cria um novo repository de usuários
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindAll retorna todos os usuários com paginação
func (r *GormUserRepository) FindAll(pagination *models.Pagination) ([]models.User, error) {
	var users []models.User

	query := r.GetDB().Model(&models.User{}).Preload("Role")
	query, err := utils.Paginate(&models.User{}, pagination, query)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID busca um usuário pelo ID
func (r *GormUserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.GetDB().First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByIDWithRole busca um usuário pelo ID e carrega o relacionamento com Role
func (r *GormUserRepository) FindByIDWithRole(id uint) (*models.User, error) {
	var user models.User
	if err := r.GetDB().Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername busca um usuário pelo nome de usuário
func (r *GormUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail busca um usuário pelo email
func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create cria um novo usuário
func (r *GormUserRepository) Create(user *models.User) error {
	return r.GetDB().Create(user).Error
}

// Update atualiza um usuário existente
func (r *GormUserRepository) Update(user *models.User) error {
	return r.GetDB().Save(user).Error
}

// Delete exclui um usuário (soft delete)
func (r *GormUserRepository) Delete(id uint) error {
	return r.GetDB().Delete(&models.User{}, id).Error
}

// ExistsByUsername verifica se existe um usuário com o nome de usuário especificado
func (r *GormUserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail verifica se existe um usuário com o email especificado
func (r *GormUserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ExistsByUsernameExcept verifica se existe um usuário com o nome de usuário especificado, exceto o usuário com o ID especificado
func (r *GormUserRepository) ExistsByUsernameExcept(username string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.User{}).Where("username = ? AND id != ?", username, id).Count(&count).Error
	return count > 0, err
}

// ExistsByEmailExcept verifica se existe um usuário com o email especificado, exceto o usuário com o ID especificado
func (r *GormUserRepository) ExistsByEmailExcept(email string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&models.User{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}

// CountByRoleID conta quantos usuários estão usando um determinado perfil
func (r *GormUserRepository) CountByRoleID(roleID uint) (int64, error) {
	var count int64
	err := r.GetDB().Model(&models.User{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}
