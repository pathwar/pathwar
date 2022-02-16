<p>Some PHP configurations even allow URL include...</p>

<p>We're going to include a file called <b>pown.php</b>:</p>

<pre>
<?php
show_source('pown.php');
?>
</pre>

<p>If you click <a href="/index.php?page=<?php echo "http://".$_SERVER['SERVER_NAME']."/pown"; ?>" target="_blank">here</a>, you'll be able to list all the files in this directory!</p>

<p>Impressed? This is nothing compared to what you'll learn later!</p>

<p><b>Happy Hacking!</b></p>