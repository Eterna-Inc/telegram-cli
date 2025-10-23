#!/usr/bin/env bash
set -e

REPO="Eterna-Inc/telegram-cli"
APP_NAME="telegram"
INSTALL_PATH="/usr/local/bin/${APP_NAME}"
CONFIG_PATH="${HOME}/.telegram_config"

show_help() {
  echo "Telegram CLI Installer"
  echo
  echo "Usage: bash install.sh [command]"
  echo
  echo "Commands:"
  echo "  install     Install or update Telegram CLI"
  echo "  update      Force update to the latest release"
  echo "  uninstall   Remove the binary and configuration"
  echo
}

ensure_root() {
  if [ "$EUID" -ne 0 ]; then
    # macOS genellikle /usr/local/bin dizinine yazabilir, sudo gerekmez
    if [ -w "/usr/local/bin" ]; then
      echo "✅ Sudo not required (you already have write access to /usr/local/bin)"
      return
    fi
    echo "⚙️  Root privileges are required. Restarting with sudo..."
    exec sudo /usr/bin/env bash "$0" "$@"
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

set_owner() {
  # Config dosyasını gerçek kullanıcıya ait hale getir
  USER_NAME=$(logname 2>/dev/null || echo "$SUDO_USER")
  GROUP_NAME=$(id -gn "$USER_NAME" 2>/dev/null || echo "staff")
  if [ -n "$USER_NAME" ]; then
    chown "$USER_NAME:$GROUP_NAME" "$1" 2>/dev/null || true
  fi
}

install_app() {
  ensure_root
  detect_platform
  echo "💻 Platform detected: $OS-$ARCH"

  echo "🔍 Fetching latest release..."
  LATEST_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest \
    | grep "browser_download_url" \
    | grep "${OS}-${ARCH}" \
    | cut -d '"' -f 4)

  if [ -z "$LATEST_URL" ]; then
    echo "⚠️  No prebuilt binary found. Building from source..."
    tmp=$(mktemp -d)
    cd "$tmp"
    git clone --depth=1 "https://github.com/$REPO" .
    CGO_ENABLED=0 go build -ldflags="-s -w" -o "$INSTALL_PATH"
  else
    echo "⬇️  Downloading latest binary..."
    curl -L "$LATEST_URL" -o "$INSTALL_PATH"
    chmod +x "$INSTALL_PATH"
  fi

  # Create config file if it doesn’t exist
  if [ ! -f "$CONFIG_PATH" ]; then
    echo "🧩 Creating configuration file..."
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
    set_owner "$CONFIG_PATH"
    chmod 600 "$CONFIG_PATH"
  fi

  echo
  echo "✅ ${APP_NAME} installed successfully at: $INSTALL_PATH"
  echo "ℹ️  Configure your Telegram credentials:"
  echo "   telegram config --token <TOKEN> --chatid <CHAT_ID>"
  echo
  echo "ℹ️  Test it:"
  echo "   telegram \"Installation complete ✅\""
  echo
}

update_app() {
  ensure_root
  detect_platform
  echo "🔄 Updating Telegram CLI..."
  rm -f "$INSTALL_PATH"
  install_app
  echo "✅ Update completed successfully!"
}

uninstall_app() {
  ensure_root
  echo "🧹 Uninstalling ${APP_NAME}..."
  rm -f "$INSTALL_PATH" 2>/dev/null || true
  rm -f "$CONFIG_PATH" 2>/dev/null || true
  echo "✅ ${APP_NAME} has been completely removed."
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
    echo "❌ Invalid command: $1"
    show_help
    ;;
esac
