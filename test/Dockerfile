FROM            alpine
RUN             apk --no-cache add curl openssl python3 make bash jq \
 &&             pip3 install -U pip setuptools httpie \
 &&             http --version \
 &&             jq --version
ENTRYPOINT      ["/bin/bash", "-c"]
CMD             ["sleep 5; ./integration.sh"]
COPY            . .
