<p>
That's because the code here's looking like this:
</p>

<pre>
<?php

show_source('index.php');
?>
</pre>

<p>So basically we would need to bypass the <b>.php</b> which is automatically happened.</p>
<p>In order to bypass it, a technique with <a href="http://php.net/manual/en/security.filesystem.nullbytes.php">null-byte terminators</a> is used.</p>

<p>So you would need to do something more like <b>/index.php?page=../../../../../etc/passwd%00</b>
<p>Now you can think about including a MySQL my.cnf config file and so on...</p>

<p><u>NB:</u> Unfortunately, this is no more exploitable since PHP 5.3.4. But since the Internet is full of deprecated and old versions of PHP, it's still something interesting to check!</p>


<p><A href="/index.php?page=include5">Want to know more ?</a></p>