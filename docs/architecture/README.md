# 🏗️ Arquitectura - API Stori

Este directorio contiene documentación de la arquitectura del sistema, diseño de API y decisiones técnicas.

## 📋 Documentos de Arquitectura

### **Sistema**
- [ ] Arquitectura General del Sistema
- [ ] Diagramas de Componentes
- [ ] Flujo de Datos

### **API Design**
- [ ] Diseño de Endpoints
- [ ] Especificación OpenAPI
- [ ] Patrones de Request/Response

### **Base de Datos**
- [ ] Diseño de Esquema
- [ ] Patrones de Acceso a Datos
- [ ] Estrategias de Migración

## 🎯 Principios Arquitectónicos

### **1. Separación de Responsabilidades**
- **Handlers**: Manejo de HTTP requests
- **Services**: Lógica de negocio
- **Repositories**: Acceso a datos
- **Models**: Entidades de dominio

### **2. Inyección de Dependencias**
- Interfaces para servicios
- Factory pattern para repositorios
- Configuración centralizada

### **3. Testing**
- Unit tests para servicios
- Integration tests para endpoints
- Load tests para rendimiento

## 📊 Diagramas de Arquitectura

### **Arquitectura Actual**
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Handlers      │    │   Services       │    │   Repositories  │
│                 │    │                  │    │                 │
│ MigrationHandler│───▶│ MigrationService │───▶│ MockDatabase    │
│ BalanceHandler  │    │ UsersService     │    │ (Future: MySQL) │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### **Arquitectura Futura (con Abstracción)**
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Handlers      │    │   Services       │    │   Repositories  │
│                 │    │                  │    │                 │
│ MigrationHandler│───▶│ MigrationService │───▶│ TransactionRepo │
│ BalanceHandler  │    │ UsersService     │    │ UserRepo        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                       │
                                               ┌───────┴───────┐
                                               │   Factory     │
                                               │ Mock | MySQL  │
                                               └───────────────┘
```

## 🔧 Decisiones Técnicas

### **Lenguaje y Framework**
- **Go 1.21+**: Performance y concurrencia
- **Gorilla Mux**: Routing HTTP
- **Standard Library**: Minimizar dependencias

### **Base de Datos**
- **Mock**: Para desarrollo y testing
- **MySQL**: Para producción
- **Abstracción**: Repository pattern

### **Testing**
- **Unit Tests**: Servicios individuales
- **Integration Tests**: Endpoints completos
- **Load Tests**: Rendimiento bajo carga
- **Performance Tests**: Métricas detalladas

## 📚 Convención de Nombrado

Los archivos de arquitectura siguen el formato:
```
YYYY-MM-DD_architecture_component_description.md
```

**Ejemplos:**
- `2025-09-26_system_architecture.md`
- `2025-09-26_api_design.md`
- `2025-09-26_database_design.md`

## 🎯 Próximos Documentos

- [ ] **Arquitectura del Sistema** - Diagramas y componentes
- [ ] **Diseño de API** - Endpoints y especificaciones
- [ ] **Diseño de Base de Datos** - Esquemas y relaciones
- [ ] **Patrones de Código** - Convenciones y estándares

---

**📅 Última actualización**: 2025-09-26  
**👨‍💻 Mantenido por**: Equipo de Desarrollo API Stori
