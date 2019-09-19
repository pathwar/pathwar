<?php
if (isset($_GET['step'])) {
    switch ($_GET['step']) {
    case '1':
        include('step1.php'); // HTTP Methods
        break;
    case '2':
        include('step2.php'); // HTTP Headers
        break;
    case '3':
        include('step3.php'); // HTTP Redirection
        break;
    case '4':
        include('step4.php'); // HTTP Cookies
        break;
    default:
        include('default.php');
        break;
    }
}
else
	include('default.php');