FROM            mysql:5.7
COPY            bootstrap.sql /docker-entrypoint-initdb.d/
ENV             MYSQL_DATABASE=training_sqli \
                MYSQL_USER=training_sqli \
                MYSQL_PASSWORD=training_sqli \
                MYSQL_RANDOM_ROOT_PASSWORD=yes
COPY            on-init /pwinit/