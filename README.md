![GitHub CI](https://github.com/rainb0w-clwn/go_auth_limiter/actions/workflows/ci.yml/badge.svg)
# Сервис "Анти-брутфорс"


## Запуск

```bash
make run
 ```

## CLI

1. Очистить бакет
```bash
 make run-cli ARGS="reset_bucket email@x.com 192.168.0.1"
 ```

2. Добавить подсеть в черный список
```bash
 make run-cli ARGS="add_cidr_to_black_list 192.168.1.1/24" 
 ```

3. Добавить подсеть в белый список
```bash
 make run-cli ARGS="add_cidr_to_white_list 192.168.1.0/24" 
 ```

4. Удалить подсеть из черного списка
```bash
 make run-cli ARGS="delete_cidr_from_black_list 192.168.1.1/24" 
 ```

5. Удалить подсеть из белого списка
```bash
 make run-cli ARGS="delete_cidr_from_white_list 192.168.1.0/24" 
 ```

## API

- [GRPC](./proto/limiter/AuthLimiter.proto) 

- [HTTP](./proto/limiter/AuthLimiter.openapi.yaml)

## Тесты

1. UNIT-Тесты
```bash
make test
 ```

2. Интеграционные тесты
```bash
make integration-test
 ```
