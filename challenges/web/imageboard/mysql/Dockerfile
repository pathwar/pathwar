FROM            mysql:5.7
COPY            bootstrap.sql /docker-entrypoint-initdb.d/
ENV             MYSQL_DATABASE=imageboard \
                MYSQL_USER=imageboard \
                MYSQL_PASSWORD=imageboard \
                MYSQL_RANDOM_ROOT_PASSWORD=yes
