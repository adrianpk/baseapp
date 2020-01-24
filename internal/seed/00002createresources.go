package seed

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
	"log"
	"strings"
	"time"
)

var (
	resources = []map[string]interface{}{
		newResourceMap("00000000-0000-0000-0000-000000000001", "users-000000000001", "user", "", "/users"),

		newResourceMap("00000000-0000-0000-0000-000000000002", "resources-000000000002", "resources", "", "/resources"),

		newResourceMap("00000000-0000-0000-0000-000000000003", "roles-000000000003", "roles", "", "/roles"),

		newResourceMap("00000000-0000-0000-0000-000000000004", "permissions-000000000004", "permissions", "", "/permissions"),
	}
)

// CreateUsers seeding
func (s *step) CreateResources() error {
	tx := s.GetTx()

	st := `INSERT INTO resources (id, slug, tenant_id, name, description, tag, path, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :name, :description, :tag, :path, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at);`

	for _, r := range resources {
		_, err := tx.NamedExec(st, r)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func newResourceMap(id, slug, name, description, path string) map[string]interface{} {
	return map[string]interface{}{
		"id":            id,           //genUUID()
		"slug":          slug,         //genSlug(username),
		"tenant_id":     systemTenant, //genSlug(username),
		"name":          name,
		"description":   description,
		"tag":           strings.ToUpper(kbs.GenTag()),
		"path":          path,
		"is_active":     true,
		"is_deleted":    false,
		"created_by_id": systemUUID,
		"updated_by_id": zeroUUID,
		"created_at":    time.Now(),
		"updated_at":    time.Time{},
	}
}
