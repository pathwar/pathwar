FROM            alpine
RUN             apk --no-cache add curl openssl python3 make bash jq \
 &&             pip3 install -U pip setuptools httpie \
 &&             http --version \
 &&             jq --version
RUN             wget https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh && chmod +x wait-for-it.sh
ENTRYPOINT      ["/bin/bash", "-c"]
COPY            . .
