package seed

import (
	"log"
	"time"
)

var (
	roles = []map[string]interface{}{
		newRoleMap("00000000-0000-0000-0000-000000000001", "system-000000000001", "system", "System role"),

		newRoleMap("00000000-0000-0000-0000-000000000002", "superadmin-000000000002", "superadmin", "Superadmin role"),

		newRoleMap("00000000-0000-0000-0000-000000000003", "admin-000000000003", "admin", "Admin role"),

		newRoleMap("00000000-0000-0000-0000-000000000004", "user-000000000004", "user", "User role"),
	}
)

// CreateUsers seeding
func (s *step) CreateRoles() error {
	tx := s.GetTx()

	st := `INSERT INTO roles (id, slug, tenant_id, name, description, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :name, :description, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at);`

	for _, r := range roles {
		_, err := tx.NamedExec(st, r)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func newRoleMap(id, slug, name, description string) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,           //genUUID()
		"slug":          slug,         //genSlug(username),
		"tenant_id":     systemTenant, //genSlug(username),
		"name":          name,
		"description":   description,
		"is_active":     true,
		"is_deleted":    false,
		"created_by_id": systemUUID,
		"updated_by_id": zeroUUID,
		"created_at":    time.Now(),
		"updated_at":    time.Time{},
	}
}
