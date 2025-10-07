# PetMatch Frontend (Angular)

## Resumen
- **Framework**: Angular 17 (standalone components) con `HttpClient` y formularios reactivos.
- **Objetivo**: Permitir registro/login, exploracion de mascotas, solicitudes de adopcion y paneles basicos para refugios y administradores.
- **Diseño**: Layout responsivo sencillo, componentes reutilizables (`PetCard`) y servicios centralizados en `core/`.

## Estructura clave (`src/app`)
- `core/models`: Tipos compartidos (`User`, `Pet`, `AdoptionRequest`, `Auth`).
- `core/services`: `api`, `auth`, `pet`, `adoption`, `admin` encapsulan llamadas al backend y cabeceras JWT.
- `core/guards/auth.guard.ts`: Protege rutas y valida roles (`adopter`, `shelter`, `admin`).
- `shared/components/pet-card`: Tarjeta reutilizable para el listado.
- `features/...`: Secciones por dominio
  - `pets/pages/home`: landing + filtros dinamicos, consumo del catalogo.
  - `pets/pages/pet-detail`: ficha + formulario de solicitud.
  - `auth/login` y `auth/register`: formularios reactivos con validaciones.
  - `adoption/pages/requests`: listado contextual (adoptante/refugio) con acciones.
  - `admin/pages/users`: panel para aprobar refugios.

## Integracion API
- Base URL inyectable via token `API_BASE_URL` (por defecto `http://localhost:8080/api/v1`).
- `AuthService` maneja sesion (localStorage) y headers `Authorization`.
- Tras inicio de sesion se refresca el usuario con `/auth/me` para sincronizar roles.

## Rutas
| Ruta | Rol | Componente |
| --- | --- | --- |
| `/` | Publico | HomeComponent |
| `/pets/:id` | Publico | PetDetailComponent |
| `/auth/login` | Publico | LoginComponent |
| `/auth/register` | Publico | RegisterComponent |
| `/adoption-requests` | Adoptante/Refugio | RequestsComponent |
| `/admin/users` | Admin | UsersComponent |

## Scripts npm
```bash
npm install        # instala dependencias
npm start          # ng serve --open
npm run build      # build de produccion (genera /dist)
```

## Variables utiles
Ajusta el token en `core/config/api.tokens.ts` o redefine el provider en `main.ts` para apuntar a otro backend.

## Mejoras futuras
- Internacionalizacion e inclusion de un estado global.
- Estados de carga/toasts centralizados.
- Integrar subida real de imagenes y perfiles ampliados.
- Añadir pruebas unitarias (Karma/Jest) para servicios y componentes.
