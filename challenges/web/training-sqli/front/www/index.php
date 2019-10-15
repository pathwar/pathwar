<?php
include('header.php');
if (isset($_GET['step'])) {
    switch ($_GET['step']) {
	case '1':
	    include('step1.php');
	    break;
	case '2':
	    include('step2.php');
	    break;
	default:
	    include('default.php');
	    break;
    }
} else {
    include('default.php');
}
?>
</html>
