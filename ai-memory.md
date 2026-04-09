# Memoria Técnica Estructurada - Gestrym (Auth)

## 1. Descripción General del Proyecto

* **Tipo de sistema:** Servicio backend / API REST. Modular (microservicios posibles) con BasePath `/gestrym-auth`.
* **Propósito del servicio actual:** Proveer rutas de autenticación y manejo de jerarquías (jwt, roles, permisos, users, catalogs, login/registro).
* **Tecnologías utilizadas:** Golang, Gin, Postgres, GORM, Viper, implementaciones custom de Logging, Swaggo (Swagger UI) y bcrypt.

## 2. Arquitectura (Diseño Personalizado de Capas)

La arquitectura sigue una estructura inspirada en capas y Clean Architecture pero con una decisión crucial del desarrollador para evitar el acoplamiento cruzado:

* **Entidades Persistentes (`common/models`)**: Todas las entidades y modelos de base de datos (`GORM structs`) se definen de manera central y transversal en este paquete. 
* **Puertos y Modelos Transitorios (Módulos Individuales - `domain`)**: Dentro del módulo (por ejemplo `core/users/domain`), residen:
  * **Ports (Interfaces)**: Abstracciones que permiten aislar el repositorio (ej. `domain/ports/IUser_repository.go`).
  * **Structs (VMs, DTOs)**: Modelos de vista, peticiones (Request) y respuestas (Response).
* **Application (`app`)**: Lógica de negocios que orquesta a través de las interfaces (`user_service.go`).
* **Infrastructure (`infra`)**: Implementaciones técnicas (ej. Repositorios de BD y Controladores HTTP). Todo acoplamiento al framework va aquí.

## 3. Estructura del Proyecto

* **Estructura de carpetas explicada:**
  * `/main.go` y `/src/app.go`: Puntos de entrada del servicio HTTP.
  * `/src/common/`: Infraestructura y Modelos de dominio global.
    * `/models`: Modelos de GORM (`Users.go`, `Role.go`, `Permission.go`, `RolePermission.go`).
    * `/config`: Conexión de base de datos (`Database.go`) y variables (`Enviroment.go`).
    * `/middleware`: Proxys JWT y validación.
  * `/src/core/<module>/`: Módulos de negocio (ej. `users`).
    * `/domain`: Declaración de reglas y formas ajenas a DB (`/ports`, structs de DTO).
    * `/app`: Servicios de implementación de casos de uso.
    * `/infra`: Puentes a la DB (`repository`) y a la red (`controller`).

## 4. Dominio (Domain Layer)

* **Entidades Globales (`common/models`)**: `User`, `Role`, `Permission`, `RolePermission`. Recientemente expandido con `UserToken`, `RefreshToken`, `UserTokenType`, `Action` y `AccessLevel`. Albergadas de manera central, abstraen relaciones jerárquicas y ciclos de vida dinámicos de tokens.
* **Structs / ViewModels**: Definidos (o próximos a ubicar de forma completa) en la carpeta del dominio local respectivo de su controlador.
* **Reglas de Negocio:** Hashing con bcrypt necesario de la contraseña, JWT caducando a las 24 hrs. Identificación unívoca basada en la interfaz y relaciones entre `Role` y `Permission`.

## 5. Casos de Uso (Application Layer) y Lógica de Negocio

* **Crear / Registrar Usuarios (`RegisterUser / CreateUser`):**
  * *Flujo sin contraseña inicial:* En la creación o registro, el servicio solicita nombres, correo y teléfono en lugar de contraseña. Al instante, se dispara una petición al microservicio `gestrym-notification` adjuntando un `confirm_token` seguro para notificar al usuario.
  * **Registro Público (`/public/auth/register`)**: Permite únicamente registrar usuarios con `role_id` 1 (Cliente) o 2 (Coach).
  * **Registro Privado Interno (`/private/users/register`)**: Guardado por middleware de roles. Valida matriz jerárquica: Un Admin (4) puede crear a todos. Un Gym (3) solo crea Coach o Clientes. Y un Coach (2) solo aprueba Clientes.
* **Autenticación (`Login / GoogleLogin / Refresh / Logout`):**
  * Login tradicional: Valida hash bcrypt, genera JWT, **y almacena el estado activo en base de datos (`UserToken`)** adjuntando un `RefreshToken` criptográfico propio.
  * Login con Google: Valida un `id_token` de Firebase/Google a través del paquete `google.golang.org/api/idtoken`. Carga auto-registro invisible si no existe, forzando rol de `Cliente` (1). Tras verificar, repite el proceso de anclar tokens locales a DB.
  * Ciclo de Vida: El token de estado es validado a voluntad en endpoints `/auth/refresh` y destrozado/marcado revocado en `/auth/logout`.
* **Catálogos Dinámicos:** Módulos que operan las opciones del proyecto para Actions, Niveles de Acceso y Tipos de Token (asegurados sólo para Admin en sus respectivos Endpoints).
* **GetMe**: Retorna datos del usuario basado en parseo del access token (`auth/validate`).

## 6. Infraestructura

* **Persistencia:** Aislada por inyección de interfaces bajo `repositories.go` (implementaciones) usando GORM.
* **Middlewares:** Setup genérico e inyección contextual de `user_id` desde el header de `Bearer Auth`.

## 7. API / Interfaces HTTP

* Rutas Públicas (`/gestrym-auth/public/auth/...`): Exponen Login y Registro, parseando DTOs desde JSON.
* Rutas Privadas Extensas: `/protected` (por API key) y `/private` (Por Bearer JWT).
* **Swagger Dinámico:** Se levanta en base path (ej: `/gestrym-auth/swagger/index.html`) proveyendo documentación automática generada con anotaciones GoDoc `swag init` alojadas en los handlers. Control central ruteador en `ServerRoutesDefinition.go`.

## 8. Seguridad

* Prevención estricta: Toda contraseña viaja hasta la `app layer` y allí es obligatoriamente pre-procesada por `bcrypt` antes de la persistencia.
* Fuerte control de Expiración JWT y separación en la validez del método criptográfico de hash de JWT (solo HS256).

## 9. Convenciones del Desarrollador (Regla de Oro)

1. **Centralización de Entidades:** **NO CREAR** modelos GORM dentro de las carpetas de negocio. Toda tabla de persistencia debe declararse en `src/common/models` y auto-migrarse globalmente. Esto elimina el espagueti relacional en GORM ante cruce de módulos.
2. **Separación de Tránsito:** La carpeta `domain` dentro de mi core se usa para alojar los **Puertos (Interfaces)** en `domain/ports` y los **Structs temporales (DTOs, Vms, Request/Response)** para mantener purificado el handler del Controller.

## 10. Reglas para Extender el Sistema

* **Pasos para crear nueva funcionalidad o endpoint:**
  1. Si hay tablas implicadas, modelarlo en `src/common/models` y auto migrarlo.
  2. Crear los Request/Response `structs` en `src/core/[modulo]/domain/[structs_o_vm]`.
  3. Declarar la interfaz Port si se consume un Repository o External Service (`domain/ports`).
  4. Implementar funcionalidad en `app/` recibiendo las abstracciones.
  5. Acoplar redimensionador y request parser HTTP en el Controler y enrutar en `ServerRoutesDefinition.go`.

## 11. Configuración de Despliegue y Control de Versiones

* **Variables de Entorno y Git ( `.gitignore` ):**
  Se centralizó que los archivos `.env`, así como logs (`*.log`), binarios (`*.exe`, `*.so`), directorios IDE (`.idea/`, `.vscode/`) y la carpeta dinámica `dist/` e incluso `node_modules/` (en caso de convivir con JS) queden estrictamente excluidos del respositorio para evitar subidas inseguras.
* **Conteneurización ( `Dockerfile` ):**
  Se reestructuró el despliegue para usar una **Construcción Multi-Etapa (Multi-stage Build)** puramente en Golang:
  1. Utiliza `golang:1.25.0-alpine` como **Builder** de los archivos binarios (`CGO_ENABLED=0` para portabilidad máxima en Alpine linux).
  2. Produce un ejecutable de poco tamaño insertado en un nuevo ambiente base ligero `alpine:latest` provisto con instaladores CAs (`ca-certificates tzdata`) para habilitar requests SSL salientes (necesarios para OAuth de Google o Webhooks al MS Node.js).
  3. Esto asegura imágenes mucho más compactas, seguras y especializadas, abandonando por completo ecosistemas innecesarios ajenos a Go.
