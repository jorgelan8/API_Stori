# 📝 Convención de Nombrado - API Stori

Este documento establece las convenciones de nombrado para archivos de documentación en el proyecto API Stori.

## 🎯 Formato General

### **Estructura: `YYYY-MM-DD_category_description.md`**

**Componentes:**
- **YYYY-MM-DD**: Fecha en formato ISO 8601
- **category**: Categoría del documento (opcional)
- **description**: Descripción clara del contenido
- **.md**: Extensión Markdown

## 📁 Categorías de Documentos

### **📋 Plans (Planes de Desarrollo)**
```
YYYY-MM-DD_plan_name.md
```

**Ejemplos:**
- `2025-09-26_database_abstraction_plan.md`
- `2025-09-30_security_improvements_plan.md`
- `2025-10-01_performance_optimization_plan.md`

### **🔧 Improvements (Mejoras Implementadas)**
```
YYYY-MM-DD_improvement_category_description.md
```

**Ejemplos:**
- `2025-09-26_best_practices_analysis.md`
- `2025-09-28_test_centralization.md`
- `2025-10-02_error_handling_improvements.md`

### **🏗️ Architecture (Arquitectura)**
```
YYYY-MM-DD_architecture_component_description.md
```

**Ejemplos:**
- `2025-09-26_system_architecture.md`
- `2025-09-26_api_design.md`
- `2025-09-26_database_design.md`

### **🚀 Deployment (Despliegue)**
```
YYYY-MM-DD_deployment_component_description.md
```

**Ejemplos:**
- `2025-09-26_docker_setup.md`
- `2025-09-26_environment_config.md`
- `2025-09-26_monitoring_setup.md`

### **🔧 Troubleshooting (Solución de Problemas)**
```
YYYY-MM-DD_troubleshooting_category_description.md
```

**Ejemplos:**
- `2025-09-26_common_issues.md`
- `2025-09-26_debugging_guide.md`
- `2025-09-26_performance_issues.md`

## 📋 Reglas de Nombrado

### **1. Fecha (YYYY-MM-DD)**
- ✅ **Correcto**: `2025-09-26`
- ❌ **Incorrecto**: `26-09-2025`, `Sept-26-2025`

### **2. Separadores**
- ✅ **Correcto**: `_` (guión bajo)
- ❌ **Incorrecto**: `-` (guión), ` ` (espacio)

### **3. Descripción**
- ✅ **Correcto**: `database_abstraction_plan`
- ❌ **Incorrecto**: `DatabaseAbstractionPlan`, `database-abstraction-plan`

### **4. Extensión**
- ✅ **Correcto**: `.md`
- ❌ **Incorrecto**: `.txt`, `.doc`, sin extensión

## 🎯 Ejemplos de Nombrado Correcto

### **Planes**
```
2025-09-26_database_abstraction_plan.md
2025-09-30_security_improvements_plan.md
2025-10-05_api_versioning_plan.md
```

### **Mejoras**
```
2025-09-26_best_practices_analysis.md
2025-09-28_test_centralization.md
2025-10-01_logging_improvements.md
```

### **Arquitectura**
```
2025-09-26_system_architecture.md
2025-09-26_api_design.md
2025-09-26_database_design.md
```

### **Despliegue**
```
2025-09-26_docker_setup.md
2025-09-26_environment_config.md
2025-09-26_monitoring_setup.md
```

### **Troubleshooting**
```
2025-09-26_common_issues.md
2025-09-26_debugging_guide.md
2025-09-26_performance_issues.md
```

## 🔄 Proceso de Creación de Documentos

### **1. Identificar Categoría**
- ¿Es un plan? → `docs/plans/`
- ¿Es una mejora? → `docs/improvements/`
- ¿Es arquitectura? → `docs/architecture/`
- ¿Es despliegue? → `docs/deployment/`
- ¿Es troubleshooting? → `docs/troubleshooting/`

### **2. Generar Nombre**
- Obtener fecha actual: `2025-09-26`
- Identificar categoría: `plan`, `improvement`, etc.
- Describir contenido: `database_abstraction`
- Formar nombre: `2025-09-26_database_abstraction_plan.md`

### **3. Crear Archivo**
```bash
# Ejemplo
touch docs/plans/2025-09-26_database_abstraction_plan.md
```

## 📊 Ventajas de esta Convención

### **✅ Orden Cronológico**
- Los archivos se ordenan automáticamente por fecha
- Fácil ver la evolución temporal del proyecto

### **✅ Descriptivo**
- El nombre explica claramente el contenido
- Fácil identificar el propósito del documento

### **✅ Escalable**
- Fácil agregar nuevas categorías
- Fácil crear versiones o actualizaciones

### **✅ Estándar**
- Formato ISO 8601 para fechas
- Convención reconocible y profesional

## 🎯 Casos Especiales

### **Actualizaciones de Documentos**
```
# Documento original
2025-09-26_database_abstraction_plan.md

# Actualización
2025-09-28_database_abstraction_plan_v2.md
```

### **Documentos por Versión**
```
2025-09-26_api_v1_design.md
2025-10-01_api_v2_design.md
```

### **Documentos por Entorno**
```
2025-09-26_development_environment_setup.md
2025-09-26_production_environment_setup.md
```

---

**📅 Última actualización**: 2025-09-26  
**👨‍💻 Mantenido por**: Equipo de Desarrollo API Stori  
**🎯 Objetivo**: Establecer convención clara y consistente para nombrado de documentos
