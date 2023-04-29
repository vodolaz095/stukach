stukach
====================

Утилита, которая извлекает нежелательные письма из почтового ящика по протоколу IMAP и отправляет 
для обучения антиспам фильтров в rspamd.

Образец конфигурации
========================

```yaml

rspamd: 
  url: "http://localhost:11334/" # где слушает webUI rspamd с красивыми графиками
  username: "Rspamd controller password"
  password: "thisIsNotAPassword134"

inputs: # задаём строки соединения с imap серверами и название директории, откуда будем выгружать спам
  - server: "imap.example.org"
    port: 993
    username: "somebody"
    password: "thisIsNotAPassword134"
    useTLS: true
    directory: "Shared/Spam"

  - server: "imap.gmail.com" # для gmail может потребоваться получить пароль приложения с доступом к почте
    port: 993
    username: "somebody@gmail.com"
    password: "thisIsNotAPassword134"
    useTLS: true
    directory: "[Gmail]/Спам"

  - server: "imap.yandex.ru"# для yandex может потребоваться получить пароль приложения с доступом к почте
    port: 993
    username: "somebody@yandex.ru"
    password: "thisIsNotAPassword134"
    useTLS: true
    directory: "Spam"


```

Как запустить?
=================
Имитация - ничего в rspamd не посылается
```shell
$ stuckach --config ./config.yaml --dry 
```

Письма просто проверяются.
```shell
$ stuckach --config ./config.yaml 
```

Письма проверяются и записываются как спам
```shell
$ stuckach --config ./config.yaml --learn
```

Как скомпилировать приложение?
===================

Вариант 1. Сборка на хост-системе:

0. Linux - код проверялся на Centos 8 Stream
1. Golang > [1.19.4](https://go.dev/dl/)
2. GNU Make > [4.2.1](https://www.gnu.org/software/make/)
3. upx > [3.96](https://upx.github.io/)

```shell

$ make build
$ su -c 'make install'

```

Вариант 2. Сборка в докер контейнере

```shell

# ./docker_build.sh
# make install

```


