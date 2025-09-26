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
- **Concurrencia**: 10, 50, 100 usuarios simultáneos
- **Duración**: Tests de 30 segundos a 2 minutos
- **Endpoints**: `/migrate`, `/balance`, `/health`

### **Stress Tests**
- **Límites máximos**: Hasta 500 usuarios concurrentes
- **Recovery testing**: Comportamiento post-sobrecarga
- **Resource monitoring**: CPU, memoria, conexiones

## 📊 Métricas Clave

- **Response Time**: Latencia promedio y percentiles (P95, P99)
- **Throughput**: Requests por segundo (RPS)
- **Error Rate**: Porcentaje de requests fallidos
- **Resource Usage**: CPU, memoria, conexiones de red

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
- Response time < 2 segundos
- Error rate < 1%
- Throughput estable
- Recursos dentro de límites

### **⚠️ Advertencia**
- Response time 2-5 segundos
- Error rate 1-5%
- Throughput variable
- Recursos altos pero manejables

### **❌ Falla**
- Response time > 5 segundos
- Error rate > 5%
- Throughput degradado
- Recursos agotados

### **Personalización**
- Ajustar `concurrency` según infraestructura
- Modificar `duration` según necesidades

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
