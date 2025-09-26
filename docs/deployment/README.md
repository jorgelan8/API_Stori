# 🚀 Despliegue y Operaciones - API Stori

Este directorio contiene guías de despliegue, configuración de entornos y operaciones del proyecto API Stori.

## 📋 Guías de Despliegue

### **Entornos**
- [ ] Desarrollo Local
- [ ] Entorno de Testing
- [ ] Entorno de Producción

### **Configuración**
- [ ] Variables de Entorno
- [ ] Configuración de Base de Datos
- [ ] Configuración de Email

### **Docker**
- [ ] Docker Compose
- [ ] Dockerfile
- [ ] Volúmenes y Redes

## 🐳 Docker

### **Comandos Básicos**
```bash
# Desarrollo
docker-compose up -d

# Producción
docker-compose -f docker-compose.prod.yml up -d

# Logs
docker-compose logs -f api

# Limpiar
docker-compose down --volumes --remove-orphans
```

### **Configuración de Entornos**

#### **Desarrollo**
```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=mock
    volumes:
      - .:/app
```

#### **Producción**
```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=mysql
      - DB_HOST=mysql
    depends_on:
      - mysql
```

## 🔧 Variables de Entorno

### **Desarrollo**
```bash
# .env.development
APP_ENV=development
PORT=8080
DB_TYPE=mock
```

### **Producción**
```bash
# .env.production
APP_ENV=production
PORT=8080
DB_TYPE=mysql
DB_HOST=mysql
DB_PASSWORD=secure_password
```

## 📊 Monitoreo

### **Health Checks**
- `GET /api/v1/health` - Estado de la API
- `GET /api/v1/metrics` - Métricas de Prometheus

### **Logs**
- **Desarrollo**: Console logs
- **Producción**: Structured logs (JSON)

## 🔄 CI/CD

### **Pipeline de Despliegue**
1. **Build**: Compilar aplicación
2. **Test**: Ejecutar tests
3. **Security**: Análisis de seguridad
4. **Deploy**: Desplegar a entorno
5. **Verify**: Verificar despliegue

## 📚 Convención de Nombrado

Los archivos de despliegue siguen el formato:
```
YYYY-MM-DD_deployment_component_description.md
```

**Ejemplos:**
- `2025-09-26_docker_setup.md`
- `2025-09-26_environment_config.md`
- `2025-09-26_monitoring_setup.md`

## 🎯 Próximos Documentos

- [ ] **Configuración de Docker** - Setup completo
- [ ] **Variables de Entorno** - Guía de configuración
- [ ] **Monitoreo** - Setup de métricas y alertas
- [ ] **CI/CD Pipeline** - Automatización de despliegue

---

**📅 Última actualización**: 2025-09-26  
**👨‍💻 Mantenido por**: Equipo de Desarrollo API Stori
