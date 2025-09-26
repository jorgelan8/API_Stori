# Balance Service - API Endpoints

Este documento describe los endpoints disponibles para el servicio de balance de usuarios.

## 🚀 Endpoints Disponibles

### 1. GET /api/v1/users/{user_id}/balance
**Descripción**: Obtiene el saldo de un usuario específico con opción de filtrar por rango de fechas.

**Request**:
- **Method**: GET
- **Path Parameters**: 
  - `user_id` (int) - ID del usuario
- **Query Parameters** (opcionales):
  - `from` (string) - Fecha de inicio en formato "YYYY-MM-DDTHH:MM:SSZ"
  - `to` (string) - Fecha de fin en formato "YYYY-MM-DDTHH:MM:SSZ"

**Ejemplos de uso**:

#### Obtener balance completo de un usuario
```bash
curl -X GET http://localhost:8080/api/v1/users/1001/balance
```

#### Obtener balance con filtro de fecha
```bash
curl -X GET "http://localhost:8080/api/v1/users/1001/balance?from=2024-01-15T00:00:00Z&to=2024-01-20T23:59:59Z"
```

#### Obtener balance desde una fecha específica
```bash
curl -X GET "http://localhost:8080/api/v1/users/1001/balance?from=2024-01-15T00:00:00Z"
```

#### Obtener balance hasta una fecha específica
```bash
curl -X GET "http://localhost:8080/api/v1/users/1001/balance?to=2024-01-20T23:59:59Z"
```

**Response**:
```json
{
  "balance": 4.95,
  "total_debits": 10.05,
  "total_credits": 15.00
}
```

**Error Responses**:

#### Usuario no encontrado (400)
```json
HTTP/1.1 400 Bad Request
User not found
```

#### Formato de fecha inválido (400)
```json
HTTP/1.1 400 Bad Request
Invalid 'from' date format. Expected: YYYY-MM-DDTHH:MM:SSZ
```

#### Rango de fechas inválido (400)
```json
HTTP/1.1 400 Bad Request
Invalid date range: 'from' date must be before 'to' date
```

#### User ID inválido (400)
```json
HTTP/1.1 400 Bad Request
Invalid user_id format
```

## 📊 Formato de Respuesta

### BalanceResponse
```json
{
  "balance": float64,      // Saldo total del usuario
  "total_debits": int,     // Número total de transacciones negativas (débitos)
  "total_credits": int     // Número total de transacciones positivas (créditos)
}
```

## 🔧 Validaciones

### Parámetros de Entrada
- **user_id**: Debe ser un número entero válido
- **from**: Debe estar en formato "YYYY-MM-DDTHH:MM:SSZ" (ISO 8601 con Z)
- **to**: Debe estar en formato "YYYY-MM-DDTHH:MM:SSZ" (ISO 8601 con Z)

### Reglas de Negocio
- Si se proporcionan ambas fechas (`from` y `to`), `from` debe ser anterior a `to`
- Si no se proporcionan fechas, se incluyen todas las transacciones del usuario
- Si solo se proporciona `from`, se incluyen transacciones desde esa fecha en adelante
- Si solo se proporciona `to`, se incluyen transacciones hasta esa fecha

### Códigos de Error
- **400 Bad Request**: 
  - Usuario no encontrado
  - Formato de fecha inválido
  - Rango de fechas inválido
  - User ID inválido

## 📝 Formato de Fechas

El servicio acepta fechas en formato ISO 8601 con zona horaria UTC:

```
YYYY-MM-DDTHH:MM:SSZ
```

**Ejemplos válidos**:
- `2024-01-15T10:30:00Z`
- `2024-12-31T23:59:59Z`
- `2024-01-01T00:00:00Z`

**Ejemplos inválidos**:
- `2024-01-15 10:30:00` (falta la T y Z)
- `2024-01-15T10:30:00` (falta la Z)
- `15-01-2024T10:30:00Z` (formato de fecha incorrecto)

## 🧪 Testing

### Casos de Prueba

1. **Usuario existente sin filtros**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/users/1001/balance
   ```

2. **Usuario no existente**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/users/9999/balance
   ```

3. **Formato de fecha inválido**:
   ```bash
   curl -X GET "http://localhost:8080/api/v1/users/1001/balance?from=2024-01-15"
   ```

4. **Rango de fechas inválido**:
   ```bash
   curl -X GET "http://localhost:8080/api/v1/users/1001/balance?from=2024-01-20T00:00:00Z&to=2024-01-15T23:59:59Z"
   ```

5. **User ID inválido**:
   ```bash
   curl -X GET http://localhost:8080/api/v1/users/abc/balance
   ```

## 📊 Cálculo del Balance

El balance se calcula de la siguiente manera:

1. **Balance**: Suma de todos los montos de las transacciones filtradas
2. **Total Débitos**: Suma de transacciones con monto negativo
3. **Total Créditos**: Suma de transacciones con monto positivo

**Ejemplo**:
- Transacción 1: +100.00 (crédito)
- Transacción 2: -50.00 (débito)
- Transacción 3: +25.21 (crédito)
- Transacción 4: -25.00 (débito)

**Resultado**:
```json
{
  "balance": 50.21,
  "total_debits": 75.00,
  "total_credits": 125.21
}
```

## 🔗 Integración con Migration Service

El Balance Service utiliza las transacciones almacenadas por el Migration Service:

1. Primero se debe ejecutar el Migration Service para cargar transacciones
2. Luego se puede consultar el balance de cualquier usuario

## 📝 Notas

- Las transacciones se almacenan en memoria y se pierden al reiniciar el servidor
- El servicio valida todos los parámetros de entrada antes de procesar
- Los errores se devuelven con códigos HTTP apropiados y mensajes descriptivos
- El formato de fecha es estricto y debe incluir la zona horaria UTC (Z)
