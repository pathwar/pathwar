<p>Some PHP configurations even allow URL include...</p>

<p>We're going to include a file called <b>pown.php</b> which is composed of the following code piece:</p>

<pre>
<?php
show_source('pown.php');
?>
</pre>

<p>So if you click <a href="/index.php?page=<?php echo "http://".$_SERVER['SERVER_NAME']."/pown"; ?>" target="_blank">here</a>, you'll be able to list all the file in this directory</p>

<p>Impressed ? That's nothing compare what you're going to learn after ! :-)</p>

<p><b>Happy Hacking !</b></p>