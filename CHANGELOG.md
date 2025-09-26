# Historial de Cambios (CHANGELOG)

Todas las versiones notables de este proyecto se documentarán en este archivo.

El formato se basa en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), y este proyecto adhiere al [Versionado Semántico (SemVer)](https://semver.org/lang/es/).

## [1.1.3] - 2025-09-25

### Añadido (Added)
- Se agregaron **pruebas de rendimiento** con benchmarking de endpoints críticos.
- Se implementó análisis de latencia y rendimiento para su optimización.
- Se añadió documentación completa de pruebas de rendimiento.

## [1.1.2] - 2025-09-25

### Añadido (Added)
- Se implementaron **pruebas de carga** con testing de concurrencia con una variabilidad de goroutine (10-1000).
- Se agregaron métricas de rendimiento bajo carga para endpoints `/migrate` y `/balance`.
- Se incluyó testing de balance con y sin rango de fechas bajo carga.
- No manda email, para evitar saturación del servidor
- Se añadió documentación de load testing con configuración específica por test.

## [1.1.0] - 2024-09-25

### Añadido (Added)
- Se agregó el **resumen de migración** que es enviado asincronamente por email después de procesar el CSV
- Se implementó funcionalidad de generación de reportes de migración.
- Se añadió formato de reporte estructurado con métricas de migración.
- Se añadieron variable de entorno

## [1.0.3] - 2025-09-24

### Añadido (Added)
- Se implementaron **pruebas de integración** con testing end-to-end de endpoints.
- Se agregó validación de flujos completos de la API.

### Cambiado (Changed)
- Se refactorizó la estructura de tests para mejorar la mantenibilidad.
- Se centralizó la configuración de tests para evitar duplicación de código.

## [1.0.2] - 2024-09-24

### Añadido (Added)
- Se agregó **centralización de tests** con funciones reutilizables en `tests/test_utils/`.
- Se implementó `SetupTestServer()` para configuración centralizada de servidor de pruebas.
- Se añadió `GenerateTestCSV()` para generación de datos de prueba consistentes.
- Se incluyó `CreateMultipartFormDataPerRequest()` para testing thread-safe de multipart data.
- Se agregó `MigrateTestData()` para migración de datos de prueba.

## [1.0.1] - 2025-09-24

### Añadido (Added)
- Se implementaron **pruebas unitarias** con cobertura de código básica.
- Se agregó mocking de dependencias para testing aislado.
- Se incluyó estructura de testing unitario para todos los componentes.

## [1.0.0] - 2024-01-24

### Añadido (Added)
- Primera versión compilada de la API Stori.
- Se implementó estructura básica del proyecto con endpoints principales.
- Se agregaron endpoints básicos: `/health`, `/migrate`, `/balance`.
- Se incluyó configuración básica de base de datos (Mock) y servidor HTTP.
- Se añadió documentación inicial del proyecto.

## [0.0.0] - 2024-01-23

### Añadido (Added)
- Se analizo el requerimiento.
- Se tomo minicurso de GO.
- Se determino la estructura para el proyecto.

