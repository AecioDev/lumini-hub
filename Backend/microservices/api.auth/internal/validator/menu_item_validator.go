package validator

import (
	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
)

// MenuItemValidator valida regras de negócio relacionadas a itens de menu
type MenuItemValidator struct {
	menuItemRepo repository.MenuItemRepository
	permRepo     repository.PermissionRepository
}

// NewMenuItemValidator cria um novo validador de itens de menu
func NewMenuItemValidator(menuItemRepo repository.MenuItemRepository, permRepo repository.PermissionRepository) *MenuItemValidator {
	return &MenuItemValidator{
		menuItemRepo: menuItemRepo,
		permRepo:     permRepo,
	}
}

func (v *MenuItemValidator) validateParentAndPermission(errors *ValidationErrors, parentID *uint, permissionID *uint, selfID *uint) {
	if parentID != nil {
		if selfID != nil && *parentID == *selfID {
			errors.AddError("parent_id", "um item de menu não pode ser pai de si mesmo")
		} else {
			parent, err := v.menuItemRepo.FindByID(*parentID)
			if err != nil {
				errors.AddError("parent_id", "erro ao validar o item pai")
			} else if parent == nil {
				errors.AddError("parent_id", "item de menu pai não encontrado")
			} else if selfID != nil && v.isDescendant(*parentID, *selfID) {
				errors.AddError("parent_id", "não é possível mover um item de menu para dentro de um de seus próprios descendentes")
			}
		}
	}

	if permissionID != nil {
		perm, err := v.permRepo.FindByID(*permissionID)
		if err != nil {
			errors.AddError("permission_id", "erro ao validar a permissão")
		} else if perm == nil {
			errors.AddError("permission_id", "permissão não encontrada")
		}
	}
}

// isDescendant verifica, subindo pela cadeia de pais a partir de candidateParentID,
// se ancestorID aparece no caminho — usado para impedir ciclos ao reatribuir o pai.
func (v *MenuItemValidator) isDescendant(candidateParentID uint, ancestorID uint) bool {
	currentID := candidateParentID
	for depth := 0; depth < 50; depth++ {
		if currentID == ancestorID {
			return true
		}
		current, err := v.menuItemRepo.FindByID(currentID)
		if err != nil || current == nil || current.ParentID == nil {
			return false
		}
		currentID = *current.ParentID
	}
	return false
}

// ValidateForCreation valida os dados para criação de um item de menu
func (v *MenuItemValidator) ValidateForCreation(req domain.CreateMenuItemRequest) error {
	var errors ValidationErrors
	v.validateParentAndPermission(&errors, req.ParentID, req.PermissionID, nil)

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de um item de menu
func (v *MenuItemValidator) ValidateForUpdate(id uint, req domain.UpdateMenuItemRequest) error {
	var errors ValidationErrors

	existing, err := v.menuItemRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		errors.AddError("id", "item de menu não encontrado")
		return errors
	}

	v.validateParentAndPermission(&errors, req.ParentID, req.PermissionID, &id)

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se um item de menu pode ser excluído
func (v *MenuItemValidator) ValidateForDeletion(id uint) error {
	var errors ValidationErrors

	existing, err := v.menuItemRepo.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		errors.AddError("id", "item de menu não encontrado")
		return errors
	}

	hasChildren, err := v.menuItemRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		errors.AddError("id", "não é possível excluir um item de menu que possui filhos — exclua ou mova os filhos primeiro")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
