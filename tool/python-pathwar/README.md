python-pathwar
==============

Pathwar API client in Python

Examples
--------

```python
from pathwar import Pathwar

## Connection
client = Pathwar(user="username", password="password")  # Connect with user+pass couple
client = Pathwar(token="token")  # Connect using generated token

## Switch organization
client.set_organization('roxxorz')  # Switch to organization with name = `roxxorz`
client.set_organization('abcdef-ghij-klmn-opqr-stuvwx')  # Switch to organization by `organization_id`
client.set_organization(session='super-final')  # switch to the organization where session is `super-final`

## Some actions
client.levels.get(name='pnu').buy()
client.coupons.validate('cool-coupon')

## A level workflow
for level in client.levels.all():
    if not level.validated:
        print('You still need to validate the level {}'.format(level.name))
        # Try to use 'toto' as passphrase
        level.validate('toto')
    elif not level.bought:
        # Try to buy the level
        level.buy()
```
