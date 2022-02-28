<?php
echo ("<p>" . _('By browsing the awesome <a href="https://www.owasp.org/index.php/Testing_for_Local_File_Inclusion">OWASP</a> website, we learn that:') ."</p>"
     ."<p>" . _('The File Inclusion vulnerability allows an attacker to include a file, usually exploiting a "dynamic file inclusion" mechanisms implemented in the target application. The vulnerability occurs due to the use of user-supplied input without proper validation.') . "</p>"
     ."<p>" . _('This can allow attackers to simply output the contents of the file onto a page, or, more seriously, enable:') . "</p>"
    . "<ul> "
    . "<li>" . _('Code execution on the web server') . "</li>"
    . "<li>" . _('Code execution on the client-side such as JavaScript which can lead to other attacks such as cross site scripting (XSS)') ."</li>"
    . "<li>" . _('Denial of Service (DoS)') . "</li>"
    . "<li>" . _('Sensitive Information Disclosure') . "</li>"
    . "</ul>"
    . "<p><a href=\"/index.php?page=include3\">" . _('Continue...') ."</a></p>");