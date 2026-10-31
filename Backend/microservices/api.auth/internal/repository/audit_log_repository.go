package repository

import (
	"strings"

	"lumini-hub/api.auth/internal/domain"
	"lumini-hub/common/database"
	"lumini-hub/common/utils"

	"gorm.io/gorm"
)

// AuditLogRepository lê a trilha de auditoria (system_logs). A escrita é do
// middlewares.LoggerMiddleware de cada serviço, não daqui.
type AuditLogRepository interface {
	Repository
	FindAllFiltered(pagination *utils.Pagination, filters domain.InGetAuditLogsFilters, hideUsersOfRole string) ([]database.SystemLog, error)
	FindEntityTypes() ([]string, error)
}

// GormAuditLogRepository implementa AuditLogRepository usando GORM
type GormAuditLogRepository struct {
	*BaseRepository
}

// NewAuditLogRepository cria um novo repository de auditoria
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &GormAuditLogRepository{BaseRepository: NewBaseRepository(db)}
}

// FindAllFiltered devolve uma página de registros filtrados, do mais novo para o mais antigo por
// padrão. hideUsersOfRole (quando não vazio) esconde o que foi feito por usuários desse Perfil,
// no SQL, para o total da paginação continuar certo (ver o Perfil DEVELOP no CLAUDE.md).
func (r *GormAuditLogRepository) FindAllFiltered(pagination *utils.Pagination, filters domain.InGetAuditLogsFilters, hideUsersOfRole string) ([]database.SystemLog, error) {
	var logs []database.SystemLog
	query := r.GetDB().Model(&database.SystemLog{})

	if search := strings.TrimSpace(filters.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("(action ILIKE ? OR details->>'username' ILIKE ?)", like, like)
	}
	if filters.EntityType != "" {
		query = query.Where("entity_type = ?", filters.EntityType)
	}
	if filters.Method != "" {
		query = query.Where("action LIKE ?", filters.Method+" %")
	}
	if filters.UserID != 0 {
		query = query.Where("user_id = ?", filters.UserID)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at < ?", *filters.DateTo)
	}
	if hideUsersOfRole != "" {
		query = query.Where(
			"(user_id IS NULL OR user_id NOT IN (SELECT id FROM users WHERE role_id IN (SELECT id FROM roles WHERE name = ?)))",
			hideUsersOfRole,
		)
	}

	paginated, err := utils.Paginate(&database.SystemLog{}, pagination, query)
	if err != nil {
		return nil, err
	}
	if err := paginated.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// FindEntityTypes devolve os tipos de entidade que já aparecem na auditoria (para o filtro da tela).
func (r *GormAuditLogRepository) FindEntityTypes() ([]string, error) {
	var types []string
	err := r.GetDB().Model(&database.SystemLog{}).
		Where("entity_type <> ''").
		Distinct().Order("entity_type").
		Pluck("entity_type", &types).Error
	return types, err
}
