# Guía de Integración API - Gestrym Auth

## 🔑 Roles Disponibles

| `role_id` | Rol | Descripción |
|:---:|:---|:---|
| `1` | **Cliente** | Usuario final que entrena |
| `2` | **Entrenador (Coach)** | Gestiona clientes |
| `3` | **Gimnasio (Gym)** | Gestiona entrenadores y clientes |
| `4` | **Admin** | Acceso total al sistema |

---

## 📋 Base URL

```
http://tu-servidor/gestrym-auth
```

---

## 1️⃣ REGISTRO — `POST /public/auth/register`

### Body base (todos los casos)

```json
{
  "email": "usuario@ejemplo.com",
  "name": "Juan Pérez",
  "password": "mi_password123",
  "prefix": "+57",
  "phone": "3001234567",
  "role_id": 1,
  "registration_source": "self"
}
```

### Campos disponibles

| Campo | Tipo | Requerido | Descripción |
| :--- | :--- | :--- | :--- |
| `email` | string | ✅ Sí | Email del usuario (único por rol). |
| `name` | string | ✅ Sí | Nombre completo. |
| `password` | string | ✅ Sí | Contraseña del usuario. |
| `prefix` | string | ✅ Sí | Prefijo telefónico (ej: `+57`). |
| `phone` | string | ✅ Sí | Número de teléfono. |
| `role_id` | uint | ✅ Sí | `1` Cliente, `2` Entrenador, `3` Gimnasio. |
| `registration_source` | string | ✅ Sí | `"self"`, `"gym"`, o `"trainer"`. |
| `source_id` | uint | ⚠️ Condicional | ID del Gimnasio o Entrenador que registra (obligatorio si no es `"self"`). |
| `city` | string | ⚠️ Solo Gym | Ciudad del gimnasio. |
| `department` | string | ⚠️ Solo Gym | Departamento del gimnasio. |
| `country` | string | ⚠️ Solo Gym | País del gimnasio. |
| `workstation` | string | No | Dirección o sede del gimnasio. |
| `primary_color` | string | No | Color primario de la marca (`#HEX`). |
| `secondary_color` | string | No | Color secundario de la marca (`#HEX`). |
| `referral_code` | string | No | Código de referido. |
| `avatar_file_id` | uint | No | ID del archivo de imagen del perfil. |

---

### Casos de uso según rol y origen

#### 🧍 Cliente se registra solo
```json
{
  "email": "cliente@ejemplo.com",
  "name": "Ana Gómez",
  "password": "pass123",
  "prefix": "+57",
  "phone": "3109876543",
  "role_id": 1,
  "registration_source": "self"
}
```

---

#### 🏋️ Entrenador se registra solo
```json
{
  "email": "coach@ejemplo.com",
  "name": "Carlos Ruiz",
  "password": "pass123",
  "prefix": "+57",
  "phone": "3201112222",
  "role_id": 2,
  "registration_source": "self",
}
```

---

#### 🏢 Gimnasio se registra solo
> ⚠️ `city`, `department` y `country` son **obligatorios** para gimnasios.

```json
{
  "email": "gym@ejemplo.com",
  "name": "GymFit Centro",
  "password": "pass123",
  "prefix": "+57",
  "phone": "3453334444",
  "role_id": 3,
  "registration_source": "self",
  "city": "Bogotá",
  "department": "Cundinamarca",
  "country": "Colombia",
}
```

---

#### 🏢➡️🏋️ Gimnasio registra un Entrenador
> `source_id` = ID del usuario Gimnasio.

```json
{
  "email": "coach2@ejemplo.com",
  "name": "Luisa Torres",
  "password": "pass123",
  "prefix": "+57",
  "phone": "3115556666",
  "role_id": 2,
  "registration_source": "gym",
  "source_id": 5
}
```

---

#### 🏢➡️🧍 Gimnasio registra un Cliente
> `source_id` = ID del usuario Gimnasio.

```json
{
  "email": "cliente2@ejemplo.com",
  "name": "Pedro Díaz",
  "password": "pass123",
  "prefix": "+57",
  "phone": "3127778888",
  "role_id": 1,
  "registration_source": "gym",
  "source_id": 5
}
```

---

#### 🏋️➡️🧍 Entrenador registra un Cliente
> `source_id` = ID del usuario Entrenador.

```json
{
  "email": "cliente3@ejemplo.com",
  "name": "María López",
  "password": "pass123",
  "prefix": "+57",
  "phone": "3139990000",
  "role_id": 1,
  "registration_source": "trainer",
  "source_id": 12
}
```

---

### ✅ Respuesta exitosa del Registro (201 Created)
```json
{
  "id": 42,
  "email": "usuario@ejemplo.com",
  "name": "Juan Pérez",
  "phone": "3001234567",
  "role_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI..."
}
```
> ⚠️ El `token` devuelto en el registro es el **token de activación de email**, **NO** el `access_token` para hacer login.
> El usuario **no puede iniciar sesión** hasta haber confirmado su email.

### ❌ Errores posibles del Registro

| Código | Mensaje | Causa |
| :---: | :--- | :--- |
| `400` | `"Key: 'RegisterRequest.Email' Error:Field validation..."` | Campo obligatorio faltante o mal formateado |
| `400` | `"source_id es requerido para registrar un entrenador desde un gimnasio"` | Falta `source_id` en source `"gym"` o `"trainer"` |
| `400` | `"city, department y country son requeridos para registrar un gimnasio"` | Faltan campos de ubicación al registrar un Gym |
| `409` | `"ya existe un usuario con ese email"` | Email ya activo con el mismo rol |
| `409` | `"ya existe un usuario activo con ese email y rol diferente"` | El email ya existe pero en otro rol |
| `500` | `"error creando nuevo usuario"` | Fallo en la base de datos |
| `500` | `"tipo de token de activación no encontrado"` | El catálogo de tipos de token no está configurado en BD |

---

## 2️⃣ LOGIN — `POST /public/login`

```json
{
  "email": "usuario@ejemplo.com",
  "password": "mi_password123"
}
```

### ✅ Respuesta exitosa (200 OK)
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
  "role_id": 1,
  "email": "usuario@ejemplo.com"
}
```
> 💾 Guarda `access_token` y `role_id` en el estado global o LocalStorage para manejar la navegación por rol.

### ❌ Errores posibles del Login

| Código | Mensaje | Causa |
| :---: | :--- | :--- |
| `400` | `"Key: 'LoginRequest.Email' Error:..."` | Falta email o password |
| `401` | `"credenciales inválidas"` | Password incorrecto o usuario no existe |
| `401` | `"usuario inactivo, confirma tu email"` | El usuario no ha confirmado el email aún |

---

## 3️⃣ CONFIRMAR EMAIL — `GET /public/auth/confirm?token=<TOKEN>`

Este endpoint se llama automáticamente cuando el usuario hace clic en el enlace del correo que recibe al registrarse.

**Ejemplo de URL:**
```
GET /gestrym-auth/public/auth/confirm?token=eyJhbGciOiJIUzI1NiIsInR5cCI...
```

### ✅ Respuesta exitosa (200 OK)
```json
{
  "message": "Email confirmado exitosamente. Tu cuenta ha sido activada.",
  "user": {
    "id": 42,
    "email": "usuario@ejemplo.com",
    "name": "Juan Pérez",
    "phone": "3001234567",
    "prefix": "+57",
    "role_id": 1,
    "role_name": "cliente",
    "is_active": true,
    "email_confirmed": true
  }
}
```

### ❌ Errores posibles

| Código | Mensaje | Causa |
| :---: | :--- | :--- |
| `400` | `"token es requerido"` | No se envió el query param `token` |
| `500` | `"token de activación inválido o no registrado"` | Token ya usado, expirado o inválido |

---

## 4️⃣ RECUPERACIÓN DE CONTRASEÑA

### Paso 1 — Solicitar recovery `POST /public/auth/password/recovery`

```json
{
  "email": "usuario@ejemplo.com"
}
```

**✅ Respuesta (200 OK):**
```json
{ "message": "Email de recuperación enviado" }
```

**❌ Errores:**

| Código | Mensaje | Causa |
| :---: | :--- | :--- |
| `400` | `"Key: 'PasswordRecoveryRequest.Email'..."` | Falta el campo `email` |
| `500` | `"usuario no encontrado"` | El email no está registrado |
| `500` | `"error enviando email de recuperación"` | Fallo en el servicio de notificaciones |

---

### Paso 2 — Restablecer password `POST /public/auth/password/reset`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI...",
  "password": "nueva_password123"
}
```

**✅ Respuesta (200 OK):** Devuelve el objeto `user` actualizado igual que en Confirmar Email.

**❌ Errores:**

| Código | Mensaje | Causa |
| :---: | :--- | :--- |
| `400` | Validación fallida | Falta `token` o `password` |
| `500` | `"token de recuperación inválido o no registrado"` | Token ya usado, expirado o inválido |

---

## 5️⃣ VALIDAR TOKEN — `GET /public/auth/validate`

Útil para verificar si el `access_token` sigue siendo válido (por ejemplo al recargar la app).

**Por Header:**
```
Authorization: Bearer <access_token>
```
**O por Query Param:**
```
GET /public/auth/validate?token=<access_token>
```

### ✅ Respuesta exitosa (200 OK)
```json
{
  "valid": true,
  "user_id": 42,
  "role_id": 1,
  "access_level_id": 1,
  "email": "usuario@ejemplo.com",
  "expires_at": 1713999999
}
```

### ❌ Errores

| Código | Mensaje | Causa |
| :---: | :--- | :--- |
| `400` | `"token es requerido"` | No se envió token |
| `401` | `"invalid token"` | Token inválido o expirado |

---

## 6️⃣ ENDPOINTS PRIVADOS (requieren JWT)

> ⚠️ Todos los endpoints `/private/*` requieren el header:
> ```
> Authorization: Bearer <access_token>
> ```

### Listar usuarios — `GET /private/auth/users`

**Query params opcionales:**
- `page` (default: 1)
- `page_size` (default: 10)
- `name`, `email`, `dni`, `role_id`

**✅ Respuesta:**
```json
{
  "page": 1,
  "page_size": 10,
  "total": 50,
  "results": [
    {
      "id": 1,
      "name": "Juan Pérez",
      "email": "juan@ejemplo.com",
      "phone": "3001234567",
      "role_id": 1,
      "role_name": "cliente"
    }
  ]
}
```

---

### Obtener usuario por ID — `GET /private/auth/users/:id`

**✅ Respuesta:**
```json
{
  "id": 42,
  "email": "usuario@ejemplo.com",
  "name": "Juan Pérez",
  "phone": "3001234567",
  "prefix": "+57",
  "role_id": 1,
  "role_name": "cliente",
  "is_active": true,
  "email_confirmed": true
}
```

---

### Actualizar usuario — `PUT /private/auth/users/:id`

Todos los campos son **opcionales** (solo se actualiza lo que se envíe):
```json
{
  "name": "Juan Carlos Pérez",
  "email": "nuevo@ejemplo.com",
  "phone": "3009998888",
  "prefix": "+57",
  "password": "nueva_pass123"
}
```

---

### Eliminar usuario (soft delete) — `DELETE /private/auth/users/:id`

No requiere body. Marca el usuario con `is_active: false`.

**✅ Respuesta:** `204 No Content`

---

### Ver relaciones — `GET /private/auth/relationships`

El comportamiento varía según el rol del usuario autenticado:

- **Si es Entrenador (role_id: 2):** Devuelve sus clientes independientes y los clientes que tiene en gimnasios.
- **Si es Gimnasio (role_id: 3):** Devuelve sus entrenadores con sus respectivos clientes.

**✅ Respuesta para Entrenador:**
```json
{
  "independent_clients": [
    { "id": 10, "name": "Ana", "email": "ana@ejemplo.com", "phone": "..." }
  ],
  "gym_clients": [
    {
      "trainer_id": 5,
      "trainer_name": "Carlos",
      "trainer_email": "carlos@ejemplo.com",
      "clients": [...]
    }
  ]
}
```

**❌ Error si el rol no tiene permiso:**
```json
{ "error": "acceso no permitido para este rol" }
```

---

## ❌ Tabla General de Códigos de Error

| Código HTTP | Significado | Cuándo ocurre |
|:---:|:---|:---|
| `400 Bad Request` | Error de validación | Faltan campos obligatorios o formato incorrecto |
| `401 Unauthorized` | No autorizado | Token inválido/expirado o credenciales incorrectas |
| `403 Forbidden` | Prohibido | Rol sin permisos para esa acción |
| `404 Not Found` | No encontrado | Recurso o ruta no existe |
| `409 Conflict` | Conflicto de datos | Email ya registrado en ese rol |
| `500 Internal Server Error` | Error del servidor | Fallo en BD o servicios externos |

---

## 🔄 Flujo Completo de Registro y Acceso

```
1. POST /public/auth/register
   └── Usuario creado con is_active: false

2. [Backend envía email con enlace de activación]

3. Usuario hace clic en el enlace del email:
   GET /public/auth/confirm?token=...
   └── Cuenta activada (is_active: true, email_confirmed: true)

4. POST /public/login
   └── Obtiene access_token + refresh_token + role_id

5. Requests a rutas protegidas:
   Header: Authorization: Bearer <access_token>
   └── Acceso a /private/* según rol
```

---

## 💡 Tips de Implementación

1. **Guarda siempre** `access_token` y `role_id` tras el login (LocalStorage o estado global).
2. **Valida el token** al cargar la app con `GET /public/auth/validate` para desloguear automáticamente si expiró.
3. **Diferencia los formularios** de registro según el `role_id` seleccionado (mostrar `city/department/country` solo si es Gym).
4. **El `source_id`** siempre es el `id` del usuario creador (no del perfil), retornado en el campo `id` de la respuesta de registro o login.
5. **Manejo de 409**: Muestra un mensaje claro al usuario si el email ya está en uso.
