#!/usr/bin/env bash
# Автоматическая сборка и установка budget-cli
# Использование: ./install.sh [--user] [--prefix DIR]

set -euo pipefail

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Настройки по умолчанию
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="budget-cli"
BUILD_DIR="."
MAIN_PACKAGE="./cmd/budget"
USE_SUDO=true

# Парсинг аргументов
while [[ $# -gt 0 ]]; do
    case $1 in
        --user)
            INSTALL_DIR="$HOME/.local/bin"
            USE_SUDO=false
            shift
            ;;
        --prefix)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -h|--help)
            echo "Использование: $0 [опции]"
            echo "Опции:"
            echo "  --user       Установить в \$HOME/.local/bin (без sudo)"
            echo "  --prefix DIR Установить в указанную директорию"
            echo "  -h, --help   Показать эту справку"
            exit 0
            ;;
        *)
            echo -e "${RED}Неизвестная опция: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${YELLOW}=== Сборка budget-cli ===${NC}"

# Проверяем, что мы в правильной папке
if [[ ! -f "go.mod" ]]; then
    echo -e "${RED}Ошибка: go.mod не найден. Запустите скрипт из корня проекта.${NC}"
    exit 1
fi

# Сборка
echo "Собираем бинарник..."
if ! go build -o "$BINARY_NAME" "$MAIN_PACKAGE"; then
    echo -e "${RED}Ошибка сборки${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Сборка успешна${NC}"

# Проверяем, что бинарник работает
if ! ./"$BINARY_NAME" --help >/dev/null 2>&1; then
    echo -e "${RED}Ошибка: собранный бинарник не запускается${NC}"
    exit 1
fi

# Создаём директорию установки если нужно
if [[ ! -d "$INSTALL_DIR" ]]; then
    echo "Создаём директорию $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR"
fi

# Установка
echo "Устанавливаем в $INSTALL_DIR..."
if [[ "$USE_SUDO" == true ]]; then
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Проверка установки
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    INSTALLED_VERSION=$("$BINARY_NAME" --help 2>&1 | head -1)
    echo -e "${GREEN}✓ Установка успешна!${NC}"
    echo -e "  Команда: ${YELLOW}$BINARY_NAME${NC}"
    echo -e "  Путь:    ${YELLOW}$(command -v "$BINARY_NAME")${NC}"
    echo -e "  Версия:  ${YELLOW}$INSTALLED_VERSION${NC}"
else
    echo -e "${YELLOW}⚠ Бинарник установлен, но не найден в PATH${NC}"
    echo -e "  Добавьте $INSTALL_DIR в PATH:"
    echo -e "  ${YELLOW}export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
    echo -e "  Или перезапустите терминал"
fi

# Очистка временного бинарника
rm -f "$BINARY_NAME"

echo -e "${GREEN}Готово!${NC}"