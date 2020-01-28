# Gtd

## Inbox

  * Update authentication middleware
    * Read user permissions tags from cookie / bearer token.
    * Compare to those required to access a resource (path)
    * Allow / deny access.
    * Implement a cache to avoid rereading resource required permissions from persistence on each request.
    * Implement a cookie refresh mechanism to update user permissions without having to go through the sigin process.

  * Update SQL queries so all of them takes care of `is_active` and `is_deleted` columns.

## Next

* Implement an authentication middleware

    * Use encrypted cookie data / JWT claims to allow/disable path access.

* Implment a cache for all RBAC resources data

    * Avoid Db rountrips on each request
    * Nor a big problem right now (in memory volatile repo)
    * Memory?
    * Redis?

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


