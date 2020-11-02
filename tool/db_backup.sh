echo Backup of pathwar db
_now=$(date +"%d%m%Y_%H%M")
_file="pathwar.$_now.sql"
/usr/bin/docker exec platform-dev_db_1 mysqldump -u root -p'uns3cur3' pathwar --skip-extended-insert > /home/$(whoami)/pwdb_backup/$_file
while IFS= read -r f; do rm "/home/$(whoami)/pwdb_backup/$f"; done < <(ls -tp /home/$(whoami)/pwdb_backup/ | grep -v '/$' | tail -n +16)
