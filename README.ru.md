# KKT 54-ФЗ Мониторинг

![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)

Система мониторинга контрольно-кассовой техники (ККТ) в соответствии с Федеральным законом 54-ФЗ "О применении контрольно-кассовой техники".

**Бренд / сайт автора:** https://run-as-daemon.ru

## Описание

`kkt-54fz-monitoring` - это комплексное решение для мониторинга и анализа работы контрольно-кассовой техники, работающей по 54-ФЗ. Система собирает данные из различных источников (логи, HTTP API ОФД), агрегирует метрики и предоставляет их для Prometheus. Включает встроенную AI-подсистему для кластеризации ошибок и рекомендаций по алертам.

## Возможности

- 📊 **Сбор метрик** из файловых логов и HTTP API ОФД
- 📈 **Prometheus Exporter** с готовыми метриками
- 🚨 **Предустановленные правила алертов** для типовых проблем
- 📉 **Графана-дашборды** для визуализации
- 🤖 **AI-подсистема** для анализа ошибок и рекомендаций
- ⚙️ **Гибкая конфигурация** через YAML
- 🔒 **Безопасность** - встроенные проверки и валидация

## Быстрый старт

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/ranas-mukminov/kkt-54fz-monitoring.git
cd kkt-54fz-monitoring

# Сборка
make build

# Или с помощью Go
go build -o kkt-monitor ./cmd/kkt-monitor
```

### Запуск

```bash
# Запуск с конфигурацией по умолчанию
./kkt-monitor --config configs/config.yaml

# Запуск с Docker
docker-compose up -d
```

### Проверка метрик

```bash
curl http://localhost:9090/metrics
```

## Архитектура

```
┌─────────────┐     ┌─────────────┐
│  File Logs  │     │  HTTP OFD   │
└──────┬──────┘     └──────┬──────┘
       │                   │
       └───────┬───────────┘
               │
        ┌──────▼──────┐
        │ Collectors  │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │   Domain    │
        │    Model    │
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │  Prometheus │
        │   Exporter  │
        └──────┬──────┘
               │
               ├────────┐
               │        │
        ┌──────▼──────┐ │
        │ Prometheus  │ │
        └─────────────┘ │
                        │
                 ┌──────▼──────┐
                 │   Grafana   │
                 └─────────────┘
```

## Конфигурация

Пример конфигурационного файла `configs/config.yaml`:

```yaml
server:
  port: 9090
  metrics_path: /metrics

collectors:
  file_log:
    enabled: true
    path: /var/log/kkt/*.log
    format: json
    poll_interval: 10s
  
  http_ofd:
    enabled: true
    url: https://ofd.example.ru/api/v1
    api_key: ${OFD_API_KEY}
    poll_interval: 30s

ai:
  provider: mock  # mock, openai, anthropic
  error_clustering:
    enabled: true
    min_cluster_size: 5
  alert_advisor:
    enabled: true

logging:
  level: info
  format: json
```

## Метрики

Система экспортирует следующие метрики:

- `kkt_status` - статус ККТ (0=недоступна, 1=работает, 2=ошибка)
- `kkt_documents_total` - общее количество фискальных документов
- `kkt_errors_total` - количество ошибок по типам
- `kkt_ofd_sync_status` - статус синхронизации с ОФД
- `kkt_shift_status` - статус смены (открыта/закрыта)
- `kkt_last_document_timestamp` - timestamp последнего документа

## Алерты

Предустановленные правила алертов находятся в `configs/alerts/kkt-alerts.yaml`:

- Недоступность ККТ более 5 минут
- Критическая ошибка фискального накопителя
- Проблемы синхронизации с ОФД
- Переполнение памяти ФН
- Истечение срока действия ФН

## Разработка

### Требования

- Go 1.21+
- Make
- Docker и Docker Compose (для локальной разработки)

### Сборка и тестирование

```bash
# Установка зависимостей
make deps

# Запуск линтера
make lint

# Запуск тестов
make test

# Запуск интеграционных тестов
make test-integration

# Проверка безопасности
make security-check

# Проверка производительности
make perf-check

# Полная проверка (lint + test + security)
make check
```

### Структура проекта

```
.
├── cmd/
│   └── kkt-monitor/        # Точка входа приложения
├── internal/
│   ├── domain/             # Доменные модели
│   ├── config/             # Загрузка и валидация конфигов
│   ├── collector/          # Коллекторы данных
│   ├── exporter/           # Prometheus exporter
│   └── ai/                 # AI-подсистема
├── pkg/
│   ├── utils/              # Утилиты
│   └── logger/             # Логирование
├── configs/
│   ├── config.yaml         # Основная конфигурация
│   ├── alerts/             # Правила алертов
│   └── dashboards/         # Графана-дашборды
├── deployments/
│   ├── docker/             # Docker-файлы
│   └── kubernetes/         # K8s манифесты
├── test/
│   ├── testdata/           # Тестовые данные
│   └── integration/        # Интеграционные тесты
└── docs/                   # Документация
```

## AI-подсистема

### Кластеризация ошибок

AI-модуль автоматически группирует похожие ошибки для упрощения анализа:

```bash
curl http://localhost:9090/api/v1/ai/error-clusters
```

### Alert Advisor

Получение рекомендаций по настройке алертов на основе исторических данных:

```bash
curl http://localhost:9090/api/v1/ai/alert-recommendations
```

### Провайдеры

Поддерживаемые AI-провайдеры:
- **mock** - заглушка для разработки и тестирования
- **openai** - OpenAI GPT API
- **anthropic** - Anthropic Claude API

## Соответствие 54-ФЗ

Система соответствует требованиям:
- Федерального закона от 22.05.2003 № 54-ФЗ
- Приказа ФНС России от 21.03.2017 № ММВ-7-20/229@
- Технических требований к форматам фискальных документов

См. файл [LEGAL](LEGAL) для подробной информации о соответствии законодательству.

## Лицензия

Apache License 2.0. См. файл [LICENSE](LICENSE) для подробностей.

## Поддержка

- 📧 Email: support@run-as-daemon.ru
- 🐛 Issues: https://github.com/ranas-mukminov/kkt-54fz-monitoring/issues
- 📖 Документация: https://github.com/ranas-mukminov/kkt-54fz-monitoring/wiki

## Автор

© 2024 [run-as-daemon.ru](https://run-as-daemon.ru)

---

[English version](README.md)
