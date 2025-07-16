# RssX Installer for Windows

# This script installs rssx on Windows using PowerShell.
#
# Usage:
#   irm https://raw.githubusercontent.com/lakerszhy/rssx/main/scripts/install.ps1 | iex
#
# The script will:
# 1. Detect the user's architecture (amd64 or arm64).
# 2. Download the latest release of rssx for Windows.
# 3. Create a directory at '$HOME\AppData\Local\rssx'.
# 4. Extract the binary to the new directory.
# 5. Add the directory to the user's PATH environment variable.

param()

$ErrorActionPreference = 'Stop'

function Main {
    # --- Detect Architecture ---
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        'AMD64' { 'amd64' }
        'ARM64' { 'arm64' }
        default { Write-Host "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" -ForegroundColor Red; exit 1 }
    }

    # --- Get Latest Version ---
    try {
        $response = Invoke-WebRequest -Uri "https://github.com/lakerszhy/rssx/releases/latest"
        $latest_version_url = $response.BaseResponse.ResponseUri.AbsoluteUri
        $latest_version = $latest_version_url.Split('/')[-1]
    } catch {
        Write-Host "Failed to fetch the latest version of rssx." -ForegroundColor Red
        exit 1
    }

    if (-not $latest_version) {
        Write-Host "Failed to fetch the latest version of rssx." -ForegroundColor Red
        exit 1
    }

    # --- Download and Extract ---
    $download_url = "https://github.com/lakerszhy/rssx/releases/download/$latest_version/rssx-$latest_version-windows_$($arch).zip"
    $install_dir = "$env:LOCALAPPDATA\rssx"
    $temp_zip_path = "$env:TEMP\rssx.zip"

    Write-Host "Downloading rssx $latest_version for Windows/$arch..." -ForegroundColor Green

    try {
        Invoke-WebRequest -Uri $download_url -OutFile $temp_zip_path
    } catch {
        Write-Host "Failed to download rssx. Please check the URL and your network connection." -ForegroundColor Red
        exit 1
    }

    Write-Host "Download complete." -ForegroundColor Green

    # --- Installation ---
    $temp_extract_dir = Join-Path -Path $env:TEMP -ChildPath "rssx_extracted"
    if (Test-Path -Path $temp_extract_dir) {
        Remove-Item -Path $temp_extract_dir -Recurse -Force
    }
    New-Item -ItemType Directory -Path $temp_extract_dir | Out-Null

    Write-Host "Extracting files..." -ForegroundColor Yellow
    Expand-Archive -Path $temp_zip_path -DestinationPath $temp_extract_dir -Force

    $archive_dir_name = "rssx-$latest_version-windows_$($arch)"
    $source_exe_path = Join-Path -Path $temp_extract_dir -ChildPath "$archive_dir_name\rssx.exe"

    if (-not (Test-Path -Path $install_dir)) {
        New-Item -ItemType Directory -Path $install_dir | Out-Null
    }

    Write-Host "Installing rssx to $install_dir..." -ForegroundColor Yellow
    Move-Item -Path $source_exe_path -Destination $install_dir -Force

    # --- Add to PATH ---
    $currentUserPath = [System.Environment]::GetEnvironmentVariable('Path', 'User')
    if (-not ($currentUserPath -split ';' -contains $install_dir)) {
        Write-Host "Adding $install_dir to your PATH." -ForegroundColor Yellow
        $newPath = "$currentUserPath;$install_dir"
        [System.Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        $env:Path = $newPath # Update for current session
    }

    # --- Cleanup ---
    Remove-Item -Path $temp_zip_path
    Remove-Item -Path $temp_extract_dir -Recurse -Force

    Write-Host "Installation complete!" -ForegroundColor Green
    Write-Host "Please restart your terminal for the PATH changes to take full effect."
    Write-Host "You can now run 'rssx' to start the application."
}

Main