#FROM pathwar/golang
FROM pathwar/golang
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

RUN echo revel run github.com/revel/revel/samples/chat > /pathwar/scripts/run
CMD /pathwar/scripts/init
COPY scripts /pathwar/scripts/
COPY level.yml /pathwar/

RUN go get github.com/revel/cmd/revel
EXPOSE 9000
