# Casa Atreides — El Linaje que Gobierna los Datos

> *"El poder de gobernar está en los pequeños detalles."*
> — Thufir Hawat, Mentat de la Casa Atreides

## La Casa Atreides en el Universo Dune

La Casa Atreides es una de las Grandes Casas del Landsraad, conocida por su
honor, lealtad y excelencia en el gobierno. Originarios de Caladan, un mundo
oceánico, los Atreides gobernaron con justicia antes de ser enviados a Arrakis
— el planeta desierto — donde su destino se entrelazó con la especia melange.

### Valores de la Casa Atreides

| Valor | Significado en Dune | Significado en Atreides API |
|---|---|---|
| **Honor** | Lealtad inquebrantable a la palabra dada | Integridad de los datos, transacciones ACID |
| **Estrategia** | El Mentat calcula cada movimiento | Clean Architecture, cada capa tiene un propósito |
| **Adaptación** | De Caladan a Arrakis, sobrevivir y prosperar | De librería-cafetería a cualquier dominio de negocio |
| **Gobierno** | Administrar recursos de un planeta | Gestionar inventario, productos, empresas |
| **Legado** | El linaje perdura generaciones | El código es mantenible, extensible, duradero |

### El Escudo de la Casa

El emblema original de la Casa Atreides es un **halcón rojo** sobre fondo verde.
En nuestro contexto digital, el halcón representa la API que vela sobre los datos:

```
    ╔══════════════════════════════════╗
    ║     CASA ATREIDES — API REST    ║
    ║   "Gobernar con honor y orden"   ║
    ╚══════════════════════════════════╝
```

## Paralelismos con la Arquitectura de la API

### El Mentat (Capa de Configuración)

> *"Un Mentat debe ser objetivo, desapasionado y preciso."*
> — Thufir Hawat

- `internal/config/` es nuestro Mentat: calcula los parámetros exactos para que
  el sistema funcione (Postgres DSN, JWT secret). Sin emoción, solo datos.

### El Maester de la Espada (Middleware)

> *"La habilidad de un guerrero se mide por su control."*
> — Gurney Halleck

- El **AuthMiddleware** es nuestro Gurney Halleck: verifica cada petición como
  un guerrero revisa sus armas. Sin token válido, no hay paso.
- El **RecoveryMiddleware** es el escudo que protege contra ataques inesperados
  (panics), manteniendo la dignidad de la Casa.

### La Base de Datos (Arrakis)

- PostgreSQL es **Arrakis**: el lugar donde se extrae y almacena el recurso más
  valioso (los datos). JSONB son los depósitos de especia; las migraciones son
  los túneles que los fremen construyen para acceder a nuevas vetas.

### Los Handlers (Consejo de Guerra)

Cada handler es un **Consejo de Guerra** donde el Duque (el endpoint) escucha
las peticiones del Landsraad (el cliente HTTP), consulta a sus mentats (usecases)
y emite órdenes (respuestas JSON).

## El Juramento de la API

> *"No temeré al error 500, porque el RecoveryMiddleware lo capturará.*
> *No temeré al token inválido, porque el AuthMiddleware lo rechazará.*
> *No temeré a los datos corruptos, porque las transacciones son ACID.*
> *Porque soy Atreides, y el orden gobierna sobre el caos."*
>
> — Juramento del Desarrollador Atreides

## Diagrama de la Casa

```
                    ┌──────────────────┐
                    │   CASA ATREIDES  │
                    │   (API REST)     │
                    └────────┬─────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
    ┌─────────▼──────┐ ┌────▼─────┐ ┌──────▼────────┐
    │   Caladan      │ │ Arrakis  │ │  El Imperio   │
    │  (Dominio)     │ │  (DB)    │ │  (Handlers)   │
    │ 18 entidades   │ │PostgreSQL│ │ 19 handlers   │
    │ puras sin deps │ │  JSONB   │ │ HTTP + CORS   │
    └────────────────┘ └──────────┘ └───────────────┘
```

Atreides no es solo una API. Es el gobierno digital de tu negocio,
con el honor y la estrategia de la más noble Casa del Landsraad.
