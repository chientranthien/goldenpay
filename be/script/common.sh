echo_info() {
  color=$(tput setaf 2)
  reset=$(tput sgr0)
  echo "${color}[INFO] $*${reset}"
}

echo_warn() {
  color=$(tput setaf 3)
  reset=$(tput sgr0)
  echo "${color}[WARN] $*${reset}"
}

echo_error() {
  # Visit this page for tput look up color code
  # https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
  color=$(tput setaf 208)
  reset=$(tput sgr0)
  echo "${color}[ERROR] $*${reset}"
}

echo_finish() {
 echo_info "===================================================================="
}
