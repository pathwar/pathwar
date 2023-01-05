FROM python:3.8-alpine

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

COPY    www/ /var/www/html/
COPY    on-init /pwinit/
WORKDIR /var/www/html/

ENTRYPOINT [ "python" ]

CMD [ "app.py" ]