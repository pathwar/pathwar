import requests

def send_mail(user, message):
    to = '{} <{}>'.format(user['login'], user['email'])
    print(to, message)

    return
    return requests.post(
        "https://api.mailgun.net/v2/sandboxa094db996c974fc7aaf47b5cd4f45d82.mailgun.org/messages",
        auth=(
            "api", "key-bb8419945197f59c961659c8c7fd7547",
        ),
        data={
            "from": "Mailgun Sandbox <postmaster@sandboxa094db996c974fc7aaf47b5cd4f45d82.mailgun.org>",
            "to": to,
            "subject": "Hello Manfred Touron",
            "text": message,
        })
