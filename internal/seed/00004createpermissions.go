package seed

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
	"log"
	"strings"
	"time"
)

var (
	permissions = []map[string]interface{}{
		newPermissionMap("00000000-0000-0000-0000-000000000001", "admin-system-task-000000000001", "admin-system-task", "Admin system tasks"),

		newPermissionMap("00000000-0000-0000-0000-000000000002", "create-user-000000000002", "create-user", "Create user"),

		newPermissionMap("00000000-0000-0000-0000-000000000003", "update-user-000000000003", "update-user", "Update user"),

		newPermissionMap("00000000-0000-0000-0000-000000000004", "delete-user-000000000004", "delete-user", "Delete user"),

		newPermissionMap("00000000-0000-0000-0000-000000000005", "admin-rbac-000000000005", "admin-rbac", "Admin RBAC"),
	}
)

// CreateUsers seeding
func (s *step) CreatePermissions() error {
	tx := s.GetTx()

	st := `INSERT INTO permissions (id, slug, tenant_id, name, description, tag, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :name, :description, :tag, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at);`

	for _, r := range permissions {
		_, err := tx.NamedExec(st, r)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func newPermissionMap(id, slug, name, description string) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,           //genUUID()
		"slug":          slug,         //genSlug(username),
		"tenant_id":     systemTenant, //genSlug(username),
		"name":          name,
		"description":   description,
		"tag":           strings.ToUpper(kbs.GenTag()),
		"is_active":     true,
		"is_deleted":    false,
		"created_by_id": systemUUID,
		"updated_by_id": zeroUUID,
		"created_at":    time.Now(),
		"updated_at":    time.Time{},
	}
}
