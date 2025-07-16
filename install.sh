#!/bin/bash

# This script installs rssx on Linux and macOS.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/l-z-h/rssx/main/scripts/install.sh | sh
#
# The script will:
# 1. Detect the user's OS and architecture.
# 2. Download the latest release of rssx for their platform.
# 3. Extract the binary and move it to /usr/local/bin (with sudo) or ~/.local/bin.
# 4. If installed to ~/.local/bin, it will add the directory to the user's PATH.

set -e

# --- Helper Functions ---

# Colors for output
red='\033[0;31m'
nc='\033[0m' # No Color
yellow='\033[0;33m'
green='\033[0;32m'

# --- Main Script ---

main() {
  # Get OS and architecture
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  arch=$(uname -m)

  case $arch in
    x86_64) arch="amd64" ;;
    aarch64) arch="arm64" ;;
    armv*) arch="arm64" ;;
  esac

  # Get the latest version tag from the GitHub releases page.
  latest_version_url=$(curl -sL -o /dev/null -w %{url_effective} https://github.com/lakerszhy/rssx/releases/latest)
  latest_version=$(basename "$latest_version_url")

  if [ -z "$latest_version" ]; then
    echo -e "${red}❌ Failed to fetch the latest version of rssx.${nc}"
    exit 1
  fi

  # Construct the download URL
  download_url="https://github.com/lakerszhy/rssx/releases/download/${latest_version}/rssx-${latest_version}-${os}_${arch}.tar.gz"

  echo -e "${green}⬇️  Downloading rssx ${latest_version} for ${os}/${arch}...${nc}"

  # Create a temporary directory for the download
  temp_dir=$(mktemp -d)
  trap 'rm -rf "$temp_dir"' EXIT

  # Download and extract the tarball
  if ! curl -L "$download_url" | tar -xz -C "$temp_dir"; then
    echo -e "${red}❌ Failed to download or extract rssx. Please check the URL and your network connection.${nc}"
    exit 1
  fi

  echo -e "${green}✅ Download complete.${nc}"

  # Define the directory name inside the archive
  archive_dir="rssx-${latest_version}-${os}_${arch}"

  # Move the binary to the installation directory
  # Try /usr/local/bin first, then fallback to ~/.local/bin
  if [ -w /usr/local/bin ]; then
    install_dir="/usr/local/bin"
    echo -e "${yellow}Installing rssx to ${install_dir}...${nc}"
    if mv "${temp_dir}/${archive_dir}/rssx" "${install_dir}/"; then
      echo -e "${green}✅ Successfully installed to ${install_dir}${nc}"
    else
      echo -e "${yellow}⚠️  Failed to install to ${install_dir}, trying ~/.local/bin...${nc}"
      install_dir="$HOME/.local/bin"
      mkdir -p "${install_dir}"
      if ! mv "${temp_dir}/${archive_dir}/rssx" "${install_dir}/"; then
        echo -e "${red}❌ Failed to install rssx to both locations.${nc}"
        exit 1
      fi
    fi
  else
    install_dir="$HOME/.local/bin"
    echo -e "${yellow}Installing rssx to ${install_dir}...${nc}"
    mkdir -p "${install_dir}"
    if ! mv "${temp_dir}/${archive_dir}/rssx" "${install_dir}/"; then
      echo -e "${red}❌ Failed to install rssx.${nc}"
      exit 1
    fi
  fi

  # Add to PATH if installed to ~/.local/bin and not already in PATH
  if [[ "$install_dir" == *".local/bin" ]] && ! [[ ":$PATH:" == *":${install_dir}:"* ]]; then
    echo -e "${yellow}Adding ${install_dir} to your PATH.${nc}"
    # Detect shell and update config
    shell_config_file=""
    case "$SHELL" in
      */bash) shell_config_file="~/.bashrc" ;;
      */zsh) shell_config_file="~/.zshrc" ;;
      */fish) shell_config_file="~/.config/fish/config.fish" ;;
      *) echo -e "${red}Unsupported shell: ${SHELL}. Please add ${install_dir} to your PATH manually.${nc}" ;;
    esac

    if [ -n "$shell_config_file" ]; then
      if [[ "$SHELL" == *"/fish" ]]; then
        echo -e "\nfish_add_path ${install_dir}" >> "$shell_config_file"
      else
        echo -e "\nexport PATH=\"${install_dir}:\$PATH\"" >> "$shell_config_file"
      fi
      echo -e "${green}Please restart your shell or run 'source ${shell_config_file}' to apply the changes.${nc}"
    fi
  fi

  echo -e "${green}🎉 Installation complete!${nc}"
  echo -e "You can now run 'rssx' to start the application."
}

main