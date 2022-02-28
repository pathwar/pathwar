<?php
echo ("<p>" . _('Some PHP configurations even allow URL include...') ."</p>"
     . "<p>" . _('We\'re going to include a file called') . " <b>pown.php</b>:</p>");
?>

<pre>
<?php
show_source('pown.php');
?>
</pre>

<?php
echo ( "<p>" ._('If you click') . ' <a href="/index.php?page=http://' .$_SERVER['SERVER_NAME'] . '/pown" target="_blank">' . _('here') . "</a>," 
        . _('you\'ll be able to list all the files in this directory!') ."</p>"
        . "<p>" . _('Impressed? This is nothing compared to what you\'ll learn later!') . "</p>"
        . "<p><b>" . _('Happy Hacking!') ."</b></p>");