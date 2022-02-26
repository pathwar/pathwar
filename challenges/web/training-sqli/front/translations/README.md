# How translate

## If original strings have been changed.
Update the translate template (.pot file)

    xgettext -n ../www/*.php -s -o messages.pot

## Translate
Open templates/messages.pot with a compatible translation tool. Poedit in my case.

Do the translation and save it on the file:

    translations/locale/LC_messages/messages.po

For the en_US translation, it is

    translations/en_US/LC_MESSAGES/messages.po

The binary translation file (.mo) will be automatically generated when the challenge instance is inited.


