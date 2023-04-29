#!/usr/bin/env bash

set -e

# можно использовать и docker, и podman
containerEngine=docker
# containerEngine=podman

rm -f build/stukach

# создаём базовый образ
"$containerEngine" build --progress plain -t stukach:latest .

# запускаем скрипт entrypoint.sh в базовом образе
"$containerEngine" run -v `pwd`/build:/build/ stukach:latest
