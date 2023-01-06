#!/bin/sh
set -e


# settings
owner="carbonetes"
repo="diggity"
version=""
githubUrl="https://github.com"
executable_folder="./bin"
format="tar.gz"

usage() (
  this=$1
  cat <<EOF
$this: download go binaries for carbonetes/diggity

Usage: $this [-d] dir [-v] [tag]
  -d  the installation directory (dDefaults to ./bin)
  -v the specific release to use (if missing, then the latest will be used)
EOF
  exit 2
)

get_arch() {
    a=$(uname -m)
    case ${a} in
        "x86_64" | "amd64" )
            echo "amd64"
        ;;
        "i386" | "i486" | "i586")
            echo "386"
        ;;
        "aarch64" | "arm64" | "arm")
            echo "arm64"
        ;;
        "mips64el")
            echo "mips64el"
        ;;
        "mips64")
            echo "mips64"
        ;;
        "mips")
            echo "mips"
        ;;
        "ppc64le")
            echo "ppc64le"
        ;;
        "s390x")
            echo "s390x"
        ;;
        *)
            echo ${NIL}
        ;;
    esac
}
get_binary_name() {
  os="$1"
  arch="$2"
  binary="$3"
  original_binary="${binary}"

  case "$1" in
    windows) binary="$3.exe" ;;
  esac

  echo "get_binary_name(os=${os}, arch=${arch}, binary=${original_binary}) returned '${binary}'"

  echo "${binary}"
}

get_os(){
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case "$os" in
        cygwin_nt*) os="windows" ;;
        mingw*) os="windows" ;;
        msys_nt*) os="windows" ;;
    esac
    echo "$os"
}
get_latest_release() {
    curl --silent "https://api.github.com/repos/$1/$2/releases/latest" |
    grep '"tag_name":' |                                           
    sed -E 's/.*"([^"]+)".*/\1/'                                    
}
install_binary() (

  # don't continue if we don't have anything to install
  if [ -z "$1" ]; then
      return
  fi

  archive_dir=$(dirname "$1")

  # unarchive the downloaded archive to the temp dir
  (cd "${archive_dir}" && extract "$1")
  # create the destination dir
  test ! -d "$3" && install -d "$2"

  # install the binary to the destination dir
  install "${archive_dir}/$3" "$2"
)

extract() (
  archive=$1
  case "$1" in
    *.tar.gz | *.tgz) tar --no-same-owner -xzf "$1" ;;
    *.tar) tar --no-same-owner -xf "$1" ;;
    *.zip) unzip -q "$1" ;;
    *.dmg) extract_from_dmg "$1" ;;
    *)
      echo "erorr extracting unknown archive format for $1"
      return 1
      ;;
  esac
)

install_diggity() {
    # parse flag
    while getopts "v:d:" arg; do
        case "${arg}" in
            d) executable_folder="$OPTARG";;
            v) version="$OPTARG";;
        esac
    done
    shift $((OPTIND - 1))
    set +u

    
    downloadFolder=$(mktemp -d -t diggity-XXXXXXXXXX)
    trap 'rm -rf -- "$downloadFolder"' EXIT
    mkdir -p ${downloadFolder} # make sure download folder exists
    os=$(get_os)
    arch=$(get_arch)
    # if version is empty
    if [ -z "$version" ]; then
        tag=$(get_latest_release ${owner} ${repo})
        version=${tag}
    fi
    
    # change format to .zip if windows
    # append .exe if windows
    final_binary="${repo}"
    case ${os} in
     windows) format=zip final_binary="${repo}.exe";;
    esac

    # init filename for binary
    file_name="${repo}_${version#v}_${os}_${arch}.${format}"
    downloaded_file="${downloadFolder}/${file_name}"
    asset_uri="${githubUrl}/${owner}/${repo}/releases/download/${version}/${file_name}"

    echo "[1/3] Download ${asset_uri} to tmp folder"
    rm -f ${downloaded_file}
    curl --fail --location --output "${downloaded_file}" "${asset_uri}"

    echo "[2/3] Install ${repo} to the ${executable_folder}"

    install_binary "${downloaded_file}" "${executable_folder}" "${final_binary}"
    exe=${executable_folder}/${repo}
    chmod +x ${exe}

    echo "[3/3] Set environment variables"
    echo "${repo} was installed successfully to ${exe}"
    if command -v $repo --version >/dev/null; then
        echo "Run '$repo --help' to get started"
    else
        echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
        echo "  export PATH=${executable_folder}:\$PATH"
        echo "Run '$exe_name --help' to get started"
    fi
}

install_diggity "$@"

# exit 0