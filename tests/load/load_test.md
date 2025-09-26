# Load Testing - API Stori

## 🎯 ¿Qué son las Pruebas de Carga?

Las pruebas de carga evalúan el comportamiento de la API bajo diferentes niveles de estrés y concurrencia, simulando condiciones reales de uso.

## 🚀 Beneficios

### **1. Detección Temprana de Problemas**
- **Bottlenecks de rendimiento** antes del despliegue
- **Memory leaks** y problemas de recursos
- **Timeouts** y errores bajo carga

### **2. Validación de Escalabilidad**
- **Límites de concurrencia** reales del sistema
- **Capacidad máxima** de procesamiento
- **Degradación gradual** vs fallos catastróficos

### **3. Optimización de Recursos**
- **Configuración óptima** de servidores
- **Ajuste de timeouts** y límites
- **Planificación de infraestructura**

### **4. Confianza en Producción**
- **SLA compliance** bajo carga real
- **Experiencia de usuario** consistente
- **Prevención de downtime** por sobrecarga

## 🧪 Tests Incluidos

### **Load Tests**
- **Concurrencia**: 10-15 goroutines simultáneas
- **Requests por goroutine**: 25-250 requests
- **Total requests**: 250-3,750 requests por test
- **Endpoints**: `/migrate`, `/balance` (con y sin date range)

### **Concurrency Tests**
- **Migration load**: 10 goroutines × 25 requests = 250 total
- **Balance load**: 10 goroutines × 250 requests = 2,500 total  
- **Balance with date range**: 15 goroutines × 250 requests = 3,750 total

## 📊 Métricas Clave

- **Duration**: Tiempo total de ejecución del test
- **Throughput**: Requests por segundo (RPS) calculado
- **Success Count**: Número de requests exitosos
- **Error Count**: Número de requests fallidos (debe ser 0)
- **Timeout**: 30 segundos máximo por test

## 🚀 Ejecución

```bash
# Ejecutar load tests
go test -v ./tests/load/...

# Con coverage
go test -v ./tests/load/... -cover

# Con timeout personalizado
go test -v ./tests/load/... -timeout 10m
```

## 📈 Interpretación de Resultados

### **✅ Éxito**
- Error count = 0 (todos los requests exitosos)
- Test completa en < 30 segundos
- Throughput > 10 RPS
- Sin timeouts

### **⚠️ Advertencia**
- Error count < 5% del total
- Test completa en < 30 segundos
- Throughput > 5 RPS
- Algunos timeouts ocasionales

### **❌ Falla**
- Error count > 5% del total
- Test timeout después de 30 segundos
- Throughput < 5 RPS
- Muchos timeouts

### **Configuración del Código**
- **Migration test**: `concurrency = 10`, `requestsPerGoroutine = 25`
- **Balance test**: `concurrency = 10`, `requestsPerGoroutine = 250`
- **Balance with date range**: `concurrency = 15`, `requestsPerGoroutine = 250`
- **Timeout**: 30 segundos máximo por test

## 📋 Mejores Prácticas

1. **Ejecutar regularmente** en CI/CD
2. **Baseline establecido** antes de cambios
3. **Monitoreo continuo** en producción
4. **Documentar resultados** y tendencias
5. **Actualizar tests** con nuevos endpoints

## 🎯 Próximos Pasos

- [ ] **Automatizacion de pruebas de carga** en CI/CD
- [ ] **Performance regression** detection
- [ ] **Monitoreo en Tiempo-Real** en la integracion
- [ ] **Load testing** contra base de datos real


---

**💡 Tip**: Ejecuta load tests después de cada cambio significativo para mantener la confianza en el rendimiento de la API.
