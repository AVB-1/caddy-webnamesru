# Этап сборки: используем официальный builder-образ Caddy для компиляции с xcaddy
FROM caddy:2.11.4-builder AS builder

# Собираем Caddy с вашим модулем webnamesru
# Обратите внимание: замените ваш-username на ваше реальное имя пользователя GitHub
RUN xcaddy build \
    --with github.com/AVB-1/caddy-webnamesru@v0.1.0

# Финальный этап: используем лёгкий официальный образ Caddy
FROM caddy:2.11.4-alpine

# Копируем собранный бинарный файл
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
