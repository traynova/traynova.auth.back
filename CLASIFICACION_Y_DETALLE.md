# 🔍 Clasificación y Detalle del Repositorio: `traynova.auth.back`

Este documento proporciona un desglose técnico, taxonómico y funcional completo de todos los archivos, componentes, módulos y modelos presentes en el repositorio **`traynova.auth.back`** (módulo Go `gestrym`).

---

## 📐 1. Clasificación General del Repositorio

El repositorio es un **microservicio en Go (v1.25+)** especializado en la **gestión de autenticación, autorización (RBAC), registro multiorigen e identidades** para la plataforma deportiva y nutricional **Traynova / Gestrym**.

### Taxonomía de Archivos y Carpetas

| Categoría | Ubicación | Descripción |
| :--- | :--- | :--- |
| **Punto de Entrada** | `main.go`, `src/app.go` | Inicializan la aplicación, leen argumentos CLI, cargan Viper, ejecutan migraciones GORM y levantan el servidor HTTP Gin. |
| **Núcleo de Dominio** | `src/core/` | Contiene los submódulos de negocio (`auth`, `login`, `jwt`, `roles`, `permissions`, `actions`, `access_levels`, `token_types`) organizados en **Arquitectura Hexagonal**. |
| **Capa Compartida** | `src/common/` | Reúne infraestructura común: conexión DB, auto-migración, middlewares JWT/CORS/RBAC, modelos relacionales GORM y utilidades. |
| **Configuración** | `.env`, `.env.example`, `deployment/env_local.yaml` | Define parámetros de ejecución local y de producción (credenciales DB, secretos JWT, URLs de servicios). |
| **Despliegue** | `Dockerfile` | Define la compilación e imagen Docker para entornos de producción (Docker / Cloud / Render). |
| **Documentación** | `README.md`, `GUIA_FRONT.md`, `IA_MEMORY.md`, `docs/` | Documentación técnica del proyecto, especificación OpenAPI/Swagger, manual de integración para Frontend y memoria para IA. |

---

## 🏗️ 2. Clasificación por Capas de Arquitectura (Hexagonal)

El código fuente en `src/` está estrictamente dividido respetando el patrón **Ports & Adapters (Arquitectura Hexagonal)**:

```
src/
├── common/             <-- Capa Cross-Cutting / Infraestructura Compartida
│   ├── config/         <-- Conexión a Base de Datos PostgreSQL y AutoMigrate
│   ├── middleware/     <-- Middlewares HTTP (JWT Auth, RBAC Role Validation, CORS)
│   ├── models/         <-- Definiciones GORM de Tablas / Entidades Relacionales
│   ├── routes/         <-- Enrutador Gin y mapeo de grupos (public, private, protected)
│   ├── shared/         <-- Interfaces o estructuras compartidas globales
│   └── utils/          <-- Encriptación Bcrypt, Logger, Constantes de Roles
└── core/               <-- Módulos de Dominio Especifico
    ├── auth/           <-- Submódulo principal (Registro, Confirmación, Grupos, Usuarios)
    ├── login/          <-- Submódulo de inicio de sesión y autenticación
    ├── jwt/            <-- Submódulo de generación y registro de tokens
    ├── roles/          <-- Submódulo CRUD para Roles (RBAC)
    ├── permissions/    <-- Submódulo CRUD para Permisos (RBAC)
    ├── actions/        <-- Submódulo CRUD para Acciones del Sistema
    ├── access_levels/  <-- Submódulo CRUD para Niveles de Acceso
    └── token_types/    <-- Submódulo CRUD para Tipos de Token (activación, recuperación)
```

En cada submódulo dentro de `src/core/` se cumple la tripleta:
- **`app/` (Application Layer)**: Implementa los servicios del sistema y la lógica de caso de uso (ej. `auth_service.go`, `login_service.go`).
- **`domain/` (Domain Layer)**: Contiene los puertos/interfaces (`IAuthRepository`, `ILoginRepository`) y las estructuras DTO de entrada y salida (`structs/`).
- **`infra/` (Infrastructure Layer)**: Contiene los adaptadores concretos:
  - Controladores HTTP Gin (ej. `auth_public_controller.go`, `auth_private_controller.go`).
  - Repositorios de base de datos GORM (ej. `auth_repository.go`, `login_repository.go`).

---

## 🗄️ 3. Clasificación del Modelo de Base de Datos (Entidades GORM)

Las entidades están definidas en `src/common/models/` y traducidas a tablas PostgreSQL mediante GORM:

### A. Sistema de Identidades y Usuarios
- **`User`**: Entidad principal que almacena `email`, `password_hash`, `full_name`, `phone_prefix`, `phone_number`, `is_active`, `email_confirmed`, `role_id` y marcas de tiempo (`created_at`, `updated_at`, `deleted_at` para soft delete).
- **`GymProfile`**: Información extendida para gimnasios (`city`, `department`, `country`, `primary_color`, `secondary_color`, `referral_code`, `workstation`).
- **`TrainerProfile`**: Información extendida para entrenadores (`primary_color`, `secondary_color`, `referral_code`, `files_id` para avatar).

### B. Control de Acceso Basado en Roles (RBAC)
- **`Role`**: Roles del sistema (`Client`=1, `Coach`=2, `Gym`=3, `Admin`=4).
- **`Permission`**: Permisos asociados a un rol y recurso.
- **`RolePermission`**: Tabla pivote para la relación muchos a muchos entre `Role` y `Permission`.
- **`Action`**: Operaciones registrables sobre recursos.
- **`AccessLevel`**: Definición de grados o niveles de acceso.

### C. Relaciones de Negocio Deportivo
- **`TrainerClient`**: Mapea la relación directa entre un Entrenador (`trainer_id`) y su Cliente (`client_id`).
- **`GymClient`**: Mapea la relación directa entre un Gimnasio (`gym_id`) y su Cliente (`client_id`).
- **`TrainerProfile` (Gym Binding)**: Asocia un entrenador a un gimnasio a través del campo `gym_id`.

### D. Grupos y Sedes de Entrenamiento
- **`UserGroup`**: Define agrupaciones de usuarios creadas por un Gimnasio (ej. "Sede Norte", "Sede Centro") o por un Entrenador (ej. "Grupo Maratón").
- **`UserGroupMember`**: Tabla relacional que vincula un `group_id` con un `user_id` (pudiendo ser clientes o entrenadores).

### E. Tokens de Seguridad y Estado
- **`UserTokenType`**: Catálogo de tipos de token (1: activación de cuenta, 2: recuperación de contraseña).
- **`UserToken`**: Registra los tokens JWT emitidos para un usuario, su tipo, estado de uso (`is_used`) y fecha de expiración.

---

## 🔑 4. Explicación Detallada de los Flujos de Negocio

### 1. Flujo de Registro y Confirmación de Cuenta
```
[ Frontend / App ] ──POST /public/auth/register──> [ AuthPublicController ]
                                                             │
                                                             ▼
                                                    [ AuthService.RegisterUser ]
                                                             │
                                      ┌──────────────────────┴──────────────────────┐
                                      ▼                                             ▼
                             [ Se crea User inactivo ]                   [ Se crea perfil según rol ]
                             (is_active: false)                          (GymProfile / TrainerProfile)
                             (email_confirmed: false)                    (TrainerClient / GymClient)
                                      │
                                      ▼
                             [ Genera JWT Activación ]
                             (UserToken tipo activation)
                                      │
                                      ▼
                             [ Servicio Notificaciones ] ──HTTP POST──> [ Template account_confirmation ]
```
- **Paso 1**: El usuario se registra enviando datos de perfil y origen (`registration_source`: `"self"`, `"gym"`, `"trainer"`).
- **Paso 2**: El usuario se almacena en base de datos con **contraseña hasheada con Bcrypt**, pero en estado **inactivo** (`is_active: false`) y **sin confirmar** (`email_confirmed: false`).
- **Paso 3**: Se genera un token JWT de activación registrado en `UserToken`.
- **Paso 4**: El servicio realiza una petición HTTP al microservicio de notificaciones usando el formato unificado de plantilla `account_confirmation` con la URL de confirmación.
- **Paso 5**: Al hacer clic en el enlace (`GET /public/auth/confirm?token=...`), el backend valida el token, marca `is_active: true`, `email_confirmed: true` e invalida el token.

### 2. Flujo de Inicio de Sesión (Login)
- **Ruta**: `POST /public/login`
- **Proceso**:
  1. Recibe email y contraseña.
  2. Busca al usuario en la base de datos por email.
  3. **Verifica estado**: Si `is_active` o `email_confirmed` es falso, rechaza el inicio de sesión.
  4. Compara la contraseña recibida contra el `password_hash` mediante `bcrypt.CompareHashAndPassword`.
  5. Si es correcta, genera un Access Token JWT con los claims `user_id`, `role_id` y `access_level_id`.

### 3. Flujo de Recuperación de Contraseña
- **Solicitud**: `POST /public/auth/password/recovery` con el email del usuario.
- **Generación**: Se crea un `UserToken` de recuperación y se envía la plantilla de correo `recovery_password` con el enlace de restauración.
- **Restablecimiento**: `POST /public/auth/password/reset` recibe el token de recuperación y la nueva contraseña. Se valida el token, se actualiza el hash de la contraseña en la BD y se marca el token como consumido (`is_used: true`).

### 4. Flujo de Grupos, Sedes y Control de Acceso Jerárquico
- **Grupos / Sedes**:
  - `POST /private/auth/groups`: Permite a Gimnasios o Entrenadores crear grupos/sedes.
  - `GET /private/auth/users?group_id=...`: Permite filtrar el listado de usuarios por sede o grupo.
- **Control de Jerarquía y Pertenencia**:
  - Los Entrenadores solo pueden visualizar y editar a los clientes asignados directamente a su código/ID (`TrainerClient`).
  - Los Gimnasios pueden visualizar a sus entrenadores y a los clientes asociados tanto directa como indirectamente (`GymClient` + `TrainerClient` vinculados a sus entrenadores).
  - Los Administradores (`role_id: 4`) poseen visibilidad global.

---

## 📄 5. Clasificación de Archivos y Documentación del Proyecto

### Archivos de Documentación
1. **`README.md`**: Presentación general del proyecto para GitHub, instrucciones de instalación, variables de entorno y mapa de rutas.
2. **`GUIA_FRONT.md`**: Manual técnico detallado para desarrolladores Frontend. Explica las respuestas JSON, errores HTTP (400, 401, 409, 500) y cómo consumir los endpoints con cabeceras `Authorization: Bearer <token>`.
3. **`IA_MEMORY.md`**: Documento vivo que preserva las decisiones arquitectónicas, reglas de no-rotura, cambios recientes y contexto histórico para ser leído por asistentes IA.
4. **`CLASIFICACION_Y_DETALLE.md`** *(este documento)*: Inventario y desglose exhaustivo del repositorio.

### Archivos de Configuración y Despliegue
- **`deployment/env_local.yaml`**: Archivo de variables YAML utilizado cuando se ejecuta en desarrollo local (`--local=true`).
- **`.env.example`**: Plantilla con las variables de entorno requeridas para ejecución en contenedores o servidores de integración continua.
- **`Dockerfile`**: Script de construcción en multi-stage para generar un binario ligero de Go y desplegar en entornos como Render, Kubernetes o AWS.

---

## ⚠️ 6. Restricciones y Reglas de Negocio a Preservar

1. **Arquitectura Hexagonal Estricta**: No realizar consultas a la base de datos directamente desde los controladores ni desde la capa `common`. Todo acceso a datos debe pasar por `app service` -> `domain port` -> `infra repository`.
2. **Estado de Usuario al Registrar**: Todo nuevo usuario registrado DEBE crearse deshabilitado (`is_active: false`) hasta que confirme su correo electrónico.
3. **Hasheo Obligatorio de Contraseñas**: Nunca almacenar ni transmitir contraseñas en texto plano; usar siempre el paquete `common/utils` con Bcrypt.
4. **Seguridad JWT**: Ningún endpoint bajo `/private/` o `/protected/` debe omitir el middleware `JWTModdleware.go`.

---

## 📈 7. Oportunidades de Evolución y Mejoras Futuras

- **Implementación de Refresh Tokens**: Incorporar rotación de `refresh_token` con expiración prolongada y `access_token` de vida corta.
- **OAuth2 / Google Sign-In**: Finalizar e integrar completamente la autenticación con Google (código actualmente preparado en controlador).
- **Pruebas Automatizadas**: Agregar suites de tests unitarios y de integración para la capa de servicios y controladores con `testify` o `go test`.
- **Métricas y Telemetría**: Integrar Prometheus/OpenTelemetry para monitoreo de llamadas a endpoints de autenticación.
