#!/bin/sh

# Scripts/hooks
run_script() {
    script="${1}"
    shift
    args="${@}"
    if [ -f "/pathwar/scripts/${script}" ]; then
        if [ $# -gt 0 ]; then
            echo "[+] Running '${script}' script (args: ${args})"
        else
            echo "[+] Running '${script}' script"
        fi
        /pathwar/scripts/${script} ${args}
    else
        echo "[-] Script '${script}' not found"
    fi
}

# Passphrase helpers
get_passphrase() {
    key="${1}"
    if [ -f /pathwar/passphrases/${key} ]; then
        cat /pathwar/passphrases/${key}
    else
        warn "No such passphrase: ${key}"
        echo "NO_SUCH_PASSPHRASE"
    fi
}

# Logging
debug() {
    echo "[?]   $@" >&2
}

info() {
    echo "[+]   $@" >&2
}

warn() {
    echo "[-]   $@" >&2
}

fatal() {
    echo "[-]   $@" >&2
    exit -1
}
