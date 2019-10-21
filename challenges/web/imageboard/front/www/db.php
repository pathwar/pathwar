<?php
$mysqli = new mysqli('mysql:3306', 'imageboard', 'imageboard', 'imageboard');

if (mysqli_connect_errno()) {
    die(mysqli_connect_error());
}
