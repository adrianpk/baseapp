# Draft

## Notes

### User

- Geolocation
  - Used mainly to detect suspicious logins from unusual locations.

## Long names on receiver method names
  * I would like to avoid them but in this case I prefer to avoid ambiguity or have to rename them with longer names later to avoid collisions. If a better approach becomes apparent we will change the strategy later.

    * i.e.: GetResourcePermissionTagsByPath



### SQL

**Drop tables**

```sql
drop table role_permissions; drop table account_roles; drop table accounts; drop table migrations; drop table permissions; drop table resources; drop table roles; drop table seeds; drop table users;
```
