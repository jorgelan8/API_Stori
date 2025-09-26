# 🔧 Solución de Problemas - API Stori

Este directorio contiene guías para diagnosticar y resolver problemas comunes del proyecto API Stori.

## 🚨 Problemas Comunes

### **Errores de Base de Datos**
- [ ] Error de conexión a MySQL
- [ ] Timeout de queries
- [ ] Problemas de migración

### **Errores de API**
- [ ] 404 Not Found
- [ ] 500 Internal Server Error
- [ ] Timeout de requests

### **Errores de Testing**
- [ ] Tests fallando
- [ ] Problemas de concurrencia
- [ ] Timeout en load tests

## 🔍 Diagnóstico

### **Logs de Aplicación**
```bash
# Ver logs en tiempo real
docker-compose logs -f api

# Ver logs específicos
docker-compose logs api | grep ERROR
```

### **Health Checks**
```bash
# Verificar estado de la API
curl http://localhost:8080/api/v1/health

# Verificar métricas
curl http://localhost:8080/api/v1/metrics
```

### **Base de Datos**
```bash
# Conectar a MySQL
docker-compose exec mysql mysql -u root -p

# Verificar tablas
SHOW TABLES;

# Verificar datos
SELECT COUNT(*) FROM transactions;
```

## 🛠️ Soluciones

### **Problema: API no responde**
```bash
# Verificar contenedores
docker-compose ps

# Reiniciar API
docker-compose restart api

# Ver logs de error
docker-compose logs api
```

### **Problema: Error de conexión a BD**
```bash
# Verificar MySQL
docker-compose exec mysql mysql -u root -p

# Verificar variables de entorno
docker-compose exec api env | grep DB_

# Reiniciar servicios
docker-compose restart
```

### **Problema: Tests fallando**
```bash
# Limpiar cache de Go
go clean -cache

# Ejecutar tests con verbose
go test -v ./tests/...

# Ejecutar tests específicos
go test -v ./tests/unit/...
```

## 📊 Herramientas de Debugging

### **Go Debugging**
```bash
# Profiling de CPU
go tool pprof http://localhost:8080/debug/pprof/profile

# Profiling de memoria
go tool pprof http://localhost:8080/debug/pprof/heap

# Trace de goroutines
go tool trace trace.out
```

### **Docker Debugging**
```bash
# Inspeccionar contenedor
docker inspect api-stori_api_1

# Ejecutar shell en contenedor
docker-compose exec api sh

# Ver uso de recursos
docker stats
```

## 📚 Convención de Nombrado

Los archivos de troubleshooting siguen el formato:
```
YYYY-MM-DD_troubleshooting_category_description.md
```

**Ejemplos:**
- `2025-09-26_common_issues.md`
- `2025-09-26_debugging_guide.md`
- `2025-09-26_performance_issues.md`

## 🎯 Próximos Documentos

- [ ] **Problemas Comunes** - Lista de errores frecuentes
- [ ] **Guía de Debugging** - Herramientas y técnicas
- [ ] **Problemas de Performance** - Optimización y tuning
- [ ] **Problemas de Seguridad** - Vulnerabilidades y fixes

---

**📅 Última actualización**: 2025-09-26  
**👨‍💻 Mantenido por**: Equipo de Desarrollo API Stori
