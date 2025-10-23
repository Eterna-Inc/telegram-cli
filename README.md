 # Telegram CLI

 Telegram dan mesaj göndermek için Çapraz platform minimal CLI uygulaması.


 ## 🚀 Kurulum

 Aşağıda Windows, Linux ve macOS için yükleme ve test adımları bulunmaktadır.

 ### Windows (PowerShell)

 UAC veya ExecutionPolicy sorunları için:

 ```powershell
 Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
 ```

 Uzaktan tek satır ile (GitHub raw üzerinden çalıştırmak):

### 🪟 Windows (Kurulum)
 ```powershell
 iwr -useb https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1 | iex
 ```

### 🪟 Windows (Günclleme)
```powershell
iwr -useb https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1 | iex -ArgumentList '-Action update'
```



 ### 🐧 Linux (Kurulum)

 ```bash
 curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash
 ```

 ### 🐧 Linux (Güncelleme)

 ```bash
 curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash -s update
 ```

 ---

 ### 🍎 macOS

 macOS için Linux ile aynı bash tabanlı kurulum çalışır (curl | bash). Örnek:

 ```bash
 curl -fsSL https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.sh | bash
 ```

 ---

 Notlar:
 - `curl | bash` veya `iwr | iex` komutları hızlı kurulum sağlar, fakat çalıştırmadan önce içeriği incelemeniz güvenlik açısından önerilir.
 - Windows PowerShell'de ExecutionPolicy ve UAC/Administrator onayı gerektiğinde script kendi kendini yükseltmeye çalışacaktır.
