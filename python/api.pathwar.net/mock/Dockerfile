FROM python:2-onbuild

RUN cd /root; pip install -e git+https://github.com/hellerstanislav/eve-mongoengine#egg=eve-mongoengine
RUN cd /root; pip install -e git://github.com/nicolaiarocci/eve.git@develop#egg=eve


CMD python run.py
EXPOSE 5000
