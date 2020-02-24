# Gtd

## Inbox

  * Implement JSON endpoints and gRPC service.

  * Tests
    * I don't expect to be able to test endpoint handlers soon but service logic.

  * Update authentication middleware
    * Implement a cache to avoid rereading resource required permissions from persistence on each request.
      * Now configure a scheduller to periodically refresh it.
      * Or, alternatively or simultaneously, a hook to force update if a change occur in a tuple resource-permission.
    * Implement a cookie refresh mechanism to update user permissions without having to wait next signin.

  * Update SQL queries so all of them takes care of `is_active` and `is_deleted` columns.
  * Update SQL queries so all of them takes care of `tenant_id` column.

## Next

* Implement Profile resource

* Implement KeyValue resource

* Implement ResourceProperty resource

  * Analyze if needed
  * See if KeyValue cannot be used instead

* Implement Image resource

* Implement Service integration tests

  * After API is completely defined.

* Implement file upload for image resource

## Someday
  * Implement JSON REST endpoints
  * Implement gRPC endpoints
  * Improve table format for index lists.
  * Update translation texts under `assets/web/embed/i18n`
  * Homogenize migrator and seeder log format with the one used by application.

## Maybe


## Done

* Update authentication middleware
  * Read user permissions tags from cookie.
  * Compare to those required to access a resource (path)
  * Allow / deny access.
  * Implement a cache to avoid rereading resource required permissions from persistence on each request.

* Implement a cache for all RBAC resources data
    * Avoid Db rountrips on each request
    * Nor a big problem right now (in memory volatile repo)
    * Memory?
    * Redis?
