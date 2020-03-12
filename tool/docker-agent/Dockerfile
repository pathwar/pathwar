FROM pathwar/pathwar:latest as pathwar
FROM devago/docker-compose
COPY --from=pathwar /bin/pathwar /bin/
ENTRYPOINT ["/bin/pathwar"]
