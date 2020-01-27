# Routes

### Auth

| Method | Path          | Handler    |
|--------|---------------|------------|
| GET    | /auth/signup  | InitSignUp |
| POST   | /auth/signup  | SignUp     |
| GET    | /auth/signin  | InitSignIn |
| POST   | /auth/signin  | SignIn     |
| GET    | /auth/signout | SignOut    |

### User

| Method | Path                          | Handler    |
|--------|-------------------------------|------------|
| GET    | /users                        | Index      |
| GET    | /users/new                    | New        |
| POST   | /users                        | Create     |
| GET    | /users/{slug}                 | Show       |
| GET    | /users/{slug}/edit            | Edit       |
| PUT    | /users/{slug}                 | Update     |
| PATCH  | /users/{slug}                 | Update     |
| POST   | /users/{slug}/init-delete     | InitDelete |
| DELETE | /users/{slug}                 | Delete     |
| GET    | /users/{slug}/{token}/confirm | Confirm    |


### Account

| Method | Path                                    | Handler           |
|--------|-----------------------------------------|-------------------|
| GET    | /accounts/{slug}/roles                  | IndexAccountRoles |
| POST   | /accounts/{slug}/roles                  | AppendAccountRole |
| DELETE | /accounts/{slug}/roles/{subslug}        | RemoveAccountRole |


### Resource

| Method | Path                                    | Handler                  |
|--------|-----------------------------------------|--------------------------|
| GET    | /resources                              | Index                    |
| GET    | /resources/new                          | New                      |
| POST   | /resources                              | Create                   |
| GET    | /resources/{slug}                       | Show                     |
| GET    | /resources/{slug}/edit                  | Edit                     |
| PUT    | /resources/{slug}                       | Update                   |
| PATCH  | /resources/{slug}                       | Update                   |
| POST   | /resources/{slug}/init-delete           | InitDelete               |
| DELETE | /resources/{slug}                       | Delete                   |
| GET    | /resources/{slug}/permissions           | IndexResourcePermissions |
| POST   | /resources/{slug}/permissions           | AppendResourcePermission |
| DELETE | /resources/{slug}/permissions/{subslug} | RemoveResourcePermission |


### Role

| Method | Path                              | Handler              |
|--------|-----------------------------------|----------------------|
| GET    | /role                             | Index                |
| GET    | /role/new                         | New                  |
| POST   | /role                             | Create               |
| GET    | /role/{slug}                      | Show                 |
| GET    | /role/{slug}/edit                 | Edit                 |
| PUT    | /role/{slug}                      | Update               |
| PATCH  | /role/{slug}                      | Update               |
| POST   | /role/{slug}/init-delete          | InitDelete           |
| DELETE | /role/{slug}                      | Delete               |
| GET    | /role/{slug}/roles                | IndexRolePermissions |
| POST   | /role/{slug}/roles                | AppendRolePermission |
| DELETE | /role/{slug}/roles/{subslug}      | RemoveRolePermission |


### Permission

| Method | Path                            | Handler    |
|--------|---------------------------------|------------|
| GET    | /permissions                    | Index      |
| GET    | /permissions/new                | New        |
| POST   | /permissions                    | Create     |
| GET    | /permissions/{slug}             | Show       |
| GET    | /permissions/{slug}/edit        | Edit       |
| PUT    | /permissions/{slug}             | Update     |
| PATCH  | /permissions/{slug}             | Update     |
| POST   | /permissions/{slug}/init-delete | InitDelete |
| DELETE | /permissions/{slug}             | Delete     |
