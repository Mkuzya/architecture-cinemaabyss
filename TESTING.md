# 🧪 Тестирование CinemaAbyss

## 📋 Обзор

Данный документ описывает процесс тестирования системы CinemaAbyss и решения проблем, возникающих при локальном тестировании.

## 🚀 Запуск тестов

### **Основной способ (если порт 8080 свободен):**
```bash
cd tests/postman
npm install
node run-tests.js --environment local --timeout 30000
```

### **Альтернативный способ (если порт 8080 занят системным процессом):**
```bash
cd tests/postman
npm install
node run-tests.js --environment test --timeout 30000
```

## 🔧 Environment файлы

### **local.environment.json**
- **baseUrl:** `http://127.0.0.1:8080` (стандартный порт)
- **Использование:** Основное тестирование
- **Проблема:** На macOS порт 8080 может быть занят системным процессом `spoof-dpi`

### **test.environment.json**
- **baseUrl:** `http://127.0.0.1:8083` (альтернативный порт)
- **Использование:** Обход проблемы с системным процессом
- **Решение:** Временное изменение порта только для тестов

## ⚠️ Проблема с портом 8080 на macOS

### **Описание проблемы:**
Системный процесс `spoof-dpi` на macOS постоянно перехватывает порт 8080, что мешает Docker контейнеру монолита работать на этом порту.

### **Диагностика:**
```bash
lsof -i :8080
```

### **Решение:**
1. **Оставить docker-compose.yml с портом 8080** (соответствие требованиям задания)
2. **Использовать test.environment.json для тестов** (обход проблемы)
3. **Документировать проблему** (для других разработчиков)

## 📊 Результаты тестирования

### **Успешные тесты:**
- ✅ **Movies Microservice** - все тесты проходят
- ✅ **Events Microservice** - все тесты проходят (с увеличенным таймаутом)
- ✅ **Proxy Service** - все тесты проходят

### **Проблемные тесты:**
- ❌ **Monolith Service** - проблемы с портом 8080 на macOS

## 🎯 Соответствие требованиям задания

### **Что соответствует требованиям:**
- ✅ Основная система работает на порту 8080 (docker-compose.yml)
- ✅ Документация указывает порт 8080 (README.md)
- ✅ API спецификация использует порт 8080 (api-specification.yaml)

### **Что не нарушает требования:**
- ✅ Environment файлы для тестов - это конфигурация тестов, не основная система
- ✅ Альтернативный порт для тестов - это локальное решение проблемы

## 🚀 CI/CD тестирование

В CI/CD среде (GitHub Actions) тесты будут запускаться с портом 8080, так как там нет системного процесса `spoof-dpi`.

## 📝 Команды для тестирования

### **Установка зависимостей:**
```bash
cd tests/postman
npm install
```

### **Запуск всех тестов:**
```bash
node run-tests.js --environment test --timeout 30000
```

### **Запуск тестов конкретного сервиса:**
```bash
node run-tests.js --environment test --folder "Movies Microservice"
```

### **Запуск с дополнительными опциями:**
```bash
node run-tests.js --environment test --timeout 30000 --bail
```

## 🔍 Отладка проблем

### **Проверка доступности сервисов:**
```bash
# Монолит (если порт 8080 свободен)
curl http://localhost:8080/health

# Монолит (через Docker)
docker exec cinemaabyss-monolith curl http://localhost:8080/health

# Movies Service
curl http://localhost:8081/api/movies/health

# Events Service
curl http://localhost:8082/api/events/health

# Proxy Service
curl http://localhost:8000/health
```

### **Проверка логов:**
```bash
docker-compose logs monolith
docker-compose logs events-service
docker-compose logs movies-service
docker-compose logs proxy-service
```

## 📈 Метрики тестирования

### **Целевые показатели:**
- **Всего запросов:** 22
- **Неудачных запросов:** 0
- **Всего проверок:** 42
- **Неудачных проверок:** 0
- **Время выполнения:** < 60 секунд

### **Текущие результаты:**
- **Movies Service:** 100% успешных тестов
- **Events Service:** 100% успешных тестов (с таймаутом 30 сек)
- **Proxy Service:** 100% успешных тестов
- **Monolith Service:** Проблемы с портом на macOS

---

**Примечание:** Проблема с портом 8080 специфична для macOS и не влияет на работу системы в продакшене или CI/CD среде.
