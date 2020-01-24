package seed

import (
	"log"
	"time"
)

var (
	userRoles = []map[string]interface{}{
		newUserRoleMap("00000000-0000-0000-0000-000000000001", "system-system-000000000001", "system-system", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000001"),

		newUserRoleMap("00000000-0000-0000-0000-000000000002", "superadmin-superadmin-000000000002", "superadmin-superadmin", "00000000-0000-0000-0000-000000000002", "00000000-0000-0000-0000-000000000002"),

		newUserRoleMap("00000000-0000-0000-0000-000000000003", "superadmin-000000000003", "admin-admin", "00000000-0000-0000-0000-000000000003", "00000000-0000-0000-0000-000000000003"),
	}
)

// CreateUsers seeding
func (s *step) CreateUserRoles() error {
	tx := s.GetTx()

	st := `INSERT INTO user_roles (id, slug, tenant_id, name,user_id, role_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :name, :user_id, :role_id, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at);`

	for _, ur := range userRoles {
		_, err := tx.NamedExec(st, ur)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func newUserRoleMap(id, slug, name, userID, roleID string) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,           //genUUID()
		"slug":          slug,         //genSlug(username),
		"tenant_id":     systemTenant, //genSlug(username),
		"name":          name,
		"user_id":       userID,
		"role_id":       roleID,
		"is_active":     true,
		"is_deleted":    false,
		"created_by_id": systemUUID,
		"updated_by_id": zeroUUID,
		"created_at":    time.Now(),
		"updated_at":    time.Time{},
	}
}
