package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
)

// EmpresaService gerencia operações de negócio relacionadas a empresas
type EmpresaService struct {
	uow       repository.UnitOfWork
	validator *validator.EmpresaValidator
}

// NewEmpresaService cria um novo serviço de empresas recebendo o Unit of Work
func NewEmpresaService(uow repository.UnitOfWork) *EmpresaService {
	return &EmpresaService{
		uow:       uow,
		validator: validator.NewEmpresaValidator(uow.Empresas()),
	}
}

// GetEmpresas retorna todas as empresas cadastradas (sem paginação — ver nota em EmpresaRepository)
func (s *EmpresaService) GetEmpresas() (*domain.ApiEmpresaList, error) {
	empresas, err := s.uow.Empresas().FindAll()
	if err != nil {
		return nil, err
	}

	empresaDTOs := make([]domain.ApiEmpresa, 0, len(empresas))
	for _, empresa := range empresas {
		empresaDTOs = append(empresaDTOs, domain.ApiEmpresaFromModel(empresa))
	}

	return &domain.ApiEmpresaList{Empresas: empresaDTOs}, nil
}

// GetEmpresaByID busca uma empresa pelo ID
func (s *EmpresaService) GetEmpresaByID(id uint) (*domain.ApiEmpresaDetail, error) {
	empresa, err := s.uow.Empresas().FindByID(id)
	if err != nil {
		return nil, err
	}
	if empresa == nil {
		return nil, utils.ErrNotFound
	}

	dto := domain.ApiEmpresaDetailFromModel(*empresa)
	return &dto, nil
}

// CreateEmpresa cria uma nova empresa sob transação atômica do Unit of Work
func (s *EmpresaService) CreateEmpresa(req domain.CreateEmpresaRequest, userID uint) (*domain.ApiEmpresa, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	empresa := domain.Empresa{
		ParentID:     req.ParentID,
		RazaoSocial:  req.RazaoSocial,
		NomeFantasia: req.NomeFantasia,
		CNPJ:         utils.RemoveMask(req.CNPJ),
		IsActive:     true,
		CreatedByID:  &userID,
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Empresas().Create(&empresa)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiEmpresaFromModel(empresa)
	return &dto, nil
}

// UpdateEmpresa atualiza uma empresa existente sob transação do Unit of Work
func (s *EmpresaService) UpdateEmpresa(id uint, req domain.UpdateEmpresaRequest, userID uint) (*domain.ApiEmpresa, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	empresa, err := s.uow.Empresas().FindByID(id)
	if err != nil {
		return nil, err
	}
	if empresa == nil {
		return nil, utils.ErrNotFound
	}

	empresa.ParentID = req.ParentID
	empresa.RazaoSocial = req.RazaoSocial
	empresa.NomeFantasia = req.NomeFantasia
	empresa.CNPJ = utils.RemoveMask(req.CNPJ)
	empresa.IsActive = req.IsActive
	empresa.UpdatedByID = &userID

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Empresas().Update(empresa)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiEmpresaFromModel(*empresa)
	return &dto, nil
}

// DeleteEmpresa exclui uma empresa sob transação do Unit of Work
func (s *EmpresaService) DeleteEmpresa(id uint) error {
	empresa, err := s.uow.Empresas().FindByID(id)
	if err != nil {
		return err
	}
	if empresa == nil {
		return utils.ErrNotFound
	}

	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Empresas().Delete(id)
	})
}
