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

### Instalación local (comandos para MacOS)
```bash
# Abre una terminal y ejecuta los siguientes comandos

# Clonar el repositorio
git clone https://github.com/jorgelan8/API_Stori.git

# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori 

# Instalar dependencias
go mod tidy

# Configurar variables de entorno
#   Creamos el archivo .env a partir del template env.example
#     Indispensable si se quiere comprobar que se envia el email con el Summary Report en el endpoint /migrate
cp env.example .env

# Editar .env con tus configuraciones, usa nano o tu editor favorito
nano .env

# Ejecutar el API (server)
go run cmd/api/main.go

# **** Ahora ya puedes hacer request a la API ****
```

### 🧪 Testing API endpoints local (Comando para MacOS)
El server local esta configurado para usar el puerto 8080
```bash
# Abrir una terminal y ejecutar los siguientes comandos

# Probar health endpoint
curl -s http://localhost:8080/api/v1/health

# Probar root endpoint...
curl -s http://localhost:8080/

# 🧪 Probar migrate endpoint con archivo CSV, asegurate de colocar la ruta correcta del archivo a cargar, el repositoro del API contiene un arhivo de ejemplo para el exito de estas pruebas
curl -X POST http://localhost:8080/api/v1/migrate -F "csv_file=@examples/sample_transactions.csv"

#Puede crear un nuevo archivo, debe asegurarse que el formato del archivo sea el correcto

# 🧪 Probar balance endpoint, debe haber cargado un archivo en el endpoint /migrate
curl -s "http://localhost:8080/api/v1/users/1001/balance"

# Si carga un archivo diferente al de ejemplo, debe ajustar el "1001" al user_id que quiere probar
```

### Usar el API con Docker (requiere estar instalado Docker) (Descargalo [aquí][UrlDocker])
```bash
# Abre una terminal y ejecuta los siguientes comandos

# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori


# Configurar variables de entorno
#   Editar archivo docker.env.development
#     Indispensable si se quiere comprobar que se envia el email con el Summary Report en el endpoint /migrate
nano docker.env.development

# El archivo traer valores por defecto, pero vencen proximamente, se recomienda sus datos


# Construir y ejecutar contenedor
docker-compose up

# si lo prefieres O usar el script (dar permiso de ejecución chmod +x start.sh)
./start.sh

# **** Listo el ya puedes hacer request hacia el API en el Contenedor ****
```

### 🧪 Testing API endpoints en Docker
#### El contenedor esta configurado para usar el puerto 8081
```bash
# Una vez que el contenedor de Docker esta activo abra una terminal y ejecutar los siguientes comandos

# Probar health endpoint
curl -s http://localhost:8081/api/v1/health

# Probar root endpoint...
curl -s http://localhost:8081/

# 🧪 Probar migrate endpoint con archivo CSV, asegurate de colocar la ruta correcta del archivo a cargar, el repositoro del API contiene un arhivo de ejemplo para el exito de estas pruebas
curl -X POST http://localhost:8081/api/v1/migrate -F "csv_file=@examples/sample_transactions.csv"

#Puede crear un nuevo archivo, debe asegurarse que el formato del archivo sea el correcto

# 🧪 Probar balance endpoint, debe haber cargado un archivo en el endpoint /migrate
curl -s "http://localhost:8081/api/v1/users/1001/balance"

# Si carga un archivo diferente al ejemplo debe ajustar el 1001 al user_id que quiere probar
```

## Desarrollo

### 🧪 Testing
- Ejecutar pruebas unitarias
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

make test-unit
```
- Ejecutar pruebas integrales
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

make test-integration   # Pruebas de integración
```
- Ejecutar todas las pruebas
```bash
# Cambiar al directorio del repositorio clonado por default el directorio es API_Stori
cd API_Stori

make test-all
````

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



