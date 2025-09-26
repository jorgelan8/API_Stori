# API Stori - Migration & Balance Service

API RESTful en Go para migración de transacciones y consulta de balance de usuarios.


## 🚀 Características

- **Migración de transacciones** desde archivos CSV
- **Consulta de balance** de usuarios con filtros de fecha
- **Reportes automáticos** por email después de la migración
- **Documentación OpenAPI** completa (Swagger UI)
- **Suite de pruebas** completa (unitarias, integración, carga, rendimiento)
- **Configuración flexible** mediante variables de entorno
- **Docker** para despliegue y desarrollo

## 📋 Endpoints de la API

### Health Check
- `GET /api/v1/health` - Estado de salud de la API

### Migración
- `POST /api/v1/migrate` - Subir y procesar archivo CSV de transacciones

### Balance
- `GET /api/v1/users/{user_id}/balance` - Obtener balance de usuario
  - Query params: `from_date`, `to_date` (opcionales)

### Documentación
- `GET /api/v1/docs` - Swagger UI interactivo
- `GET /api/v1/swagger.yaml` - Especificación OpenAPI en YAML
- `GET /api/v1/swagger.json` - Especificación OpenAPI en JSON

## 🛠️ Instalación y Uso

### Prerrequisitos
- Go 1.21+ (Descargalo [aqui][UrlGo])
- Docker (opcional) (Solo si se quiere probar localmente) (Descargalo [aquí][UrlDocker])

### Instalación local
```bash
# Clonar el repositorio
git clone https://github.com/jorgelan8/API_Stori.git
cd API_Stori 

# Instalar dependencias
go mod tidy

# Configurar variables de entorno
#   Indispensable si se quiere comprobar enviando un email
cp env.example .env
# Editar .env con tus configuraciones

# Ejecutar el API (server)
go run cmd/api/main.go
```

### Usar el API con Docker (requiere estar instalado Docker) (Descargalo [aquí][UrlDocker])
```bash
# Construir y ejecutar contenedor
docker-compose up

# O usar el script
./start.sh
```

## 🧪 Testing

```bash
# Ejecutar todas las pruebas
make test-all

# Pruebas específicas
make test-unit          # Pruebas unitarias
make test-integration   # Pruebas de integración
```

## 📧 Configuración de Email

```bash
# Configurar email interactivamente
./configure_email.sh

# Probar envío de email
./test_email_report.sh
```

## 🔧 Variables de Entorno

Ver `env.example` para todas las variables disponibles.

### Principales:
- `PORT` - Puerto del servidor (default: 8080)
- `APP_ENV` - Entorno (development/production)
- `SMTP_HOST` - Servidor SMTP para reportes
- `SMTP_USER` - Usuario SMTP
- `SMTP_PASS` - Contraseña SMTP
- `TO_EMAILS` - Emails destino para reportes


## 🐳 Docker

### Desarrollo
```bash
docker-compose -f docker-compose.dev.yml up
```

### Producción
```bash
docker-compose up
```

## 📚 Documentación Adicional

- [Documentación de Pruebas](tests/README.md)
- [Documentación de Endpoints](api/docs/)
- [Documentación Técnica](docs/)

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.


[UrlGo]:https://go.dev/doc/install "Golang"
[UrlDocker]:https://www.docker.com/products/docker-desktop/ "Docker"



