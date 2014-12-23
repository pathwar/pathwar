#!/bin/sh

run_script() {
  script="${1}"
  shift
  args="${@}"
  if [ -f "/pathwar/scripts/${script}" ]; then
      echo "[+] Running '${script}' script (args: ${args}"
      /pathwar/scripts/${script} ${args}
  else
      echo "[-] Script '${script}' not found"
  fi
}
