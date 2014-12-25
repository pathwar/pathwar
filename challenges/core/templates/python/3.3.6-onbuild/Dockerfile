# pathwar/python:3.3.6-onbuild pathwar/python:3.3-onbuild
FROM python:3.3.6
MAINTAINER Pathwar Team <team@pathwar.net> (@pathwar_net)

# Add Pathwar tools
COPY skeleton /pathwar

# Override pathwar scripts
RUN echo python3 main.py > /pathwar/scripts/run

# Onbuild
ONBUILD COPY scripts /pathwar/scripts/
ONBUILD COPY level.yml /pathwar/
ONBUILD CMD /pathwar/scripts/init

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ONBUILD COPY app/requirements.txt /usr/src/app/
ONBUILD RUN pip install -r requirements.txt
ONBUILD COPY app /usr/src/app
