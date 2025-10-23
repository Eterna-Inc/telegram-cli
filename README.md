# Telegram CLI

[![Windows](https://img.shields.io/badge/Windows-0078D6?style=flat-square&logo=windows&logoColor=white)](#windows)
[![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat-square&logo=linux&logoColor=black)](#linux)
[![macOS](https://img.shields.io/badge/macOS-000000?style=flat-square&logo=apple&logoColor=white)](#macos)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](#)

---

## 🇹🇷 Türkçe

Telegram'dan mesaj göndermek için çapraz platform, minimal ve güçlü CLI uygulaması.  
Bir kez yapılandırdıktan sonra sadece terminalden yazmanız yeterlidir 👇

```bash
telegram "deploy bitti 🚀"
```


### ⚙️ Özellikler

*   📨 Mesaj, komut çıktısı, kod bloğu veya dosya gönderimi
*   🕐 Zamanlanmış (scheduled) mesajlar (`--in 10m`, `--at 18:00`)
*   🧩 Çoklu profil desteği (`--profile <name>`)
*   🧾 Markdown veya HTML format desteği (`--format markdown|html`)
*   💬 Komut çıktısını direkt Telegram’a yollama (`-exec "ls -la"`)
*   🧠 Basit, hızlı, platform bağımsız (Windows, Linux, macOS)

***

## 🚀 Kurulum

### 🪟 Windows (Kurulum)

```bash
iwr -useb https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1 | iex
```

### 🪟 Windows (Güncelleme)

```powershell 
Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1' -OutFile $env:TEMP\telegram-install.ps1 Start-Process powershell -Verb runAs -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `$env:TEMP\telegram-install.ps1` -Action update"
```

### 🐧 Linux (Kurulum)

```bash
curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash
```

### 🐧 Linux (Güncelleme)

```bash
curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash -s update
```

### 🍎 macOS

macOS için de aynı bash tabanlı kurulum geçerlidir:

```bash
curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash
```

***

## 🧩 İlk Kurulum Sonrası

İlk kullanımda Telegram API token’ınızı ve chat ID’nizi ayarlayın:

```bash
telegram config --token 123456:ABCDEF --chatid 987654321 --profile default
```

***

## 🧠 Kullanım Örnekleri

| Amaç | Komut |
| --- | --- |
| Basit mesaj gönder | telegram "Deploy tamam ✅" |
| Kod bloğu olarak gönder | telegram --code "npm run build && pm2 restart all" |
| Markdown mesajı | telegram --format markdown "*Deploy* **OK** ✅" |
| HTML mesajı | telegram --format html "<b>Deploy</b> <i>OK</i> ✅" |
| Komut çıktısı gönder | telegram --exec "df -h" |
| 10 dakika sonra mesaj gönder | telegram --in 10m "Reminder: check logs" |
| Belirli saatte mesaj gönder | telegram --at 18:00 "Eve gitme vakti!" |
| Belirli profile mesaj gönder | telegram --profile prod "Prod deploy complete" |

***

## 🧾 Config Dosyası

Varsayılan olarak oluşturulan yapı `~/.telegram_config` dosyasında tutulur:

`{   "active_profile": "default",   "profiles": {     "default": {       "token": "123456:ABCDEF",       "chat_id": "987654321"     }   },   "threads": {}}`

***

## 🧰 Geliştiriciler İçin

Kaynak koddan derlemek istersen:

```bash
git clone https://github.com/Eterna-Inc/telegram-cli.git 
cd telegram-cli
go build -o telegram .
```

***

## 💡 Notlar

*   `curl | bash` veya `iwr | iex` komutları hızlı kurulum sağlar,  
    ancak güvenlik için çalıştırmadan önce içeriğini incelemeniz önerilir.
*   Windows PowerShell’de `ExecutionPolicy` veya UAC gerekirse script kendini yükseltir.
*   macOS’ta `/usr/local/bin` dizinine erişiminiz varsa `sudo` gerekmez.
*   Config dosyası kullanıcıya ait olacak şekilde otomatik ayarlanır.
    

***

## 🇬🇧 English

Minimal, cross-platform CLI app for sending Telegram messages directly from your terminal.  
Once configured, simply type:

`telegram "deploy finished 🚀"`

### ⚙️ Features

*   📨 Send text, command output, or code blocks
*   🕐 Scheduled messages (`--in 10m`, `--at 18:00`)
*   🧩 Multi-profile support (`--profile <name>`)
*   🧾 Markdown or HTML message formatting
*   💬 Run commands and send their output (`--exec "uptime"`)
*   ⚡ Lightweight, fast, and cross-platform
    

***

## 🚀 Installation

### 🪟 Windows

```bash
iwr -useb https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1 | iex
```

### 🐧 Linux / 🍎 macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash
```

***

## 🧩 Configuration

First, set your Telegram credentials:

```bash
telegram config --token 123456:ABCDEF --chatid 987654321 --profile default
```

***

## 🧠 Usage Examples

| Action | Command |
| --- | --- |
| Simple message | telegram "Hello from CLI" |
| Code block | telegram --code "systemctl restart app" |
| Markdown | telegram --format markdown "*Deploy OK* ✅" |
| HTML | telegram --format html "<b>Deploy OK</b> ✅" |
| Command output | telegram --exec "ls -lah" |
| Scheduled message | telegram --in 10m "Time to check logs" |
| Send to specific profile | telegram --profile prod "Production updated" |

***

## 🧾 Config File

Stored at `~/.telegram_config`:

`{   "active_profile": "default",   "profiles": {     "default": {       "token": "123456:ABCDEF",       "chat_id": "987654321"     }   },   "threads": {}}`

***

## 🧰 Developers

```bash
git clone https://github.com/Eterna-Inc/telegram-cli.git 
cd telegram-cli
go build -o telegram .
```

* * *

## 🧩 Notes

*   Always inspect scripts before running (`curl | bash` / `iwr | iex`).
*   On Windows, PowerShell may request elevation (UAC).
*   macOS and Linux installers auto-detect privileges.
*   The config file is automatically owned by the correct user.