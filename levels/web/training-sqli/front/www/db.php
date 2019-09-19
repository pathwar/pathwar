<?php
$mysqli = new mysqli('mysql:3306', 'training_sqli', 'training_sqli', 'training_sqli');

if (mysqli_connect_errno()) {
    die(mysqli_connect_error());
}
