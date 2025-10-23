<#
Telegram CLI Installer for Windows
Usage:
  iwr -useb https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1 | iex
#>

param(
  [ValidateSet("install","update","uninstall")]
  [string]$Action="install"
)

$Repo        = "Eterna-Inc/telegram-cli"
$AppName     = "telegram.exe"
$InstallDir  = "$Env:ProgramFiles\TelegramCLI"
$ConfigPath  = "$Env:USERPROFILE\.telegram_config"
$Arch        = if([Environment]::Is64BitOperatingSystem){"amd64"}else{"386"}
$OS          = "windows"

Write-Host "Telegram CLI Installer" -ForegroundColor Cyan

function Require-Admin {
  $id = [Security.Principal.WindowsIdentity]::GetCurrent()
  $p  = New-Object Security.Principal.WindowsPrincipal($id)
  if(-not $p.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)){
  Write-Host "Restarting as Administrator..."
    try {
      # If the script was invoked from a file, re-run that file elevated.
      if ($PSCommandPath -and (Test-Path $PSCommandPath)) {
        Start-Process powershell -Verb runAs -ArgumentList "-ExecutionPolicy Bypass -File `"$PSCommandPath`" -Action $Action"
      } else {
        # When the script was piped (iwr ... | iex), $PSCommandPath may be empty.
        # Download a temporary copy and run it elevated with the same action.
        $scriptUrl = "https://raw.githubusercontent.com/Eterna-Inc/telegram-cli/main/install.ps1"
        $tmp = Join-Path $env:TEMP "telegram-cli-install.ps1"
  Write-Host "Downloading installer to $tmp"
        Invoke-WebRequest -Uri $scriptUrl -OutFile $tmp -UseBasicParsing -ErrorAction Stop
        Start-Process powershell -Verb runAs -ArgumentList "-ExecutionPolicy Bypass -File `"$tmp`" -Action $Action"
      }
    } catch {
      Write-Host "Failed to restart elevated: $($_.Exception.Message)" -ForegroundColor Red
    }
    exit
  }
}

function Detect-LatestURL {
  param($repo,$os,$arch)
  $h=@{}
  if($env:GHTOKEN){$h["Authorization"]="token $env:GHTOKEN"}
  try{
    $r=Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -Headers $h -ErrorAction Stop
    $r.assets.browser_download_url | Where-Object {$_ -match "$os-$arch"}
  }catch{
    Write-Host "Cannot fetch release: $($_.Exception.Message)" -ForegroundColor Red
    return $null
  }
}

function Write-Config {
@'
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
'@ | Out-File -Encoding utf8 $ConfigPath
}

function Install-App {
  Require-Admin
  $url=Detect-LatestURL $Repo $OS $Arch
  if(-not $url){
    Write-Host "⚠️ No binary found for $OS-$Arch. Build manually with 'go build'." -ForegroundColor Yellow
    exit 1
  }
  if(-not(Test-Path $InstallDir)){New-Item $InstallDir -ItemType Directory |Out-Null}
  $dest=Join-Path $InstallDir $AppName
  Write-Host "Downloading latest release..."
  Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing
  Write-Host "Installed to $dest"

  # PATH
  $sysPath=[Environment]::GetEnvironmentVariable("Path",[EnvironmentVariableTarget]::Machine)
  if($sysPath -notmatch [Regex]::Escape($InstallDir)){
    [Environment]::SetEnvironmentVariable("Path","$sysPath;$InstallDir",[EnvironmentVariableTarget]::Machine)
    Write-Host "Added to PATH"
  }

  if(-not(Test-Path $ConfigPath)){
    Write-Host "Creating config..."
    Write-Config
  }

  Write-Host "`nTelegram CLI ready!"
  Write-Host "Run: telegram.exe \"Hello from Windows!\"`n"
}

function Update-App {
  Require-Admin
  Write-Host "Updating..."
  $t=Join-Path $InstallDir $AppName
  if(Test-Path $t){Remove-Item -Force $t}
  Install-App
}

function Uninstall-App {
  Require-Admin
  Write-Host "Uninstalling..."
  $t=Join-Path $InstallDir $AppName
  if(Test-Path $t){Remove-Item -Force $t}
  if(Test-Path $ConfigPath){Remove-Item -Force $ConfigPath}
  Write-Host "Uninstalled"
}

switch($Action){
  "install"   {Install-App}
  "update"    {Update-App}
  "uninstall" {Uninstall-App}
  default     {Write-Host "Usage: install.ps1 [install|update|uninstall]"}
}
