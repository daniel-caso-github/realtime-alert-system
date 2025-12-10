# Realtime Alert System

### Structure Project

```
realtime-alerting-system/
├── cmd/                        # Puntos de entrada (main.go)
│   └── api/                    # Servidor API principal
│
├── internal/                   # Código privado (no importable externamente)
│   ├── domain/                 # 🟢 NÚCLEO - Entidades y reglas de negocio
│   │   ├── entity/             # Estructuras de dominio (Alert, User, etc.)
│   │   ├── repository/         # Interfaces de repositorios (ports)
│   │   └── service/            # Interfaces de servicios de dominio
│   │
│   ├── application/            # 🔵 CASOS DE USO - Lógica de aplicación
│   │   ├── dto/                # Data Transfer Objects
│   │   ├── usecase/            # Implementación de casos de uso
│   │   └── service/            # Servicios de aplicación
│   │
│   ├── infrastructure/         # 🟠 ADAPTADORES - Implementaciones externas
│   │   ├── config/             # Configuración (Viper)
│   │   ├── database/           # PostgreSQL, Redis
│   │   ├── messaging/          # Redis Streams, NATS
│   │   ├── notification/       # Slack, Email, SMS
│   │   └── logger/             # Logging estructurado
│   │
│   └── presentation/           # 🟣 INTERFAZ - HTTP, WebSocket
│       ├── http/               # Handlers REST API
│       │   ├── handler/
│       │   ├── middleware/
│       │   └── router/
│       └── websocket/          # WebSocket server
│
├── pkg/                        # Código público reutilizable
│   └── utils/                  # Utilidades compartidas
│
├── deployments/                # Configuraciones de despliegue
│   ├── docker/                 # Dockerfiles adicionales
│   └── kubernetes/             # Manifests K8s, Helm charts
│
├── terraform/                  # Infraestructura como código
│   ├── modules/                # Módulos reutilizables
│   └── environments/           # Dev, staging, prod
│
├── scripts/                    # Scripts de automatización
├── migrations/                 # Migraciones de base de datos
├── docs/                       # Documentación adicional
└── test/                       # Tests de integración/e2e
```