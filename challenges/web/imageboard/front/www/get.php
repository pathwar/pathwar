<?php
if (isset($_GET['file'])) {
	header('Content-Type: image/jpeg');
	header("Content-Transfer-Encoding: binary");
	readfile('/var/www/html/images/'.$_GET['file']);
}
?>
