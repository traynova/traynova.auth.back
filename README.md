# 🔐 Traynova / Gestrym - Microservicio de Autenticación (Auth Back)

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Framework](https://img.shields.io/badge/Gin--Gonic-v1.12.0-008080?style=for-the-badge&logo=gin&logoColor=white)
![ORM](https://img.shields.io/badge/GORM-v1.25.10-blue?style=for-the-badge)
![Database](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Architecture](https://img.shields.io/badge/Architecture-Hexagonal-orange?style=for-the-badge)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI%203.0-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

**Traynova Auth Back** (`gestrym-auth`) es el microservicio centralizado de **Autenticación, Autorización y Gestión de Identidades (IAM)** para la plataforma de entrenamiento deportivo y gestión nutricional Traynova / Gestrym.

Proporciona un mecanismo seguro y escalable basado en **Arquitectura Hexagonal (Ports & Adapters)** para gestionar usuarios, roles, permisos granularizados (RBAC), tokens de activación y recuperación de contraseña, relaciones de negocio (Gimnasio, Entrenador, Cliente) y agrupaciones/sedes.

---

## 🚀 Características Principales

- 👤 **Registro Multiorigen:**
  - Registro público (*Self-signup*).
  - Registro de clientes o entrenadores por parte de un Gimnasio o Entrenador (*Gym / Trainer source*).
- 📧 **Confirmación de Cuenta por Email:**
  - Registro con cuenta inactiva (`is_active: false`, `email_confirmed: false`).
  - Generación de token JWT de activación invalidadle tras su uso.
  - Integración mediante servicio de notificaciones pro por plantillas HTML.
- 🔑 **Autenticación y Tokens JWT:**
  - Emisión y validación de tokens JWT HMAC firmados con claims contextuales (`user_id`, `role_id`, `access_level_id`).
  - Endpoint de validación pública de tokens.
- 🔐 **Control de Acceso Basado en Roles (RBAC):**
  - Roles soportados: **Cliente (1)**, **Coach/Entrenador (2)**, **Gimnasio (3)**, **Administrador (4)**.
  - Middlewares de protección de rutas por JWT y verificación estricta de roles.
- 🔄 **Flujo de Recuperación de Contraseña:**
  - Solicitud de recuperación vía correo electrónico.
  - Restablecimiento mediante token con expiración programada.
- 🏢 **Gestión de Grupos y Sedes (`UserGroup` / `UserGroupMember`):**
  - Creación y administración de grupos o sedes para organizar clientes y personal deportivo.
  - Filtrado de usuarios por sede/grupo.
- 🛡️ **Seguridad y Auditoría:**
  - Hashing seguro de contraseñas con `bcrypt`.
  - Habilitación/deshabilitación manual de usuarios (`ToggleUserStatus`).
  - *Soft delete* de usuarios para preservar integridad histórica.

---

## 🏗️ Arquitectura del Proyecto

El microservicio está implementado siguiendo los principios de la **Arquitectura Hexagonal (Ports & Adapters)**, garantizando el desacoplamiento entre las reglas de dominio y los detalles técnicos como la base de datos o el framework HTTP.

```text
src/
├── app.go                       # Punto de entrada principal y bootstrap del servidor Gin
├── common/                      # Capa transversal y modelos compartidos
│   ├── config/                  # Carga de variables de entorno (Viper)
│   ├── middleware/              # JWT, CORS y middleware de validación de roles
│   ├── models/                  # Entidades GORM compartidas
│   ├── routes/                  # Definición global de grupos de rutas (Public, Private, Protected)
│   ├── shared/                  # Cliente HTTP e integración externa (Notificaciones)
│   └── utils/                   # Utilidades, constantes y formateadores
└── core/                        # Bounded Contexts / Módulos de Dominio
    ├── access_levels/           # Gestión de Niveles de Acceso
    ├── actions/                 # Gestión de Acciones de Permisos
    ├── auth/                    # Módulo principal de Usuarios, Registro y Perfiles
    │   ├── app/                 # Servicios de aplicación y lógica de negocio
    │   ├── domain/              # Puertos (Interfaces) y DTOs
    │   └── infra/               # Adaptadores HTTP (Controllers) y Repositorios GORM
    ├── jwt/                     # Generación y verificación de Tokens JWT
    ├── login/                   # Módulo de Autenticación y Login
    ├── permissions/             # Definición y asignación de Permisos
    ├── roles/                   # Definición de Roles del Sistema
    └── token_types/             # Catálogo de tipos de token (Activación, Reset, etc.)
```

---

## 🛠️ Tecnologías y Librerías

| Tecnología | Categoría | Propósito |
| :--- | :--- | :--- |
| **Go 1.25** | Lenguaje | Lenguaje de programación base |
| **Gin Gonic v1.12** | Web Framework | Enrutamiento HTTP rápido y ligero |
| **GORM v1.25** | ORM | Mapeo objeto-relacional para PostgreSQL |
| **PostgreSQL** | Base de Datos | Persistencia relacional de datos |
| **Viper v1.21** | Configuración | Gestión de variables de entorno y archivos `.env` |
| **Golang JWT v4** | Autenticación | Emisión y verificación de JSON Web Tokens |
| **Bcrypt** | Seguridad | Encriptación y hashing de contraseñas |
| **Gin Swagger** | Documentación | Generación automática de especificación OpenAPI / Swagger |

---

## 🌐 Resumen de Endpoints HTTP

### 🔓 Endpoints Públicos (`/public`)
- `POST /public/auth/register`: Registro de nuevos usuarios (*Self-signup*).
- `GET /public/auth/confirm`: Confirmación de email mediante token de activación.
- `GET /public/auth/validate`: Validación de integridad de tokens JWT.
- `POST /public/auth/password/recovery`: Solicitud de enlace de recuperación de contraseña.
- `POST /public/auth/password/reset`: Restablecimiento de contraseña con token.
- `POST /public/login`: Inicio de sesión y generación de tokens de acceso.

### 🔒 Endpoints Privados / Autenticados (`/private`)
- `POST /private/auth/register`: Registro interno de clientes/entrenadores por Gym o Coach.
- `GET /private/auth/users`: Listado de usuarios con filtros (rol, búsqueda, `group_id`).
- `GET /private/auth/users/:id`: Obtener detalles de un usuario específico.
- `PUT /private/auth/users/:id`: Actualizar datos de usuario.
- `DELETE /private/auth/users/:id`: Desactivación (*soft delete*) de usuario.
- `PATCH /private/auth/users/:id/status`: Habilitar o deshabilitar cuenta de usuario.
- `GET /private/auth/clients`: Listado de clientes asociados al entrenador o gimnasio autenticado.
- `POST /private/auth/groups`: Crear grupo o sede.
- `GET /private/auth/groups`: Listar grupos o sedes pertenecientes al usuario autenticado.
- `PUT /private/auth/groups/:id`: Actualizar grupo y asignación de miembros.
- `DELETE /private/auth/groups/:id`: Eliminar un grupo/sede.

### 🛡️ Endpoints Protegidos / Catálogos Admin (`/protected`)
- Rutas administrativas protegidas por rol Admin (`RequireRoles(4)`) para gestionar `roles`, `permissions`, `actions`, `access_levels` y `token_types`.

---

## ⚙️ Configuración y Variables de Entorno

Crea un archivo `.env` en la raíz del proyecto basado en `.env.example`:

```env
# Servidor HTTP
PORT=8080
ENV=local

# Base de Datos PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=traynova_auth
DB_SSLMODE=disable

# Autenticación JWT
JWT_KEY=tu_clave_secreta_jwt_super_segura
JWT_EXPIRATION_HOURS=24

# Microservicio de Notificaciones Externo
NOTIFICATION_SERVICE_URL=http://localhost:8081
X_API_KEY=tu_api_key_de_notificaciones
DASHBOARD_URL=http://localhost:3000
```

---

## ⚡ Ejecución Local

### 1. Clonar el repositorio e instalar dependencias
```bash
git clone https://github.com/tu-usuario/traynova.auth.back.git
cd traynova.auth.back
go mod download
```

### 2. Ejecutar la aplicación
```bash
go run main.go
```
El servidor iniciará por defecto en `http://localhost:8080`.

### 3. Ejecutar con Docker
```bash
docker build -t gestrym-auth .
docker run -p 8080:8080 --env-file .env gestrym-auth
```

---

## 📚 Documentación Swagger / OpenAPI

La documentación interactiva de la API está integrada mediante Swagger UI:
- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **Especificación JSON:** [swagger.json](file:///Users/jhonnier.gomez/Documents/Personal/traynova.auth.back/docs/swagger.json)

Para actualizar la documentación tras modificar comentarios en el código:
```bash
swag init -g src/app.go
```

---

## 📄 Guías Adicionales

- 📘 [CLASIFICACION_Y_DETALLE.md](file:///Users/jhonnier.gomez/Documents/Personal/traynova.auth.back/CLASIFICACION_Y_DETALLE.md): Explicación exhaustiva componente por componente, entidades de dominio y flujos.
- 📱 [GUIA_FRONT.md](file:///Users/jhonnier.gomez/Documents/Personal/traynova.auth.back/GUIA_FRONT.md): Guía detallada de integración para aplicaciones cliente (Web / Mobile).
- 🧠 [IA_MEMORY.md](file:///Users/jhonnier.gomez/Documents/Personal/traynova.auth.back/IA_MEMORY.md): Memoria técnica y decisiones de arquitectura para asistentes IA.

---

## 📝 Licencia y Autores

Desarrollado como parte del ecosistema de microservicios **Traynova / Gestrym**. Todos los derechos reservados.
