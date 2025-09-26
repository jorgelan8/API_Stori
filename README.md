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

##### [Ver Historila de cambios](CHANGELOG.md)

## 🛠️ Instalación y Uso

### Prerrequisitos
- Go 1.21+ (Descargalo [aqui][UrlGo])
- Docker (opcional) (Solo si se quiere probar localmente) (Descargalo [aquí][UrlDocker])

### Instalación local (Comandos para MacOS)
<details>
  <summary> Clic to details </summary>
#### Abre una terminal y ejecuta los siguientes comandos

- Clonar el repositorio
```bash
git clone https://github.com/jorgelan8/API_Stori.git
```

- Instalar dependencias
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori 

go mod tidy
```

- Configurar variables de entorno
```bash
# Debe crearse el archivo .env a partir del template env.example
#   Indispensable si se quiere comprobar que se envia el email con el Summary Report en el endpoint /migrate
cp env.example .env

# Editar .env con tus configuraciones, usa nano o tu editor favorito
nano .env
```
</details>

- Ejecutar el API (server)
```bash
go run cmd/api/main.go
```
*** Ahora ya puedes hacer request a la API ***

### 🧪 Testing API endpoints local
#### El server local esta configurado para usar el puerto 8080
##### Abrir una terminal y ejecutar los siguientes comandos

- Probar health endpoint
```bash
curl -s http://localhost:8080/api/v1/health
```

- Probar root endpoint...
```bash
curl -s http://localhost:8080/
```
- Probar migrate endpoint con archivo CSV ([Ver Doc][EPmigrate])
```bash
# Asegurate de colocar la ruta correcta del archivo a cargar
# el repositoro del API contiene un archivo de ejemplo para el exito de estas pruebas

curl -X POST http://localhost:8080/api/v1/migrate -F "csv_file=@examples/sample_transactions.csv"

# Puede crear un nuevo archivo, debe asegurarse que el formato del archivo sea el correcto
```

- Probar balance endpoint ([Ver Doc][EPBalance]), primero debió haber cargado un archivo en el endpoint migrate
```bash
curl -s "http://localhost:8080/api/v1/users/1001/balance"

# Si carga un archivo diferente al del ejemplo, debe ajustar el "1001" al user_id que quiere probar
```

### Usar el API con Docker (Comandos para MacOS)
#### (requiere estar instalado Docker) (Descargalo [aquí][UrlDocker]) 

##### Abre una terminal y ejecuta los siguientes comandos
- Configurar variables de entorno
  - Editar el archivo docker.env.development, con sus datos
  - Indispensable para comprobar el envío de email con el Summary Report en el endpoint migrate
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

nano docker.env.development

# El archivo traer valores por defecto, pero vencen proximamente
```

- Construir y ejecutar contenedor
```bash
docker-compose up
```

- si lo prefiere puede usar este script (dar permiso de ejecución chmod +x start.sh)
```bash
./start.sh
```
*** Listo el ya puedes hacer request hacia el API en el Contenedor ***


### 🧪 Testing API endpoints en Docker
#### El contenedor esta configurado para usar el puerto 8081
##### Abrir una terminal y ejecutar los siguientes comandos
- Probar health endpoint
```bash
curl -s http://localhost:8081/api/v1/health
```

- Probar root endpoint
``` bash
curl -s http://localhost:8081/
```

- Probar migrate endpoint con archivo CSV ([Ver Doc][EPmigrate])
```bash
# Asegurate de colocar la ruta correcta del archivo a cargar
# el repositoro del API contiene un arhivo de ejemplo para el exito de estas pruebas

curl -X POST http://localhost:8081/api/v1/migrate -F "csv_file=@examples/sample_transactions.csv"

#Puede crear un nuevo archivo, debe asegurarse que el formato del archivo sea el correcto
```

- Probar balance endpoint ([Ver Doc][EPBalance]), primero debió haber cargado un archivo en el endpoint migrate
```bash
curl -s "http://localhost:8081/api/v1/users/1001/balance"

# Si carga un archivo diferente al de ejemplo, debe ajustar el "1001" al user_id que quiere probar
```



## Desarrollo

### 🧪 Testing
- Ejecutar pruebas unitarias
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

go test -v ./internal/services/... ./internal/handlers/...
```
- Ejecutar pruebas integrales
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

go test -v ./tests/integration/...
```
- Ejecutar todas las pruebas
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

# dar permiso de ejecución chmod +x run_test.sh
./run_tests.sh
```

## Testing sofisticados

- ⚡ Tests de Carga(rendimiento bajo stress)

- **`load_test.go`** - Tests de carga con múltiples requests concurrentes ([Ver Doc][LoadTest])


### Ejecutar tests de carga:
```bash
go test -v ./tests/load/...
```

- 📊 Tests de Rendimiento

- **`performance_test.go`** - Tests de rendimiento y benchmarks ([Ver Doc][PerfTest])

### Ejecutar tests de rendimiento:
```bash
go test -v ./tests/performance/...
```

## 🔧 Variables de Entorno

Ver `env.example` para todas las variables disponibles.

### Principales:
- `APP_ENV` - Entorno (development/production)
- `PORT` - Puerto del servidor (default: 8080)
- `SMTP_HOST` - Servidor SMTP para reportes
- `SMTP_USER` - Usuario SMTP
- `SMTP_PASS` - Contraseña SMTP
- `TO_EMAILS` - Emails destino para reportes

## 📚 Documentación Técnica

### 🏗️ Arquitectura y Diseño
- [Arquitectura del Sistema](docs/architecture/) - Diseño y componentes del sistema
- [Planes de Desarrollo](docs/plans/) - Roadmaps y estrategias
- [Mejoras Implementadas](docs/improvements/) - Mejoras y optimizaciones

### 🚀 Operaciones
- [Guía de Despliegue](docs/deployment/) - Docker, entornos y configuración
- [Solución de Problemas](docs/troubleshooting/) - Debugging y troubleshooting

### 📋 API Endpoints
- [Documentación Endpoint /migrate][EPmigrate]
- [Documentación Endpoint /users/{user_id}/balance][EPBalance]
- [Documentación pruebas de stress][LoadTest]
- [Documentacion pruebas de performance][PerfTest]


## 🎯 Próximos Pasos

- [X] **Pruebas de Stress**: Swagger/OpenAPI para documentación interactiva
- [X] **Pruebas de Performance**: Swagger/OpenAPI para documentación interactiva
- [ ] **Base de datos**: Guardado permanente de las transacciones
- [ ] **CI/CD**: Integración continua con pruebas automáticas
- [ ] **Usar Secrets**: Cambiar las variables de entorno por secrets
- [ ] **Monitoreo**: Mejorar de logging, monitoreo y control de errores
- [ ] **Reportes de comportamiento de usuarios**: Promedio de saldo al mes, identificacion de temporada baja
- [ ] **Soportar conexiones fuera del dominio**: Middleware para CORS

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
[EPmigrate]:api/docs/migration_endpoints.md "Endpoint /migrate"
[EPBalance]:api/docs/balance_endpoints.md "Endpoint users/{user_id}/balance"
[LoadTest]:tests/load/load_test.md "Load Test"
[PerfTest]:tests/performance/performance_test.md "Performance Test"



