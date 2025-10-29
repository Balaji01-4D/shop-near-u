# Shop-owner Frontend API & Models

This document describes the backend endpoints, request/response models, authentication, and recommended TypeScript data contracts and frontend patterns for a shop-owner-only web application (built against the `shop-near-u` backend in this repo).

Files used to generate this spec: `internal/server/routes.go`, `internal/shop/*`, `internal/product/*`, `internal/productCatlog/*`, `internal/middlewares/*`, `internal/models/*`, `internal/utils/*`.

## High-level contract
- Base URL: same origin as API server (e.g., `http://localhost:PORT`). Server CORS already allows `http://localhost:5173` and credentials.
- Auth: server sets an `Authorization` cookie (HttpOnly, SameSite=Lax) on successful register/login. All protected endpoints require that cookie and check role `shop_owner`.
- Response envelope: success responses use `{ success: boolean, message: string, data?: any }`. Error responses use `{ success: false, message: string, error?: string }`.

## Endpoints (shop-owner relevant)

1) Register shop
 - Method: POST
 - Path: /shop/register
 - Auth: none
 - Request JSON (ShopRegisterDTORequest):
   {
     "name": string,
     "owner_name": string,
     "type": string,
     "email": string,
     "mobile": string,
     "password": string,
     "address": string,
     "latitude": number,
     "longitude": number
   }
 - Success response: 201
   {
     "success": true,
     "message": "Shop registered successfully",
     "data": {
       "id": number,
       "name": string,
       "owner_name": string,
       "type": string,
       "email": string,
       "mobile": string,
       "address": string,
       "latitude": number,
       "longitude": number,
       "token": string
     }
   }
 - Side effect: server sets `Authorization` cookie (HttpOnly) using `Set-Cookie`.
 - Errors: 400 validation, 409 already exists, 500 server error.

2) Login (shop)
 - Method: POST
 - Path: /shop/login
 - Auth: none
 - Request JSON (ShopLoginDTORequest): { "email": string, "password": string }
 - Success: 200, same response shape as register (includes `token` field) and sets `Authorization` cookie.
 - Errors: 400 validation, 401 invalid credentials, 500 server error.

3) Get shop profile
 - Method: GET
 - Path: /shop/profile
 - Auth: RequireShopOwnerAuth (Authorization cookie with shop_owner role)
 - Response: 200
   { success: true, message: "Shop profile retrieved successfully", data: ShopRegisterDTOResponse }
 - Errors: 401 unauthorized, 500 server error.

4) Nearby shops (public)
 - Method: GET
 - Path: /shop
 - Query params: lat (required), lon (required), radius (optional, default 5000), limit (optional)
 - Response: 200, { success: true, message: "Nearby shops retrieved successfully", data: [NearByShopsDTORespone] }
 - Use-case: search for shops near a location (not strictly needed by owner UI but available).

5) Catalog products (search & create)
 - POST /api/catalog-products/  (Create catalog product)
   - Body: CreateCatalogProductDTO { name, brand, category, description, image_url }
   - Auth: none (controller currently exposes it publicly)
   - Response: 201, success (returns created product name in data)
 - GET /api/catalog-products/suggest?keyword=...&limit=...
   - Response: 200 { data: { products: CatalogProduct[] } }
   - Use: autocomplete/autosuggest when adding a product.

6) Shop products (owner-only)
 All endpoints below are under `/shop/products` and are protected with RequireShopOwnerAuth middleware.

 - POST /shop/products
   - Adds a new ShopProduct for the authenticated shop.
   - Body (AddProductDTORequest):
     { "catalog_id": number, "price": number, "stock": number, "discount": number, "is_available": boolean }
   - Response: 201 created, { success:true, message:"Product added successfully", data: null }
   - Errors: 400 validation, 401 unauthorized, 500 server error.

 - GET /shop/products
   - Returns list of products for the authenticated shop: 200 { data: ShopProduct[] }

 - GET /shop/products/:id
   - Params: id
   - Returns single ShopProduct
   - Errors: 400 invalid id, 404 not found

 - PUT /shop/products
   - Body (ProductUpdateDTORequest): { "id": number, "price"?, "stock"?, "discount"?, "is_available" }
   - Response: 200, product updated

 - DELETE /shop/products/:id
   - Deletes product
   - Response: 200 success, or 404 if not found

## Error and success payloads
- Success: { success: true, message: string, data?: any }
- Error: { success: false, message: string (usually "Request failed"), error?: string }

Frontend should always check the HTTP status code and the `success` boolean before trusting `data`.

## Authentication details and frontend implications
- The server sets an `Authorization` cookie via `Set-Cookie` with HttpOnly = true and SameSite=Lax. That means JavaScript cannot read the cookie directly.
- For API calls from the browser you must include credentials:
  - fetch: { credentials: 'include' }
  - axios: axiosInstance.defaults.withCredentials = true or request config { withCredentials: true }
- The register/login endpoints also return `token` inside the JSON response. You may use this token for client-side state (e.g., store in memory or use as confirmation), but rely on the cookie for actual auth with the server.
- Since token is HttpOnly on cookie, front-end guards should assume server-side validation via cookie. Use a `GET /shop/profile` call on app start to validate session and get the shop profile (or a dedicated /me endpoint if added).

## TypeScript interfaces (recommended)

// API envelope
interface ApiResponse<T> {
  success: boolean;
  message: string;
  data?: T;
}

interface ApiError {
  success: false;
  message: string; // "Request failed"
  error?: string;
}

// Shop
interface Shop {
  id: number;
  name: number;
  owner_name: string;
  type: string;
  email: string;
  mobile: string;
  address: string;
  latitude: number;
  longitude: number;
  created_at?: string;
}

interface ShopRegisterRequest {
  name: string;
  owner_name: string;
  type: string;
  email: string;
  mobile: string;
  password: string;
  address: string;
  latitude: number;
  longitude: number;
}

interface ShopRegisterResponse extends Shop {
  token: string;
}

// Catalog product
interface CatalogProduct {
  id: number;
  name: string;
  brand?: string;
  category: string;
  description?: string;
  image_url?: string;
  created_at?: string;
}

// Shop product (join of catalog + shop-specific fields)
interface ShopProduct {
  id: number;
  shop_id: number;
  catalog_id: number;
  price: number;
  stock: number;
  is_available: boolean;
  discount: number;
  created_at?: string;
  updated_at?: string;
  catalog_product?: CatalogProduct; // populated by GORM if preloaded
}

// DTOs
interface AddProductDTO {
  catalog_id: number;
  price: number;
  stock: number;
  discount?: number;
  is_available?: boolean;
}

interface UpdateProductDTO {
  id: number;
  price?: number;
  stock?: number;
  discount?: number;
  is_available?: boolean;
}

## Frontend pages and components (shop-owner-only)
- Auth pages
  - Login page (POST /shop/login)
  - Register page (POST /shop/register)

- Protected area (requires successful /shop/profile)
  - Dashboard: summary (total products, stock low alerts)
  - Products list: table with server-synced products (GET /shop/products)
    - Actions: Edit (open modal), Delete (confirm) -> optimistic UI for delete
  - Product form: create (POST /shop/products) & update (PUT /shop/products)
    - Add product flow should include a Catalog product autosuggest (GET /api/catalog-products/suggest)
  - Shop profile: GET /shop/profile and edit (not currently provided by API; if needed add endpoint)

UI/UX patterns and guidance
- Always call GET /shop/profile on app load to confirm session and populate global shop state.
- Use `credentials: 'include'` on all fetch/axios requests.
- Use server success envelope to show friendly messages; show `error` details when available.
- Optimistic updates: on create/update/delete, update UI optimistically but roll back on 4xx/5xx.
- Pagination & caching: server returns full list; if product lists grow, implement server-side pagination (not currently supported) or incremental queries.

## Edge cases & security
- Cookie is HttpOnly: you cannot read the token from JS; rely on server cookie for authentication.
- Token expiration: token set to expire in 7 days by server. Server returns 401 on expired token — redirect to login.
- CSRF: cookie-based auth + SameSite=Lax mitigates some CSRF risk, but if you allow cross-site POSTs (e.g., third-party forms), consider CSRF tokens.
- Validation errors: server returns 400 with error message string. Display validation message inline if possible.

## Sample fetch snippets (browser)

// Login
fetch('/shop/login', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
}).then(r => r.json());

// Get products
fetch('/shop/products', { credentials: 'include' }).then(r => r.json());

// Add product
fetch('/shop/products', {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ catalog_id, price, stock, discount, is_available })
}).then(r => r.json());

// Suggest catalog products
fetch('/api/catalog-products/suggest?keyword=' + encodeURIComponent(q) + '&limit=10')
  .then(r => r.json());

## Next steps (recommended priorities)
1. Implement frontend auth wrapper that calls `GET /shop/profile` on app start and routes to login if 401.
2. Build products list + add/edit/delete UI using the DTOs above.
3. Implement catalog product autosuggest when adding products.
4. Add client-side validation matching server rules (required fields, numeric ranges).
5. Consider adding server endpoints for updating shop profile and server-side pagination for products.

---
This file is saved in `docs/shop-owner-frontend.md`. Ask me to generate a TypeScript client module (API helper + typed hooks for React/React Query) and I will add it under `frontend/` with examples wired to `credentials: 'include'` and a small test harness.
