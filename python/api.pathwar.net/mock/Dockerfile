FROM python:2-onbuild
RUN cd /root; pip install -e git://github.com/nicolaiarocci/eve.git@develop#egg=eve
RUN cd /root; pip install fake-factory factory_boy
CMD python run.py
EXPOSE 5000
