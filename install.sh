#!/bin/sh

## USAGE: bash install.sh $INSTALL_DIR
set -eu

# OS detection
KERNEL="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH=$(uname -m)

if [ "$KERNEL" = "darwin" ]; then
  if [ "$ARCH" = "x86_64" ]; then
    PLATFORM='darwin_amd64'
  else
    PLATFORM='darwin_386'
  fi
elif [ "$KERNEL" = "linux" ]; then
  if [ "$ARCH" = "x86_64" ]; then
    PLATFORM='linux_amd64'
  else
    PLATFORM='linux_386'
  fi
fi

if [ "$0" = 'sh' ]; then
  # I.e. when piping from curl
  INSTALL_DIR=${HOME}/.taxonkit
elif [ "$0" = 'install.sh' ]; then
  # I.e. when running directly
  INSTALL_DIR=${HOME}/.taxonkit
else
  # when an argument is passed
  INSTALL_DIR="$0"
fi

DOWNLOAD_URL=$(curl -s https://api.github.com/repos/shenwei356/taxonkit/releases/latest \
  | grep browser_download_url \
  | grep -i $PLATFORM \
  | cut -d '"' -f 4)

echo >&2 "==> Installing TaxonKit to:"
echo >&2 "    ${INSTALL_DIR}"
echo >&2

mkdir -p "${INSTALL_DIR}/bin"
curl -SL "$DOWNLOAD_URL" | tar zxf - -C "${INSTALL_DIR}/bin"

echo >&2
echo >&2 "==> TaxonKit successfully installed."

echo >&2
echo >&2 "==> Downloading Taxonomy Dump from NCBI"
echo >&2

TAX_DUMP="ftp://ftp.ncbi.nih.gov/pub/taxonomy/taxdump.tar.gz"

curl -SL "${TAX_DUMP}" | tar zxf - -C "${INSTALL_DIR}"

### Check which SHELL and then test different profile files
case $SHELL in
*/zsh)
  # assume Zsh
  if test -e "${HOME}/.zshrc"; then
    DOT_FILE=${HOME}/.zshrc
  elif test -e "${HOME}/.zprofile"; then
    DOT_FILE=${HOME}/.zprofile
  elif test -e "${HOME}/.profile"; then
    DOT_FILE=${HOME}/.profile
  fi
  ;;
*/bash)
  # assume Bash
  if test -e "${HOME}/.bashrc"; then
    DOT_FILE=${HOME}/.bashrc
  elif test -e "${HOME}/.bash_profile"; then
    DOT_FILE=${HOME}/.bash_profile
  elif test -e "${HOME}/.profile"; then
    DOT_FILE=${HOME}/.profile
  fi
  ;;
*)
  if test -e "${HOME}/.profile"; then
    DOT_FILE=${HOME}/.profile
  fi
esac

if [ -z ${DOT_FILE+x} ]; then
  # DOT File hasn't been set.
  echo >&2
  echo >&2 '==> No profile files were found.'
  echo >&2 '    Please create one and add the following line to that file:'
  echo >&2
  echo >&2 '    export PATH="'"${INSTALL_DIR}"'/bin:${PATH}"'
else
  echo >&2 'export PATH="'"${INSTALL_DIR}"'/bin:${PATH}"' >> "${DOT_FILE}"
  echo >&2
  echo >&2 "==> Added TaxonKit to your PATH in ${DOT_FILE}"
  echo >&2
  echo >&2 "==> Run \`taxonkit -h\` in a new window to get started."
fi

echo >&2
