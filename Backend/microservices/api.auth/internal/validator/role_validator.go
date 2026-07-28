package validator

import (
	"strings"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/api.auth/internal/repository"
	"lumini-hub/common/utils"
)

// RoleValidator valida regras de negócio relacionadas a perfis
type RoleValidator struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
	permRepo repository.PermissionRepository
}

// NewRoleValidator cria um novo validador de perfis
func NewRoleValidator(roleRepo repository.RoleRepository, userRepo repository.UserRepository, permRepo repository.PermissionRepository) *RoleValidator {
	return &RoleValidator{
		roleRepo: roleRepo,
		userRepo: userRepo,
		permRepo: permRepo,
	}
}

// isReservedRoleName impede que um requisitante que não seja DEVELOP
// crie/renomeie um Perfil pra "DEVELOP" — esse nome é o que concede bypass
// total no RBAC (checado por string, não por ID — ver utils.RoleDeveloper),
// então precisa ficar reservado independentemente de já existir uma linha
// com esse nome na tabela roles (ambientes sem o DEVELOP semeado manualmente
// ainda ficariam vulneráveis se essa checagem dependesse só da constraint
// unique).
func isReservedRoleName(name string, requesterRole string) bool {
	return requesterRole != utils.RoleDeveloper && strings.EqualFold(name, utils.RoleDeveloper)
}

// ValidateForCreation valida os dados para criação de um perfil
func (v *RoleValidator) ValidateForCreation(req domain.CreateRoleRequest, requesterRole string) error {
	var errors ValidationErrors

	if isReservedRoleName(req.Name, requesterRole) {
		errors.AddError("name", "nome de perfil inválido")
		return errors
	}

	// Verificar se o nome já existe
	exists, err := v.roleRepo.ExistsByName(req.Name)
	if err != nil {
		return err
	}
	if exists {
		errors.AddError("name", "nome de perfil já está em uso")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForUpdate valida os dados para atualização de um perfil
func (v *RoleValidator) ValidateForUpdate(id uint, req domain.UpdateRoleRequest, requesterRole string) error {
	var errors ValidationErrors

	// Verificar se o perfil existe
	role, err := v.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("id", "perfil não encontrado")
		return errors
	}

	// Verificar se o nome já está em uso por outro perfil
	if req.Name != "" && req.Name != role.Name {
		if isReservedRoleName(req.Name, requesterRole) {
			errors.AddError("name", "nome de perfil inválido")
			return errors
		}

		exists, err := v.roleRepo.ExistsByNameExcept(req.Name, id)
		if err != nil {
			return err
		}
		if exists {
			errors.AddError("name", "nome de perfil já está em uso")
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidateForDeletion valida se um perfil pode ser excluído.
//
// O Perfil "ADMIN" nunca pode ser excluído (é o perfil raiz do sistema —
// removê-lo tornaria a empresa inadministrável). Um Perfil só pode ser
// excluído se não tiver nenhuma permissão vinculada (força um passo
// explícito de "esvaziar as permissões antes de excluir", evitando perder
// de vista um Perfil com acesso configurado só porque ainda não foi
// atribuído a ninguém) e não estiver em uso por nenhum usuário.
func (v *RoleValidator) ValidateForDeletion(id uint) error {
	var errors ValidationErrors

	// Verificar se o perfil existe (com permissões, pra checagem abaixo)
	role, err := v.roleRepo.FindByIDWithPermissions(id)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("id", "perfil não encontrado")
		return errors
	}

	if strings.EqualFold(role.Name, utils.RoleAdmin) {
		errors.AddError("id", "o perfil ADMIN não pode ser excluído")
		return errors
	}

	if len(role.Permissions) > 0 {
		errors.AddError("id", "não é possível excluir um perfil que ainda tem permissões vinculadas")
	}

	// Verificar se o perfil está sendo usado por usuários
	count, err := v.userRepo.CountByRoleID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		errors.AddError("id", "não é possível excluir um perfil que está sendo usado por usuários")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

// ValidatePermissionUpdate valida a atualização de permissões de um perfil.
//
// Quem não é DEVELOP não pode atribuir a nenhum Perfil uma permissão do
// módulo utils.DeveloperModule (roles.*/permissions.*/admin.create_permissions)
// — do contrário um ADMIN poderia montar um Perfil customizado que concede,
// sem querer (ou de propósito), poderes de DEVELOP pra quem for atribuído a
// ele. A UI já esconde essas permissões do seletor (ver
// PermissionService.excludeModuleFor), isso aqui é a garantia server-side.
func (v *RoleValidator) ValidatePermissionUpdate(id uint, permissionIDs []uint, requesterRole string) error {
	var errors ValidationErrors

	// Verificar se o perfil existe
	role, err := v.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		errors.AddError("id", "perfil não encontrado")
		return errors
	}

	// Verificar se todas as permissões existem
	permissions, err := v.permRepo.FindByIDs(permissionIDs)
	if err != nil {
		return err
	}
	if len(permissions) != len(permissionIDs) {
		errors.AddError("permission_ids", "uma ou mais permissões não existem")
	}

	if requesterRole != utils.RoleDeveloper {
		for _, perm := range permissions {
			if perm.Module == utils.DeveloperModule {
				errors.AddError("permission_ids", "não é possível atribuir permissões do módulo Develop a um perfil")
				break
			}
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
