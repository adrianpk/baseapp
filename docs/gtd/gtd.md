# Gtd

## Inbox

* Implement RBAC resources

  * Resource
  * Role
  * Permission

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


