# -*- coding: utf-8 -*-

import urllib

import requests
from flask.ext.mail import Mail, Message

mail = Mail()


def send_mail(recipients, subject, message):
    with mail.connect() as conn:
        for recipient in recipients:
            to = '{} <{}>'.format(urllib.quote_plus(recipient['login']), recipient['email'])
            msg = Message(
                body=message,
                subject=subject,
                sender=("Pathwar bot", "team@pathwar.net"),
                recipients=[to],
            )
        conn.send(msg)
