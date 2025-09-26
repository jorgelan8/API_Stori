# Migration Service - API Endpoints

Este documento describe los endpoints disponibles para el servicio de migración de transacciones.

## 🚀 Endpoints Disponibles

### 1. POST /api/v1/migrate
**Descripción**: Procesa un archivo CSV con transacciones y las almacena en la base de datos.

**Request**:
- **Method**: POST
- **Content-Type**: multipart/form-data
- **Body**: Archivo CSV con las columnas: `id`, `user_id`, `amount`, `datetime`

**Ejemplo de uso con curl**:
```bash
curl -X POST http://localhost:8080/api/v1/migrate \
  -F "csv_file=@sample_transactions.csv"
```

**Response**:
```
200 OK
```

## 📁 Formato del Archivo CSV

El archivo CSV debe tener las siguientes columnas en el orden especificado:

```csv
id,user_id,amount,datetime
1,1001,150.50,2024-01-15 10:30:00
2,1001,-75.25,2024-01-15 14:45:00
```

### Columnas:
- **id** (int): Identificador único de la transacción
- **user_id** (int): ID del usuario propietario
- **amount** (float): Monto de la transacción (puede ser positivo o negativo)
- **datetime** (string): Fecha y hora en formato "YYYY-MM-DDTHH:MM:SSZ"

### Formatos de Fecha Soportados:
- `2006-01-02 15:04:05` (formato estándar)
- `2006-01-02T15:04:05` (formato ISO)
- `2006-01-02` (solo fecha)


## 📊 Características

- ✅ **Procesamiento de CSV** con validación de estructura
- ✅ **Almacenamiento en memoria** (mock de base de datos)
- ✅ **Manejo de errores** detallado por línea
- ✅ **Múltiples formatos de fecha** soportados
- ✅ **Thread-safe** para operaciones concurrentes
- ✅ **Validación de tipos** de datos
- ✅ **Envío de summary** por email, si se configuran adecuadamente las variable de entorno

## 🧪 Testing

Puedes usar el archivo de ejemplo incluido en `examples/sample_transactions.csv` para probar la funcionalidad.

## 📝 Notas

- Las transacciones se almacenan en memoria y se pierden al reiniciar el servidor
- El servicio está diseñado para ser reemplazado fácilmente por una base de datos real
- Los IDs se auto-incrementan si no se proporcionan en el CSV
- El servicio valida la estructura del CSV antes de procesar los datos
