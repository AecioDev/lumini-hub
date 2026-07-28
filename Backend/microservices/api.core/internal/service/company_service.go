package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
	"lumini-hub/common/utils"
)

// CompanyService gerencia operações de negócio relacionadas a empresas
type CompanyService struct {
	uow       repository.UnitOfWork
	validator *validator.CompanyValidator
}

// NewCompanyService cria um novo serviço de empresas recebendo o Unit of Work
func NewCompanyService(uow repository.UnitOfWork) *CompanyService {
	return &CompanyService{
		uow:       uow,
		validator: validator.NewCompanyValidator(uow.Companies()),
	}
}

// GetCompanies retorna todas as empresas cadastradas (sem paginação — ver nota em CompanyRepository)
func (s *CompanyService) GetCompanies() (*domain.ApiCompanyList, error) {
	companies, err := s.uow.Companies().FindAll()
	if err != nil {
		return nil, err
	}

	companyDTOs := make([]domain.ApiCompany, 0, len(companies))
	for _, company := range companies {
		companyDTOs = append(companyDTOs, domain.ApiCompanyFromModel(company))
	}

	return &domain.ApiCompanyList{Companies: companyDTOs}, nil
}

// GetCompanyByID busca uma empresa pelo ID
func (s *CompanyService) GetCompanyByID(id uint) (*domain.ApiCompanyDetail, error) {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, utils.ErrNotFound
	}

	dto := domain.ApiCompanyDetailFromModel(*company)
	return &dto, nil
}

// CreateCompany cria uma nova empresa sob transação atômica do Unit of Work
func (s *CompanyService) CreateCompany(req domain.CreateCompanyRequest, userID uint) (*domain.ApiCompany, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	company := domain.Company{
		ParentID:    req.ParentID,
		LegalName:   req.LegalName,
		TradeName:   req.TradeName,
		TaxID:       utils.RemoveMask(req.TaxID),
		IsActive:    true,
		CreatedByID: &userID,
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Create(&company)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFromModel(company)
	return &dto, nil
}

// UpdateCompany atualiza uma empresa existente sob transação do Unit of Work
func (s *CompanyService) UpdateCompany(id uint, req domain.UpdateCompanyRequest, userID uint) (*domain.ApiCompany, error) {
	if err := s.validator.ValidateForUpdate(id, req); err != nil {
		return nil, err
	}

	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, utils.ErrNotFound
	}

	company.ParentID = req.ParentID
	company.LegalName = req.LegalName
	company.TradeName = req.TradeName
	company.TaxID = utils.RemoveMask(req.TaxID)
	company.IsActive = req.IsActive
	company.UpdatedByID = &userID

	err = s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Update(company)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyFromModel(*company)
	return &dto, nil
}

// DeleteCompany exclui uma empresa sob transação do Unit of Work
func (s *CompanyService) DeleteCompany(id uint) error {
	company, err := s.uow.Companies().FindByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return utils.ErrNotFound
	}

	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.Companies().Delete(id)
	})
}
