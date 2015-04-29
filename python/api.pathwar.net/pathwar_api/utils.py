# -*- coding: utf-8 -*-

import random

from flask import current_app


def request_get_user(request):
    auth = request.authorization
    if auth.get('username'):
        if auth.get('password'):  # user:pass
            return current_app.data.driver.db['raw-users'].find_one({
                'login': auth.get('username'),
                # 'active': True,  # FIXME: Reenable later
            })
        else:  # token
            user_token = current_app \
                .data \
                .driver \
                .db['raw-user-tokens'] \
                .find_one({
                    'token': auth.get('username')
                })
            if user_token:
                return current_app \
                    .data \
                    .driver \
                    .db['raw-users'] \
                    .find_one({'_id': user_token['user']})
    return None


def generate_name():
    """ Python port of https://github.com/docker/docker/blob/master/pkg/namesgenerator/names-generator.go.
    """
    left = [
        "admiring",
        "adoring",
        "agitated",
        "angry",
        "backstabbing",
        "berserk",
        "boring",
        "clever",
        "cocky",
        "compassionate",
        "condescending",
        "cranky",
        "desperate",
        "determined",
        "distracted",
        "dreamy",
        "drunk",
        "ecstatic",
        "elated",
        "elegant",
        "evil",
        "fervent",
        "focused",
        "furious",
        "gloomy",
        "goofy",
        "grave",
        "happy",
        "high",
        "hopeful",
        "hungry",
        "insane",
        "jolly",
        "jovial",
        "kickass",
        "lonely",
        "loving",
        "mad",
        "modest",
        "naughty",
        "nostalgic",
        "pensive",
        "prickly",
        "reverent",
        "romantic",
        "sad",
        "serene",
        "sharp",
        "sick",
        "silly",
        "sleepy",
        "stoic",
        "stupefied",
        "suspicious",
        "tender",
        "thirsty",
        "trusting",
    ]
    right = [
        "albattani",
        "almeida",
        "archimedes",
        "ardinghelli",
        "babbage",
        "banach",
        "bardeen",
        "brattain",
        "shockley",
        "bartik",
        "bell",
        "blackwell",
        "bohr",
        "brown",
        "carson",
        "colden",
        "cori",
        "cray",
        "curie",
        "darwin",
        "davinci",
        "einstein",
        "elion",
        "engelbart",
        "euclid",
        "fermat",
        "fermi",
        "feynman",
        "franklin",
        "galileo",
        "goldstine",
        "goodall",
        "hawking",
        "heisenberg",
        "hodgkin",
        "hoover",
        "hopper",
        "hypatia",
        "jang",
        "jones",
        "kilby",
        "noyce",
        "kirch",
        "kowalevski",
        "lalande",
        "leakey",
        "lovelace",
        "lumiere",
        "mayer",
        "mccarthy",
        "mcclintock",
        "mclean",
        "meitner",
        "mestorf",
        "morse",
        "newton",
        "nobel",
        "payne",
        "pare",
        "pasteur",
        "perlman",
        "pike",
        "poincare",
        "poitras",
        "ptolemy",
        "ritchie",
        "thompson",
        "rosalind",
        "sammet",
        "sinoussi",
        "stallman",
        "swartz",
        "tesla",
        "torvalds",
        "turing",
        "wilson",
        "wozniak",
        "wright",
        "yalow",
        "yonath",
    ]
    return random.choice(left) + '_' + random.choice(right)
