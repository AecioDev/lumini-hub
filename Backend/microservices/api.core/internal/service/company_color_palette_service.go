package service

import (
	"lumini-hub/api.core/internal/domain"
	"lumini-hub/api.core/internal/repository"
	"lumini-hub/api.core/internal/validator"
)

// CompanyColorPaletteService gerencia operações de negócio relacionadas às
// paletas de cores personalizadas de uma empresa. Diferente de
// CompanyVisualConfigService, esta é uma lista de verdade — tem Delete, não
// tem Update (ver comentário em domain.CompanyColorPalette).
type CompanyColorPaletteService struct {
	uow       repository.UnitOfWork
	validator *validator.CompanyColorPaletteValidator
}

// NewCompanyColorPaletteService cria um novo serviço de paletas de cores personalizadas
func NewCompanyColorPaletteService(uow repository.UnitOfWork) *CompanyColorPaletteService {
	return &CompanyColorPaletteService{
		uow:       uow,
		validator: validator.NewCompanyColorPaletteValidator(uow.CompanyColorPalettes(), uow.Companies()),
	}
}

// GetByCompanyID lista as paletas personalizadas salvas de uma empresa
func (s *CompanyColorPaletteService) GetByCompanyID(companyID uint) ([]domain.ApiCompanyColorPalette, error) {
	palettes, err := s.uow.CompanyColorPalettes().FindAllByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	dtos := make([]domain.ApiCompanyColorPalette, 0, len(palettes))
	for _, palette := range palettes {
		dtos = append(dtos, domain.ApiCompanyColorPaletteFromModel(palette))
	}
	return dtos, nil
}

// CreateCompanyColorPalette salva uma nova paleta personalizada sob transação do Unit of Work
func (s *CompanyColorPaletteService) CreateCompanyColorPalette(req domain.CreateCompanyColorPaletteRequest) (*domain.ApiCompanyColorPalette, error) {
	if err := s.validator.ValidateForCreation(req); err != nil {
		return nil, err
	}

	palette := domain.CompanyColorPalette{
		CompanyID:      req.CompanyID,
		Name:           req.Name,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: nilIfEmpty(req.SecondaryColor),
		AccentColor:    nilIfEmpty(req.AccentColor),
	}

	err := s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyColorPalettes().Create(&palette)
	})
	if err != nil {
		return nil, err
	}

	dto := domain.ApiCompanyColorPaletteFromModel(palette)
	return &dto, nil
}

// DeleteCompanyColorPalette exclui uma paleta personalizada (soft delete)
func (s *CompanyColorPaletteService) DeleteCompanyColorPalette(id uint) error {
	if err := s.validator.ValidateForDeletion(id); err != nil {
		return err
	}

	return s.uow.Execute(func(uow repository.UnitOfWork) error {
		return uow.CompanyColorPalettes().Delete(id)
	})
}
