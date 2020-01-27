package seed

import (
	"log"
	"time"
)

var (
	resourcePermissions = []map[string]interface{}{
		newResourcePermissionMap("00000000-0000-0000-0000-000000000001", "system-admin-system-task-000000000001", "system-admin-system-task", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000001"),

		newResourcePermissionMap("00000000-0000-0000-0000-000000000002", "system-create-user-000000000002", "system-create-user", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"),

		newResourcePermissionMap("00000000-0000-0000-0000-000000000003", "system-update-user-000000000003", "system-update-user", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000003"),

		newResourcePermissionMap("00000000-0000-0000-0000-000000000004", "system-delete-user-000000000004", "system-delete-user", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000004"),

		newResourcePermissionMap("00000000-0000-0000-0000-000000000005", "system-admin-rbac-000000000005", "system-admin-rbac", "00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000005"),
	}
)

// CreateUsers seeding
func (s *step) CreateResourcePermissions() error {
	tx := s.GetTx()

	st := `INSERT INTO resource_permissions (id, slug, tenant_id, name, resource_id, permission_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :name, :resource_id, :permission_id, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at);`

	for _, ur := range resourcePermissions {
		_, err := tx.NamedExec(st, ur)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func newResourcePermissionMap(id, slug, name, resourceID, permissionID string) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,   //genUUID()
		"slug":          slug, //genSlug(name),
		"tenant_id":     systemTenant,
		"name":          name,
		"resource_id":   resourceID,
		"permission_id": permissionID,
		"is_active":     true,
		"is_deleted":    false,
		"created_by_id": systemUUID,
		"updated_by_id": zeroUUID,
		"created_at":    time.Now(),
		"updated_at":    time.Time{},
	}
}
