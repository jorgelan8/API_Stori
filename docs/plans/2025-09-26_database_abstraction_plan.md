# 🗄️ **Plan de Implementación: Abstracción de Base de Datos**

## 📋 **Análisis de la Situación Actual**

### **✅ Estado Actual Identificado:**
- **MockDatabase**: Implementación completa con 6 métodos principales
- **Servicios**: MigrationService y UsersService dependen directamente de MockDatabase
- **Métodos actuales**:
  - `SaveTransaction()` - Guardar transacciones
  - `GetTransaction()` - Obtener por ID
  - `GetTransactionsByUserID()` - Obtener por usuario
  - `GetTransactionsByUserIDWithDateRange()` - Con filtros de fecha
  - `GetAllTransactions()` - Obtener todas
  - `GetTransactionCount()` - Contar transacciones
  - `ClearTransactions()` - Limpiar (para testing)

---

## 🎯 **Objetivos del Plan**

1. **Mantener MockDatabase** como implementación actual
2. **Crear capa de abstracción** con interfaces
3. **Implementar repositorio MySQL** real
4. **Configuración flexible** mock/BD real por variable
5. **BD externa** (fuera del contenedor)
6. **Tests flexibles** (mock o BD real)
7. **Migraciones** de base de datos

---

## 🏗️ **Arquitectura Propuesta**

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Services      │    │   Interfaces     │    │  Implementations│
│                 │    │                  │    │                 │
│ MigrationService│───▶│ TransactionRepo  │◄───│ MockTransactionRepo│
│ UsersService    │    │ UserRepo         │    │ MySQLTransactionRepo│
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

---

## 📅 **Plan de Implementación Detallado**

### **🔴 Fase 1: Diseño y Interfaces (Día 1-2)**

#### **1.1 Crear Interfaces de Repositorio**
```go
// internal/repositories/interfaces.go
package repositories

import (
    "api-stori/internal/models"
    "context"
    "time"
)

type TransactionRepository interface {
    Save(ctx context.Context, transaction *models.UserTransaction) (*models.UserTransaction, error)
    GetByID(ctx context.Context, id int) (*models.UserTransaction, error)
    GetByUserID(ctx context.Context, userID int) ([]models.UserTransaction, error)
    GetByUserIDWithDateRange(ctx context.Context, userID int, from, to *time.Time) ([]models.UserTransaction, error)
    GetAll(ctx context.Context) ([]models.UserTransaction, error)
    Count(ctx context.Context) (int, error)
    Clear(ctx context.Context) error
}

type UserRepository interface {
    GetByID(ctx context.Context, id int) (*models.User, error)
    Exists(ctx context.Context, id int) (bool, error)
}
```

#### **1.2 Crear Factory Pattern**
```go
// internal/repositories/factory.go
package repositories

import (
    "api-stori/internal/config"
    "context"
)

type RepositoryFactory struct {
    config *config.DatabaseConfig
}

func NewRepositoryFactory(config *config.DatabaseConfig) *RepositoryFactory {
    return &RepositoryFactory{config: config}
}

func (f *RepositoryFactory) CreateTransactionRepository(ctx context.Context) (TransactionRepository, error) {
    switch f.config.Type {
    case "mock":
        return NewMockTransactionRepository(), nil
    case "mysql":
        return NewMySQLTransactionRepository(f.config.MySQL), nil
    default:
        return NewMockTransactionRepository(), nil
    }
}
```

### **🟡 Fase 2: Implementación MySQL (Día 3-4)**

#### **2.1 Configuración de Base de Datos**
```go
// internal/config/database.go
package config

type DatabaseConfig struct {
    Type  string      `json:"type"`  // "mock" | "mysql"
    MySQL MySQLConfig `json:"mysql"`
}

type MySQLConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
    Password string `json:"password"`
    Charset  string `json:"charset"`
    MaxConns int    `json:"max_conns"`
    MaxIdle  int    `json:"max_idle"`
}
```

#### **2.2 Implementación MySQL**
```go
// internal/repositories/mysql_transaction_repository.go
package repositories

import (
    "api-stori/internal/config"
    "api-stori/internal/models"
    "context"
    "database/sql"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
)

type MySQLTransactionRepository struct {
    db *sql.DB
}

func NewMySQLTransactionRepository(config config.MySQLConfig) (*MySQLTransactionRepository, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
        config.Username, config.Password, config.Host, config.Port, 
        config.Database, config.Charset)
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configurar pool de conexiones
    db.SetMaxOpenConns(config.MaxConns)
    db.SetMaxIdleConns(config.MaxIdle)
    db.SetConnMaxLifetime(time.Hour)
    
    // Verificar conexión
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return &MySQLTransactionRepository{db: db}, nil
}

func (r *MySQLTransactionRepository) Save(ctx context.Context, transaction *models.UserTransaction) (*models.UserTransaction, error) {
    query := `INSERT INTO transactions (user_id, amount, datetime) VALUES (?, ?, ?)`
    result, err := r.db.ExecContext(ctx, query, transaction.UserID, transaction.Amount, transaction.DateTime)
    if err != nil {
        return nil, err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }
    
    transaction.ID = int(id)
    return transaction, nil
}

// ... implementar otros métodos
```

#### **2.3 Migraciones de Base de Datos**
```go
// internal/migrations/migrations.go
package migrations

import (
    "database/sql"
    "fmt"
)

func RunMigrations(db *sql.DB) error {
    migrations := []string{
        createTransactionsTable,
        createIndexes,
    }
    
    for i, migration := range migrations {
        if _, err := db.Exec(migration); err != nil {
            return fmt.Errorf("migration %d failed: %v", i+1, err)
        }
    }
    
    return nil
}

const createTransactionsTable = `
CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    datetime DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_datetime (datetime),
    INDEX idx_user_datetime (user_id, datetime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`

const createIndexes = `
-- Índices adicionales para optimización
CREATE INDEX IF NOT EXISTS idx_user_amount ON transactions (user_id, amount);
CREATE INDEX IF NOT EXISTS idx_datetime_range ON transactions (datetime, user_id);
`
```

### **🟢 Fase 3: Refactorización de Servicios (Día 5-6)**

#### **3.1 Actualizar Servicios**
```go
// internal/services/migration_service.go
type MigrationService struct {
    transactionRepo repositories.TransactionRepository
    reportService   *ReportService
}

func NewMigrationService(transactionRepo repositories.TransactionRepository) *MigrationService {
    return &MigrationService{
        transactionRepo: transactionRepo,
        reportService:   NewReportServiceWithMockMode(defaultConfig),
    }
}

func (ms *MigrationService) ProcessCSV(reader io.Reader) (*MigrationStats, error) {
    // ... lógica existente ...
    
    // Usar repositorio en lugar de database directo
    savedTransaction, err := ms.transactionRepo.Save(ctx, &transaction)
    if err != nil {
        stats.UpdateError(lineNumber, err)
        continue
    }
    
    // ... resto de la lógica ...
}
```

#### **3.2 Actualizar Configuración Principal**
```go
// internal/config/config.go
type Config struct {
    App      AppConfig      `json:"app"`
    Email    EmailConfig    `json:"email"`
    Report   ReportConfig   `json:"report"`
    Database DatabaseConfig `json:"database"` // ← NUEVO
}

func loadDatabaseConfig() DatabaseConfig {
    return DatabaseConfig{
        Type: getEnvOrDefault("DB_TYPE", "mock"),
        MySQL: MySQLConfig{
            Host:     getEnvOrDefault("DB_HOST", "localhost"),
            Port:     getEnvIntOrDefault("DB_PORT", 3306),
            Database: getEnvOrDefault("DB_NAME", "api_stori"),
            Username: getEnvOrDefault("DB_USER", "root"),
            Password: os.Getenv("DB_PASSWORD"),
            Charset:  getEnvOrDefault("DB_CHARSET", "utf8mb4"),
            MaxConns: getEnvIntOrDefault("DB_MAX_CONNS", 10),
            MaxIdle:  getEnvIntOrDefault("DB_MAX_IDLE", 5),
        },
    }
}
```

### **🔵 Fase 4: Testing y Docker (Día 7-8)**

#### **4.1 Tests Flexibles**
```go
// tests/test_utils/database_factory.go
package test_utils

import (
    "api-stori/internal/config"
    "api-stori/internal/repositories"
    "context"
    "os"
)

func SetupTestDatabase(ctx context.Context) (repositories.TransactionRepository, error) {
    dbType := os.Getenv("TEST_DB_TYPE")
    if dbType == "" {
        dbType = "mock" // Default para tests
    }
    
    config := &config.DatabaseConfig{
        Type: dbType,
        MySQL: config.MySQLConfig{
            Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
            Port:     getEnvIntOrDefault("TEST_DB_PORT", 3306),
            Database: getEnvOrDefault("TEST_DB_NAME", "api_stori_test"),
            Username: getEnvOrDefault("TEST_DB_USER", "root"),
            Password: os.Getenv("TEST_DB_PASSWORD"),
        },
    }
    
    factory := repositories.NewRepositoryFactory(config)
    return factory.CreateTransactionRepository(ctx)
}
```

#### **4.2 Docker Compose Actualizado**
```yaml
# docker-compose.yml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_TYPE=mysql
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_NAME=api_stori
      - DB_USER=root
      - DB_PASSWORD=password
    depends_on:
      - mysql
    networks:
      - api-network

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: api_stori
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./internal/migrations:/docker-entrypoint-initdb.d
    networks:
      - api-network

volumes:
  mysql_data:

networks:
  api-network:
    driver: bridge
```

### **🟣 Fase 5: Variables de Entorno (Día 9)**

#### **5.1 Actualizar env.example**
```bash
# Database Configuration
DB_TYPE=mock                    # mock | mysql
DB_HOST=localhost
DB_PORT=3306
DB_NAME=api_stori
DB_USER=root
DB_PASSWORD=your_password
DB_CHARSET=utf8mb4
DB_MAX_CONNS=10
DB_MAX_IDLE=5

# Test Database Configuration
TEST_DB_TYPE=mock               # mock | mysql
TEST_DB_HOST=localhost
TEST_DB_PORT=3306
TEST_DB_NAME=api_stori_test
TEST_DB_USER=root
TEST_DB_PASSWORD=your_password
```

#### **5.2 Scripts de Testing**
```bash
# test_with_mock.sh
#!/bin/bash
export TEST_DB_TYPE=mock
go test -v ./tests/...

# test_with_mysql.sh
#!/bin/bash
export TEST_DB_TYPE=mysql
export TEST_DB_HOST=localhost
export TEST_DB_PASSWORD=password
go test -v ./tests/...
```

---

## 📊 **Estructura de Archivos Propuesta**

```
internal/
├── config/
│   ├── config.go          # Configuración principal
│   └── database.go        # Configuración de BD
├── repositories/
│   ├── interfaces.go      # Interfaces de repositorio
│   ├── factory.go         # Factory pattern
│   ├── mock_transaction_repository.go  # Mock (refactorizado)
│   └── mysql_transaction_repository.go # MySQL implementation
├── migrations/
│   ├── migrations.go      # Migraciones
│   └── schema.sql         # Schema SQL
└── services/
    ├── migration_service.go  # Refactorizado
    └── users_service.go      # Refactorizado

tests/
├── test_utils/
│   └── database_factory.go  # Factory para tests
└── integration/
    └── database_test.go     # Tests de BD real
```

---

## 🎯 **Cronograma de Implementación**

| Día | Fase | Tareas | Entregables |
|-----|------|--------|-------------|
| 1-2 | **Interfaces** | Crear interfaces y factory | Interfaces, Factory pattern |
| 3-4 | **MySQL** | Implementar MySQL repo | Repositorio MySQL, Migraciones |
| 5-6 | **Refactor** | Actualizar servicios | Servicios refactorizados |
| 7-8 | **Testing** | Tests flexibles, Docker | Tests, Docker compose |
| 9 | **Config** | Variables de entorno | env.example, scripts |

---

## ✅ **Criterios de Éxito**

### **Funcionalidad**
- [ ] MockDatabase funciona igual que antes
- [ ] MySQL repository implementa todas las operaciones
- [ ] Servicios funcionan con ambos repositorios
- [ ] Tests pasan con mock y MySQL

### **Configuración**
- [ ] Variable `DB_TYPE` controla el tipo de BD
- [ ] Variable `TEST_DB_TYPE` controla tests
- [ ] Docker compose conecta con MySQL externo
- [ ] Migraciones se ejecutan automáticamente

### **Calidad**
- [ ] Interfaces bien definidas
- [ ] Error handling consistente
- [ ] Logging de operaciones de BD
- [ ] Tests de integración con BD real

---

## 🚀 **Próximos Pasos Inmediatos**

1. **Crear interfaces** de repositorio
2. **Refactorizar MockDatabase** para implementar interface
3. **Actualizar configuración** para soporte de BD
4. **Implementar factory pattern**

---

## 📚 **Dependencias Adicionales Requeridas**

### **Go Modules**
```go
// go.mod - Agregar dependencias
require (
    github.com/go-sql-driver/mysql v1.7.1
    github.com/golang-migrate/migrate/v4 v4.16.2
)
```

### **Docker**
```dockerfile
# Dockerfile - Agregar para migraciones
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/internal/migrations ./migrations
CMD ["./main"]
```

---

## 🔧 **Comandos de Desarrollo**

### **Desarrollo Local**
```bash
# Con MockDatabase (default)
go run cmd/api/main.go

# Con MySQL
export DB_TYPE=mysql
export DB_PASSWORD=your_password
go run cmd/api/main.go
```

### **Testing**
```bash
# Tests con Mock
make test-all

# Tests con MySQL
export TEST_DB_TYPE=mysql
export TEST_DB_PASSWORD=password
make test-all
```

### **Docker**
```bash
# Desarrollo con MySQL
docker-compose up -d

# Solo API (BD externa)
docker-compose up api
```

---

**📅 Fecha de creación**: 2024-01-XX  
**👨‍💻 Creado por**: Claude (AI Assistant)  
**🎯 Objetivo**: Implementar abstracción de base de datos con soporte Mock/MySQL  
**⏱️ Duración estimada**: 9 días de desarrollo
