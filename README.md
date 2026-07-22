# 🛡️ Traynova Auth Microservice (gestrym-auth)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin_Gonic-008080?style=flat-square&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/ORM-GORM_PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white)](https://gorm.io/)
[![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-orange?style=flat-square)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
[![License](https://img.shields.io/badge/License-Proprietary-red?style=flat-square)]()

Microservicio central de **Autenticación, Autorización y Gestión de Identidades** para la plataforma deportiva y nutricional **Traynova / Gestrym**. 

El sistema gestiona usuarios, roles jerárquicos, permisos dinámicos, emisión de tokens JWT, flujos de confirmación por email, restablecimiento de contraseñas y organización por grupos/sedes de entrenamiento.

---

## 📋 Tabla de Contenidos

- [📌 Visión General y Objetivos](#-visión-general-y-objetivos)
- [🏗️ Arquitectura del Sistema](#️-arquitectura-del-sistema)
- [✨ Características Principales](#-características-principales)
- [🛠️ Tecnologías y Librerías](#️-tecnologías-y-librerías)
- [📂 Estructura del Proyecto](#-estructura-del-proyecto)
- [🔐 Funcionalidades y Flujos de Negocio](#-funcionalidades-y-flujos-de-negocio)
- [🌐 Enrutamiento y Endpoints API](#-enrutamiento-y-endpoints-api)
- [⚙️ Variables de Entorno](#️-variables-de-entorno)
- [🚀 Despliegue y Ejecución Local](#-despliegue-y-ejecución-local)
- [📄 Documentación Swagger](#-documentación-swagger)

---

## 📌 Visión General y Objetivos

El microservicio `traynova.auth.back` actúa como la **fuente de verdad** para la identidad de los usuarios en la plataforma. Está diseñado bajo el patrón de **Arquitectura Hexagonal (Ports & Adapters)** para garantizar un backend completamente desacoplado de las bases de datos o los frameworks HTTP.

### Objetivos clave de negocio:
- **Gestión centralizada de identidades**: Soporta clientes, entrenadores (coaches), gimnasios (gyms) y administradores.
- **Relaciones de negocio automatizadas**: Asocia de forma transparente clientes a entrenadores (`TrainerClient`), clientes a gimnasios (`GymClient`) y entrenadores a gimnasios (`TrainerProfile`).
- **Seguridad y Control de Acceso**: Control de acceso basado en roles (RBAC) con tokens JWT y expiración configurable.
- **Notificaciones por Plantillas Estructuradas**: Integración mediante webhooks/HTTP con el servicio pro de notificaciones.

---

## 🏗️ Arquitectura del Sistema

El proyecto sigue estrictamente los principios de **Arquitectura Hexagonal**:

```
[ HTTP Requests (Gin Controllers) ]
               │
               ▼
   [ Core Application Services ]
               │
               ▼
    [ Domain Interfaces / Ports ]
               │
               ▼
[ Infrastructure Repositories (GORM / Postgres) ]
```

### Capas del Proyecto:
- **`src/common`**: Configuración global de base de datos, middlewares (JWT y RBAC), modelos compartidos de base de datos, utilidades de encripción/logging y enrutador Gin.
- **`src/core/<modulo>`**:
  - **`app`**: Lógica de aplicación, casos de uso y orquestación de servicios de dominio.
  - **`domain`**: Contratos de puertos (`ports`), DTOs de Request/Response e interfaces.
  - **`infra`**: Adaptadores concretos de infraestructura (controladores HTTP Gin y repositorios GORM).

---

## ✨ Características Principales

1. **Autenticación & Login**:
   - Login por email y contraseña hash con Bcrypt.
   - Validación de cuenta activa (`is_active: true`) y correo confirmado (`email_confirmed: true`).
   - Emisión de tokens de acceso JWT HMAC firmados con clave secreta (`JWT_KEY`).

2. **Registro Multiorigen (Self / Gym / Trainer)**:
   - Registro público para automuestreo (`registration_source: "self"`).
   - Registro delegado por Gimnasio o Entrenador (`gym` / `trainer`).
   - Creación automática de perfiles detallados (`GymProfile`, `TrainerProfile`) con colores institucionales, códigos de referido y avatares.

3. **Confirmación de Correo Electrónico**:
   - Creación de cuenta en estado inactivo (`is_active: false`, `email_confirmed: false`).
   - Emisión de `UserToken` de activación JWT.
   - Notificación vía plantilla `account_confirmation` al servicio externo de emails.
   - Confirmación por URL `GET /public/auth/confirm?token=...`.

4. **Restablecimiento y Recuperación de Contraseña**:
   - Envío de token de recuperación vía correo con plantilla `recovery_password`.
   - Restablecimiento seguro consumiendo token válido `POST /public/auth/password/reset`.

5. **Organización en Grupos y Sedes**:
   - Creación de grupos/sedes (`UserGroup` / `UserGroupMember`).
   - Permite a Gimnasios organizar entrenadores y clientes por sedes físicas.
   - Permite a Entrenadores agrupar a sus clientes en programas o equipos.

6. **Control de Estado de Usuarios (Toggle & Soft Delete)**:
   - Endpoint `PATCH /private/auth/users/:id/status` para deshabilitar/habilitar usuarios manualmente.
   - Eliminación suave (*soft delete*) preservando la integridad referencial.

---

## 🛠️ Tecnologías y Librerías

- **Lenguaje**: [Go 1.25+](https://go.dev/)
- **Framework Web HTTP**: [Gin Gonic v1.12](https://github.com/gin-gonic/gin)
- **ORM / Base de Datos**: [GORM v1.25](https://gorm.io/) con driver PostgreSQL (`pgx/v5`)
- **Configuración**: [Viper v1.21](https://github.com/spf13/viper) (soporta YAML y variables de entorno)
- **Seguridad & Crypt**: [Golang JWT v4](https://github.com/golang-jwt/jwt) y [Bcrypt](https://golang.org/x/crypto)
- **Documentación API**: [Gin-Swagger / Swaggo v1.16](https://github.com/swaggo/gin-swagger)

---

## 📂 Estructura del Proyecto

```
traynova.auth.back/
├── .env.example                # Plantilla de variables de entorno
├── Dockerfile                  # Construcción de imagen de contenedor
├── GUIA_FRONT.md               # Guía de consumo de API para desarrolladores Frontend
├── IA_MEMORY.md                # Memoria de arquitectura y decisiones técnicas para IA
├── CLASIFICACION_Y_DETALLE.md  # Desglose exhaustivo de componentes y mapa del repositorio
├── main.go                     # Punto de entrada de la aplicación
├── go.mod                      # Módulo de Go y dependencias
├── deployment/
│   └── env_local.yaml          # Configuración de entorno local
├── docs/                       # Especificaciones OpenAPI / Swagger generadas
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── src/
    ├── app.go                  # Inicializador del servidor y migración DB
    ├── common/
    │   ├── config/             # Conexión DB y auto-migraciones GORM
    │   ├── middleware/         # Auth JWT Middleware & Role Authorization
    │   ├── models/             # Entidades GORM relacionales
    │   ├── routes/             # Definición centralizada de rutas HTTP
    │   └── utils/              # Hashers, logger y constantes de roles
    └── core/                   # Módulos de dominio (Hexagonal Architecture)
        ├── access_levels/      # Niveles de acceso
        ├── actions/            # Acciones del sistema
        ├── auth/               # Servicio principal (Registro, Usuarios, Grupos)
        ├── jwt/                # Generador de tokens JWT
        ├── login/              # Controladores y servicio de inicio de sesión
        ├── permissions/        # Gestión de permisos
        ├── roles/              # Gestión de roles
        └── token_types/        # Gestión de tipos de token
```

---

## 🌐 Enrutamiento y Endpoints API

### 🔓 Endpoints Públicos (`/gestrym-auth/public/*`)
- `POST /public/login` - Inicio de sesión y emisión de tokens.
- `POST /public/auth/register` - Autorregistro público de usuarios.
- `GET  /public/auth/confirm` - Confirmación de cuenta mediante token enviado por correo.
- `POST /public/auth/validate` - Validación de token JWT.
- `POST /public/auth/password/recovery` - Solicitud de correo de recuperación de contraseña.
- `POST /public/auth/password/reset` - Restablecimiento de contraseña con token.

### 🔒 Endpoints Privados (`/gestrym-auth/private/*`)
- `POST   /private/auth/register` - Registro de usuarios por parte de Entrenadores o Gimnasios.
- `GET    /private/auth/users` - Listado de usuarios activos con filtros (soporta `group_id`).
- `GET    /private/auth/users/:id` - Obtener detalle de usuario por ID.
- `PUT    /private/auth/users/:id` - Actualizar información de usuario.
- `DELETE /private/auth/users/:id` - Soft delete de usuario.
- `PATCH  /private/auth/users/:id/status` - Activar/Desactivar estado de usuario.
- `GET    /private/auth/clients` - Listar clientes/entrenadores vinculados según el rol del solicitante.
- `POST   /private/auth/groups` - Crear grupo o sede.
- `GET    /private/auth/groups` - Listar grupos o sedes creados por el usuario.
- `PUT    /private/auth/groups/:id` - Editar grupo/sede y sus miembros.
- `DELETE /private/auth/groups/:id` - Eliminar grupo o sede.

### 🛡️ Endpoints Protegidos / Catálogos (`/gestrym-auth/protected/*`)
- Catálogos administrativos protegidos por rol Admin (`RequireRoles(4)`): `/roles`, `/permissions`, `/actions`, `/access-levels`, `/token-types`.

---

## ⚙️ Variables de Entorno

Configurables mediante archivo `.env` o variables de ambiente del contenedor:

| Variable | Descripción | Valor por Defecto / Ejemplo |
| :--- | :--- | :--- |
| `GESTRYM_SERVER_ADDRESS` | Puerto y host donde corre la API | `:8080` |
| `GIN_MODE` | Modo de ejecución de Gin (`debug` / `release`) | `debug` |
| `GESTRYM_DB_HOST` | Host de la base de datos PostgreSQL | `localhost` |
| `GESTRYM_DB_USER` | Usuario PostgreSQL | `postgres` |
| `GESTRYM_DB_PASSWORD` | Contraseña PostgreSQL | `postgres` |
| `GESTRYM_DB_NAME` | Nombre de la base de datos | `gestrym_auth_db` |
| `GESTRYM_DB_PORT` | Puerto PostgreSQL | `5432` |
| `JWT_KEY` | Clave secreta para firmar tokens JWT | `secret_key_change_me` |
| `NOTIFICATION_SERVICE_URL` | Endpoint del servicio de notificaciones pro | `http://localhost:8443` |
| `DASHBOARD_URL` | URL base de la aplicación web / dashboard | `http://localhost:3000` |
| `X_API_KEY` | Llave API para autorizar el envío de notificaciones | `your_api_key` |

---

## 🚀 Despliegue y Ejecución Local

### 1. Prerrequisitos
- [Go 1.25+](https://go.dev/dl/) instalado.
- Instancia activa de **PostgreSQL**.

### 2. Clonar e Instalar Dependencias
```bash
git clone https://github.com/tu-usuario/traynova.auth.back.git
cd traynova.auth.back
go mod download
```

### 3. Ejecución Local (Desarrollo)
Crea o edita `deployment/env_local.yaml` con tus credenciales locales de PostgreSQL y ejecuta:
```bash
go run main.go --local=true
```

### 4. Ejecución en Producción / Docker
Construir y correr la imagen Docker:
```bash
docker build -t traynova-auth-back .
docker run -p 8080:8080 --env-file .env traynova-auth-back
```

---

## 📄 Documentación Swagger

El microservicio expone su documentación Swagger generada dinámicamente. 

Con el servidor corriendo, accede en el navegador a:
👉 **`http://localhost:8080/swagger/index.html`** (o en la ruta configurada `/gestrym-auth/swagger/index.html`).

---

<p align="center">
  Desarrollado con ❤️ para la plataforma <b>Traynova</b>.
</p>
