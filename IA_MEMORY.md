# IA Memory

## 1. 📌 Project Overview

* Descripción del sistema
  - Microservicio de autenticación para una plataforma de entrenamiento deportivo y nutricional.
  - Gestiona usuarios, roles, permisos, autenticación JWT y tokens.
  - Integra lógica de registro de usuarios, login y listado de usuarios.

* Objetivo del negocio
  - Proveer un servicio de AUTH centralizado que permita gestionar identidades y accesos para clientes, entrenadores y gimnasios.
  - Soportar distintos métodos de registro y roles de usuario, manteniendo el backend desacoplado mediante arquitectura hexagonal.

## 2. 🏗️ Architecture

* Tipo de arquitectura
  - Arquitectura hexagonal / ports and adapters.

* Capas identificadas
  - `common`: configuración, middleware, rutas, utilidades y modelos compartidos.
  - `core/*/app`: capa de aplicación / servicios de dominio.
  - `core/*/domain`: definiciones de puertos e structs de request/response.
  - `core/*/infra`: adaptadores de infraestructura: controladores HTTP y repositorios GORM.

* Flujo de dependencias
  - Controllers → app services → domain ports → infra repositories
  - Middleware JWT y Role en `common/middleware` inyecta contexto en requests.
  - Rutas definen grupos `public`, `private`, `protected` con middleware específico.

## 3. 📦 Modules / Services

* Descripción del servicio AUTH
  - Servicio central de AUTH ubicado en `src/core/auth`.
  - Usa `authService` y `IAuthRepository` para operaciones de usuario.
  - Consume `core/jwt/app` para generación y registro de tokens.

* Responsabilidades actuales
  - Registro de usuarios (`RegisterUser`).
  - Generación de JWT de activación y persistencia de token de activación.
  - Listado de usuarios activos con filtros.
  - Validación de email y creación/actualización de usuario.
  - Exposición de endpoints HTTP en controladores públicos y privados.

## 4. 🧩 Domain Model

* Entidades principales
  - `User`: email, password, nombre completo, prefijo, teléfono, estado activo, role_id, confirmación de email.
  - `Role`: roles de usuario con nombre y descripción.
  - `Permission`: permisos por rol y recurso.
  - `RolePermission`: relación muchos-a-muchos entre roles y permisos.
  - `Action`: acciones que pueden registrarse para permisos.
  - `AccessLevel`: niveles de acceso.
  - `UserTokenType` / `UserToken`: tipos de token y tokens registrados.
  - `TrainerProfile`: asociación de entrenador/gimnasio.

* Relaciones clave
  - `User` tiene un `Role` mediante `RoleID`.
  - `Permission` tiene referencia a `Role`.
  - `RolePermission` conecta `Role` y `Permission`.
  - `UserToken` asociado a `User` y a `UserTokenType`.
  - `TrainerProfile` representa asociación de un entrenador con su gimnasio.

## 5. 🔐 Auth & Security

* Cómo funciona JWT actualmente
  - El token se genera a partir de `core/jwt/app` usando `GenerateJwtTokenRequest`.
  - Claims personalizados incluyen `user_id`, `role_id` y `access_level_id`.
  - `common/middleware/JWTModdleware.go` valida bearer token HMAC con `JWT_KEY`.
  - Si el token es válido, se pone en contexto `user_id`, `role_id`, `access_level_id`.

* Roles y permisos
  - Rol definidos en `common/middleware/RoleMiddleware.go`: Cliente=1, Coach=2, Gym=3, Admin=4.
  - Existen roles adicionales en constantes de utilidades como `ROLE_CLIENT`, `ROLE_COACH`, `ROLE_GYM`, `ROLE_COACH_GYM`.
  - `RequireRoles(4)` se usa para proteger catálogos administrativos.
  - El servicio tiene endpoints para gestión de `roles`, `permissions`, `actions`, `access_levels`, `token_types`.

## 6. 🔄 Current Flows

* Registro de usuario
  - Registro público disponible en `/public/auth/register` para self-signup.
  - Registro privado disponible en `/private/auth/register` para gym o entrenador.
  - Se recibe `RegisterRequest` con campos `registration_source` y `source_id` para diferenciar self/gym/trainer.
  - El servicio permite reactivar usuarios existentes con el mismo email y el mismo rol.
  - Se crean asociaciones de negocio: `TrainerProfile` para entrenadores de gym, `TrainerClient` para clientes de entrenador y `GymClient` para clientes registrados por un gym.
  - **El usuario se crea con `is_active: false` y `email_confirmed: false` hasta confirmar el email.**
  - **La contraseña es requerida en el registro.** El usuario se activa SOLO al confirmar el email.
  - Se genera JWT de activación y se registra como `UserToken` tipo activación.
  - Se envía solicitud HTTP al servicio de notificaciones (URL configurable por env var `NOTIFICATION_SERVICE_URL`).
  - El email enviado incluye el nombre real del usuario (`user.FullName`) en lugar de un placeholder.
  - Existe un endpoint de confirmación de email que recibe el token de activación y activa al usuario.
  - Se agregó soporte para obtener usuario por ID, modificar datos de usuario y soft delete de usuario (`is_active = false`).
  - Se agregó flujo de recuperación de contraseña: solicitud de recuperación por email y restablecimiento con token.

* Login
  - La lógica de login tradicional y login con Google está presente en el controlador público pero actualmente está comentada.
  - El diseño sugiere login por email/password y Google ID token con la generación de access/refresh tokens.

* Generación de token
  - Se genera un JWT desde `authService` con payload de usuario, rol y nivel de acceso.
  - `core/jwt/app` es responsable de la generación y registro de tokens.
  - El token de activación se almacena en `UserToken`.

## 7. ⚠️ Constraints (MUY IMPORTANTE)

* Qué NO se debe romper
  - No modificar el código existente directamente sin justificación.
  - No romper la arquitectura hexagonal ni los contratos de puertos/adaptadores.
  - No cambiar modelos existentes sin justificar la migración.
  - Mantener la separación entre capas `app`, `domain`, `infra`.

* Reglas técnicas del proyecto
  - Respetar el flujo controller → service → repository.
  - Usar middleware para autenticación y roles en lugar de lógica ad hoc en controladores.
  - No introducir dependencias cruzadas entre dominios.
  - Mantener `common` para utilidades, configuración, middleware y modelos compartidos.

## 8. 🚧 Technical Decisions

* Decisiones importantes ya tomadas
  - Uso de arquitectura hexagonal con puertos (`domain/ports`) y adaptadores (`infra`).
  - Registro de rutas en `common/routes/ServerRoutesDefinition.go` y separación en grupos `public`, `private`, `protected`.
  - JWT HMAC centralizado a través de `JWT_KEY` y middleware de validación.
  - Almacenamiento de tokens y tipos de token como entidad de dominio.
  - Email de confirmación mediante llamada HTTP a un servicio externo.
  - Registro de usuarios con origen (`self`, `gym`, `trainer`) y creación de relaciones de negocio en los repositorios.

* Librerías usadas
  - `gin-gonic/gin` para HTTP.
  - `gorm.io/gorm` para ORM/Postgres.
  - `github.com/spf13/viper` para configuración.
  - `github.com/golang-jwt/jwt/v4` para JWT.
  - `golang.org/x/crypto/bcrypt` para hashing de contraseñas.
  - `github.com/swaggo/gin-swagger` para documentación Swagger.

* Patrones aplicados
  - Repository pattern.
  - Service layer / application service.
  - Hexagonal architecture / ports and adapters.
  - Middleware para autenticación y autorización.
  - DTOs de request/response en `domain/structs`.

## 9. 📈 Possible Improvements (SIN IMPLEMENTAR)

* Hacer que el endpoint de registro sea realmente público o separar `self-signup` de registro interno.
* Reactivar e implementar los métodos de `Login` y `GoogleLogin` en `AuthPublicController`.
* Robustecer la validación de `ValidateEmail` y evitar nil pointer cuando el usuario no existe.
* Normalizar y documentar roles/constantes entre `common/utils` y `common/middleware`.
* Ya se implementó soporte de `registration_source` y `source_id` para diferenciar registros self/gym/trainer.
* Ya se añadieron modelos de asociación para `TrainerProfile`, `TrainerClient` y `GymClient` en el registro.
* Extender el modelo para manejar explícitamente asociaciones `coach-client`, `gym-coach`, `gym-client`.
* Añadir refresh tokens y expiración de access tokens en el flujo de login.
* Implementar autorización basada en permisos dinámicos además de roles estáticos.
* Mejorar el manejo de errores en controladores para diferenciar 400/401/500.
* Centralizar llamadas a servicios externos de notificación y evitar llamadas HTTP directas en el service.
* Agregar tests de integración para rutas de autenticación y middleware JWT.
* Añadir validaciones en los request DTOs y reglas de negocio de roles de registro.
* Documentar claramente el contrato de los claims JWT y los context keys usados (`user_id`, `role_id`, `access_level_id`).

## 10. ✅ Cambios recientes implementados

* Nuevo endpoint público para confirmación de email: `/public/auth/confirm?token=...`
* Nuevo endpoint público para validación de token JWT: `/public/auth/validate` (usa header Authorization o query token)
* Nuevo flujo de recuperación de contraseña:
  - `/public/auth/password/recovery`
  - `/public/auth/password/reset`
* Nuevos endpoints privados de usuario:
  - `GET /private/auth/users/:id`
  - `PUT /private/auth/users/:id`
  - `DELETE /private/auth/users/:id`
* Soporte CRUD de usuario en `authService` y repositorio:
  - `GetUserByID`
  - `UpdateUser`
  - `DeleteUser` (soft delete)
  - `ActivateUser`
  - `RequestPasswordRecovery`
  - `ResetPassword`
* Nuevos DTOs para requests/response:
  - `PasswordRecoveryRequest`
  - `PasswordResetRequest`
  - `UpdateUserRequest`
  - `GetUserResponse`
* Extensión de modelos de perfiles para gym/trainer:
  - `GymProfile` ahora guarda `city`, `department`, `country`, `primary_color`, `secondary_color`, `referral_code` y `workstation`.
  - `TrainerProfile` ahora guarda `primary_color`, `secondary_color`, `referral_code` y `files_id` para avatar.
* Se añadió la lógica de guardado de perfil de gym y perfil de entrenador en el registro.
* Se documentó este flujo adicional en el IA Memory para mantener el diseño alineado con la implementación.
* Se creó la estructura de carpeta para el endpoint de login en `src/core/login` con los submódulos `app`, `domain`, e `infra`.

## 11. ✅ Cambios más recientes (2026-04-11)

* **Flujo de confirmación de email implementado y corregido:**
  - El usuario se registra con todos sus datos incluyendo contraseña.
  - Al registrarse, `is_active: false` y `email_confirmed: false` — la cuenta queda **inactiva para todos los usuarios**.
  - Se genera JWT de activación, se guarda en `UserToken` tipo `activation`, y se envía al servicio de notificaciones.
  - El email incluye el nombre real del usuario (`FullName`) en lugar del placeholder `"ACTIVE_USER"`.
  - Al confirmar email vía `GET /public/auth/confirm?token=...`, el usuario se activa (`is_active: true`, `email_confirmed: true`) y el token se invalida.
  - El login (`POST /public/login`) ya valida `is_active` y `email_confirmed` antes de generar tokens.

* **URLs del servicio de notificaciones externalizadas:**
  - `sendConfirmationEmail` y `sendPasswordRecoveryEmail` ahora leen `NOTIFICATION_SERVICE_URL` y `DASHBOARD_URL` desde Viper (variables de entorno).
* **Validación de duplicados en registro:**
  - El registro ahora rechaza con `409 Conflict` cuando ya existe un usuario activo con el mismo email y rol.
  - Si el email existe con un rol distinto, se retorna un error específico y no se sobrescribe el usuario.
  - Fallback a `http://localhost:8443` y `http://localhost:3000` si no se configuran.
  - Nuevas variables documentadas en `.env.example`.

* **Swagger docs agregadas a `ConfirmEmail`:**
  - El endpoint `GET /public/auth/confirm` ahora tiene anotaciones Swagger completas.
  - La respuesta incluye `message` y `user` en lugar de solo el objeto usuario.

* **Corrección de imports del módulo login:**
  - Los imports en `src/core/login/` apuntaban erróneamente a `gestrym/src/core/auth/login/`.
  - Corregidos en `login_service.go`, `login_repository.go`, `login_controller.go` y `ServerRoutesDefinition.go`.
  - El proyecto compila exitosamente (`go build ./...`).

