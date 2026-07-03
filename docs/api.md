# WindNav API

统一响应：

```json
{
  "data": {},
  "error": null,
  "meta": {}
}
```

## Public

- GET /api/public/summary
- GET /api/public/categories
- GET /api/public/sites?page=1&page_size=20&q=&category=&tag=
- POST /api/public/sites/:id/click

## Auth

- POST /api/auth/login
- GET /api/auth/me
- POST /api/auth/logout

后台接口需要 Authorization: Bearer token。

## Admin

- GET /api/admin/categories
- POST /api/admin/categories
- PUT /api/admin/categories/:id
- DELETE /api/admin/categories/:id
- GET /api/admin/sites
- POST /api/admin/sites
- PUT /api/admin/sites/:id
- DELETE /api/admin/sites/:id
- GET /api/admin/tags
- POST /api/admin/tags
- PUT /api/admin/tags/:id
- DELETE /api/admin/tags/:id
- GET /api/admin/settings
- PUT /api/admin/settings
