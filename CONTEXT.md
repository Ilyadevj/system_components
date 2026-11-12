cat << 'EOF' > CONTEXT.md
# Контекст проекта: Лабораторная работа №4 (Расчет времени отклика веб-приложения)

## Предметная область и архитектура
- Тема: Расчет времени отклика веб-приложения по задержкам компонентов (gateway, backend, database).
- Формула задержки: T = Σ(t_i * w(p_i)), где t_i — задержка узла, p_i — приоритет.
- Стек: Go (Gin, GORM), PostgreSQL, Redis, MinIO, Docker Compose.
- Клиенты: React SPA, нативное приложение Tauri.

## Реализованный функционал (ЛР-4)
1. Аутентификация и сессии:
   - Хранение паролей: bcrypt (CompareHashAndPassword / GenerateFromPassword).
   - Токены: JWT, передача через HTTP-only cookie `jwt_token`.
   - Сессии в Redis:
     - Активные сессии: `active_user:{token}` (TTL 24 часа).
     - Черный список токенов: `jwt_blacklist:{token}` (TTL 24 часа при Logout).
   - Middleware: `AuthMiddleware` проверяет отсутствие токена в `jwt_blacklist` перед вызовом контроллера.

2. API (10 REST-эндпоинтов):
   - POST /api/users/register
   - POST /api/users/login
   - POST /api/users/logout
   - GET  /api/system-components (фильтры: priority, grouping, name)
   - GET  /api/system-components/feed (параметры: id, next)
   - GET  /api/system-components/draft (получение по cookie)
   - POST /api/system-components (multipart: name, image, video -> MinIO)
   - PUT  /api/system-components/draft/publish (расчет времени отклика, публикация)
   - DELETE /api/system-components/:id (Soft delete через SQL UPDATE)
   - POST /api/system-components/:id/like (проставление/снятие лайка)

3. Модели БД (PostgreSQL):
   - User (id, username, password_hash)
   - SystemComponent (id, user_id, name, priority, grouping, response_time_ms, image_url, video_url, deleted_at)
   - Like (user_id, component_id)