# PetMatch API (Go + Gin)

## Arquitectura general
- **Framework**: Gin (HTTP) + GORM (ORM) sobre SQLite por defecto.
- **Estructura**: `cmd/` para el bootstrap, `internal/` con capas separadas de config, database, models, repositories, services, handlers, middleware y router.
- **Autenticacion**: JWT firmado (HS256) con expiracion de 24h. El secreto, puerto, ruta de base de datos y credenciales del admin se leen desde variables de entorno (`PETMATCH_*`).
- **Migraciones**: `database.Migrate` ejecuta `AutoMigrate` para User, Pet y AdoptionRequest al iniciar el servicio.

## Modelos clave
- `User`: roles `adopter`, `shelter`, `admin`; refugios requieren aprobacion manual (`is_approved`).
- `Pet`: perfiles publicados por refugios, con estado (`available`, `adopted`).
- `AdoptionRequest`: solicitudes con estados `pending`, `approved`, `rejected`.

## Endpoints principales (`/api/v1`)
- `POST /auth/register` – Registro de adoptantes/refugios (hash bcrypt).
- `POST /auth/login` / `GET /auth/me` – Inicio de sesion y recuperacion del usuario autenticado.
- `GET /pets` / `GET /pets/{id}` – Catalogo publico con filtros (`species`, `location`, `minAge`, `maxAge`, `status`).
- `POST|PUT|DELETE /pets` – CRUD para refugios autenticados y aprobados.
- `POST /pets/{id}/adoption-requests` – Crear solicitud (solo adoptantes).
- `GET /adoption-requests` – Listado contextual (adoptante o refugio).
- `PATCH /adoption-requests/{id}` – Actualizar estado (refugio propietario).
- `GET /admin/users` / `POST /admin/shelters/{id}/approve` – Moderacion basica para administradores.

Errores estandar devuelven `{ "error": string }` y codigos HTTP adecuados.

## Configuracion y ejecucion
1. Instala Go >= 1.21.
2. Desde `Backend/` instala dependencias: `go mod tidy` (ya ejecutado).
3. Ejecuta la API:
   ```bash
   go run ./cmd/server
   ```
4. Variables de entorno relevantes:
   ```bash
   PETMATCH_DB_PATH=petmatch.db
   PETMATCH_HTTP_PORT=8080
   PETMATCH_JWT_SECRET=change-me
   PETMATCH_ADMIN_EMAIL=admin@petmatch.local
   PETMATCH_ADMIN_PASSWORD=admin123
   ```

> La primera ejecucion crea automaticamente un admin con las credenciales configuradas.

## Pruebas rapidas
- `go build ./...` para validar compilacion.
- El proyecto usa SQLite embebido; cambia el DSN ajustando `PETMATCH_DB_PATH`.

## Siguientes pasos sugeridos
- Implementar paginacion/buscador avanzado.
- Agregar pruebas unitarias en services y handlers.
- Permitir carga de imagenes reales (hoy solo URL).
- Externalizar configuracion a archivos o variables segun entorno.
