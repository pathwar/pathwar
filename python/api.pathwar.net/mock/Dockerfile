FROM python:2-onbuild
RUN cd /root; pip install -e git://github.com/nicolaiarocci/eve.git@develop#egg=eve
RUN cd /root; pip install -e git://github.com/charlesflynn/eve-docs.git@master#egg=eve_docs
CMD python run.py
EXPOSE 5000
