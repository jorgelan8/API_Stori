# API_Stori

Proyecto de API RESTful en Go - Estructura inicial

## 📁 Estructura del Proyecto

```
api-stori/
├── api/
│   ├── docs/           # Documentación de la API
│   └── swagger/        # Archivos de Swagger/OpenAPI
├── cmd/
│   └── api/            # Punto de entrada principal de la aplicación
├── internal/           # Código privado de la aplicación
│   ├── config/         # Configuración de la aplicación
│   ├── handlers/       # Manejadores HTTP (controladores)
│   ├── middleware/     # Middleware personalizado
│   ├── models/         # Modelos de datos y estructuras
│   ├── routes/         # Definición de rutas
│   └── services/       # Lógica de negocio
├── pkg/                # Código que puede ser reutilizado
│   ├── logger/         # Utilidades de logging
│   └── utils/          # Utilidades generales
├── scripts/            # Scripts de automatización
├── tests/              # Tests del proyecto
│   ├── integration/    # Tests de integración
│   └── unit/           # Tests unitarios
└── README.md           # Documentación del proyecto
```
