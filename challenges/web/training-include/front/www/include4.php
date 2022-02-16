<p>
That's because the code that retrieves the page looks like this:
</p>

<pre>
<?php

show_source('index.php');
?>
</pre>

<p>Notice the <b>.php</b> suffix that's automatically added.</p>

<p>In order for our request to override it, a technique with <a href="http://php.net/manual/en/security.filesystem.nullbytes.php">null-byte terminators</a> is used, and we request <b>/index.php?page=../../../../../etc/passwd%00</b> instead.</p>

<p>Depending on your background, perhaps now you can envision a more nefarious file inclusion, e.g. a MySQL my.cnf config file.</p>

<p><b>Note:</b> The null-byte termination vulnerability was fixed in PHP 5.3.4. But since the Internet is full of deprecated PHP, this is not obsolete knowledge!</p>

<p><a href="/index.php?page=include5">Continue</a></p>
