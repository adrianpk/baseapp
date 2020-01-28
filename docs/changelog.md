# Changelog

## <a name="20200128"></a>20200128 - Authorization middleware.

  * Users have associated roles
  * Roles have associated permissions
  * Permissions can be retrieved by a tag.
  * Transitively users have all the permissions associated with each of their roles (tags).
  * Resources require permissions to be accessed (resource-permissions)
  * Resources are defined by path
  * Resources also are defined by a tag (but it has no direct influence in this use case)
  * Middleware then checks on each request that user has all resource required permissions (tags).
  * Not fully tested yet
  * Some mechanism for direct access to persistence and also a cache will be developed to optimize this strategy.

## <a name="20200127"></a>20200127 - RBAC administration.

  * Resource registration (paths)
  * Role registration
  * Permission registration
  * Attach and remove required permissions to resources.
  * Attach and remove permissions to roles.
  * Attach and remove roles to permissions.

See [routes](routes.md)
