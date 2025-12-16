![GitHub CI](https://github.com/rainb0w-clwn/go_auth_limiter/actions/workflows/ci.yml/badge.svg)
# Сервис "Анти-брутфорс"

## [ТЗ](./docs/TASK.md)

## Общее описание

Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе.

Сервис вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.

Предполагается, что сервис используется только для server-server, т.е. скрыт от конечного пользователя.

## Развертывание
* clone/download repo
* init config ($cp configs/config.example.yml configs/config.yml)
* run service ($make up)
