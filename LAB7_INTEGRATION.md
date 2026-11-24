# Лабораторная работа №7: Внедрение адаптивности, PWA и развертывание приложения

## 1. Назначение и стек
- **Предметная область:** Расчет времени отклика веб-приложения по задержкам компонентов (gateway, backend, database, cache, storage).
- **Бэкенд:** Go (Gin, GORM), PostgreSQL (реляционные данные), Redis (сессии и blacklist JWT), MinIO (S3 хранилище медиа).
- **Фронтенд:** React 19, TypeScript, Redux Toolkit, Vite PWA, Axios.
- **Сетевой контур и безопасность:**
  - Локальный Wi-Fi контур (`192.168.x.x`)
  - Node.js HTTPS Reverse Proxy с самоподписанным SSL-сертификатом (SAN: `localhost`, `127.0.0.1`, IP ноутбука)
  - GitHub Pages развертывание для статического PWA в режиме mock-данных.

## 2. Схема взаимодействия узлов при работе со смартфоном (ЛР7)
```
+-------------------------------------------------------------+
|              МОБИЛЬНЫЙ КЛИЕНТ (Смартфон)                    |
|   - PWA Standalone App / Chrome Mobile                      |
|   - URL: https://<IP_НОУТБУКА>:8443                         |
+-------------------------------------------------------------+
                              |
                              | HTTPS (TLS v1.2 / TLS v1.3)
                              v
+-------------------------------------------------------------+
|                НОУТБУК (Локальная сеть Wi-Fi)               |
|                                                             |
|  [Node.js HTTPS Reverse Proxy (:8443)]                      |
|    - SSL Termination (cert.crt, cert.key)                   |
|    - Раздача статики PWA (SPA fallback, Service Worker)     |
|    - /api/*               -> Проксирование на :8080         |
|    - /system-components/* -> Проксирование на :9000         |
|                                                             |
|  [Go Gin Web Service (:8080)]                               |
|    - REST API (10 эндпоинтов)                               |
|    - JWT Аутентификация, валидация прав                     |
|                                                             |
|  [PostgreSQL (:5432)] <--- SQL транзакции / ORM             |
|  [Redis (:6379)]      <--- Сессии, TTL 24h, JWT Blacklist   |
|  [MinIO S3 (:9000)]   <--- Изображения и видео              |
+-------------------------------------------------------------+
```

## 3. Запуск контура для показа ЛР7
1. Запуск БД, Redis и MinIO:
   ```bash
   docker-compose up -d
   ```
2. Запуск Go бэкенда:
   ```bash
   go run ./cmd/system_components/main.go
   ```
3. Запуск фронтенда и HTTPS реверс-прокси:
   ```bash
   cd ../system-components-frontend
   npm run build
   npm run serve:https
   ```
4. Открытие со смартфона в той же сети:
   `https://<IP_НОУТБУКА>:8443`
