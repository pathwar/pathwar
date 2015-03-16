# -*- coding: utf-8 -*-


import requests
import slumber


class BaseSDK(object):

    def __init__(self, auth_token=None, base_url=None, verify_ssl=True):
        self.auth_token = auth_token
        self.base_url = base_url
        self.verify_ssl = verify_ssl

    def make_requests_session(self):
        """ Attaches an Authorization header to requests.Session. """
        session = requests.Session()
        if self.auth_token:
            session.headers.update({
                # 'Authorization': base64 ':{}'.format(self.auth_token),
            })
        if not self.verify_ssl:
            session.verify = False
        return session

    def query(self):
        """ Gets a configured slumber.API object. """
        return slumber.API(
            self.base_url,
            session=self.make_requests_session(),
        )


class PathwarSDK(BaseSDK):
    """ Interacts with Pathwar API. """

    base_url = 'https://api.pathwar.net'
