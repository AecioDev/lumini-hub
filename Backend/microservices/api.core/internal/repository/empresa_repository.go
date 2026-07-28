package repository

import (
	"errors"

	"lumini-hub/api.core/internal/domain"
	commonrepo "lumini-hub/common/repository"

	"gorm.io/gorm"
)

// EmpresaRepository define as operações de acesso a dados para empresas.
//
// Sem método de busca paginada/filtrada de propósito: o volume esperado de
// empresas por tenant é baixo (um grupo com sua Matriz e algumas filiais),
// não justifica o padrão POST /filter usado em entidades de alto volume
// (Customer/Supplier) — ver lumini_hub_backend_architecture skill.
type EmpresaRepository interface {
	commonrepo.Repository[domain.Empresa]
	FindByID(id uint) (*domain.Empresa, error) // Sobrescreve para incluir preloads
	FindByCNPJ(cnpj string) (*domain.Empresa, error)
	ExistsByCNPJ(cnpj string) (bool, error)
	ExistsByCNPJExcept(cnpj string, id uint) (bool, error)
	ExistsRoot() (bool, error)
	ExistsRootExcept(id uint) (bool, error)
	CountByParentID(parentID uint) (int64, error)
}

// GormEmpresaRepository implementa EmpresaRepository usando GORM e Generics
type GormEmpresaRepository struct {
	*commonrepo.GormRepository[domain.Empresa]
}

// NewEmpresaRepository cria um novo repository de empresas
func NewEmpresaRepository(db *gorm.DB) EmpresaRepository {
	return &GormEmpresaRepository{
		GormRepository: commonrepo.NewGormRepository[domain.Empresa](db),
	}
}

// FindByID busca uma empresa pelo ID pré-carregando as empresas vinculadas diretamente abaixo dela
func (r *GormEmpresaRepository) FindByID(id uint) (*domain.Empresa, error) {
	var empresa domain.Empresa
	err := r.GetDB().Preload("Children").First(&empresa, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &empresa, nil
}

// FindByCNPJ busca uma empresa pelo CNPJ
func (r *GormEmpresaRepository) FindByCNPJ(cnpj string) (*domain.Empresa, error) {
	var empresa domain.Empresa
	err := r.GetDB().Where("cnpj = ?", cnpj).First(&empresa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &empresa, nil
}

// ExistsByCNPJ verifica se existe uma empresa com o CNPJ especificado
func (r *GormEmpresaRepository) ExistsByCNPJ(cnpj string) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Empresa{}).Where("cnpj = ?", cnpj).Count(&count).Error
	return count > 0, err
}

// ExistsByCNPJExcept verifica se existe uma empresa com o CNPJ especificado, exceto a de ID informado
func (r *GormEmpresaRepository) ExistsByCNPJExcept(cnpj string, id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Empresa{}).Where("cnpj = ? AND id != ?", cnpj, id).Count(&count).Error
	return count > 0, err
}

// ExistsRoot verifica se já existe alguma empresa Matriz cadastrada (ParentID nulo)
func (r *GormEmpresaRepository) ExistsRoot() (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Empresa{}).Where("parent_id IS NULL").Count(&count).Error
	return count > 0, err
}

// ExistsRootExcept verifica se já existe outra empresa Matriz cadastrada, exceto a de ID informado
func (r *GormEmpresaRepository) ExistsRootExcept(id uint) (bool, error) {
	var count int64
	err := r.GetDB().Model(&domain.Empresa{}).Where("parent_id IS NULL AND id != ?", id).Count(&count).Error
	return count > 0, err
}

// CountByParentID conta quantas empresas estão vinculadas diretamente abaixo da empresa informada
func (r *GormEmpresaRepository) CountByParentID(parentID uint) (int64, error) {
	var count int64
	err := r.GetDB().Model(&domain.Empresa{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}
