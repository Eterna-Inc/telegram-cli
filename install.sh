#!/usr/bin/env bash
set -e

REPO="Eterna-Inc/telegram-cli"
APP_NAME="telegram"
INSTALL_PATH="/usr/local/bin/${APP_NAME}"
CONFIG_PATH="${HOME}/.telegram_config"

show_help() {
  echo "Telegram CLI Installer"
  echo
  echo "Usage:  bash install.sh [command]"
  echo
  echo "Commands:"
  echo "  install    Install or update Telegram CLI"
  echo "  update     Force update to latest release"
  echo "  uninstall  Remove binary and config"
  echo
}

ensure_root() {
  if [ "$EUID" -ne 0 ]; then
    echo "⚙️  Root yetkisi gerekli, sudo ile yeniden başlatılıyor..."
    exec sudo bash "$0" "$@"
  fi
}

detect_platform() {
  ARCH=$(uname -m)
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l) ARCH="arm" ;;
  esac
}

install_app() {
  ensure_root
  detect_platform
  echo "💻 Platform: $OS-$ARCH"

  echo "🔍 En son sürüm aranıyor..."
  LATEST_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest \
    | grep "browser_download_url" \
    | grep "${OS}-${ARCH}" \
    | cut -d '"' -f 4)

  if [ -z "$LATEST_URL" ]; then
    echo "⚠️ Binary bulunamadı, kaynaktan derleniyor..."
    tmp=$(mktemp -d)
    cd "$tmp"
    git clone --depth=1 "https://github.com/$REPO" .
    go build -o "$INSTALL_PATH"
  else
    echo "⬇️  Binary indiriliyor..."
    curl -L "$LATEST_URL" -o "$INSTALL_PATH"
    chmod +x "$INSTALL_PATH"
  fi

  # Config oluştur
  if [ ! -f "$CONFIG_PATH" ]; then
    echo "🧩 Config oluşturuluyor..."
    cat > "$CONFIG_PATH" <<EOF
{
  "active_profile": "default",
  "profiles": {
    "default": {
      "token": "",
      "chat_id": ""
    }
  },
  "threads": {},
  "encrypted": false
}
EOF
    chown $(logname):$(logname) "$CONFIG_PATH"
    chmod 600 "$CONFIG_PATH"
  fi

  echo "✅ ${APP_NAME} yüklendi: $INSTALL_PATH"
  echo "ℹ️  telegram config --token <TOKEN> --chatid <CHAT_ID>"
  echo "ℹ️  telegram \"Kurulum tamamlandı ✅\""
}

update_app() {
  ensure_root
  detect_platform
  echo "🔄 Güncelleme başlatıldı..."
  rm -f "$INSTALL_PATH"
  install_app
  echo "✅ Güncelleme tamamlandı!"
}

uninstall_app() {
  ensure_root
  echo "🧹 ${APP_NAME} kaldırılıyor..."
  rm -f "$INSTALL_PATH"
  rm -f "$CONFIG_PATH"
  echo "✅ ${APP_NAME} tamamen kaldırıldı."
}

case "$1" in
  install|"")
    install_app
    ;;
  update)
    update_app
    ;;
  uninstall|remove)
    uninstall_app
    ;;
  -h|--help)
    show_help
    ;;
  *)
    echo "❌ Geçersiz komut: $1"
    show_help
    ;;
esac
