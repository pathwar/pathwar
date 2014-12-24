<?php

session_start();

function is_valid_char($c) {
  return (
    ($c > 47 && $c < 58) ||
    ($c > 64 && $c < 91) ||
    ($c > 96 && $c < 123)
  );
}

function make_password($length=10) {
  $str = '';
  for ($i = 0; $i < $length; $i++) {
    do {
       $c = rand() % 128;
    } while (!is_valid_char($c));
    $str .= chr($c);
  }
  return $str;
}

$_SESSION['password'] = make_password(3);

// Image creation
$img = imagecreate(300, 150);
$bgcolor = imagecolorallocate($img, rand() % 255, rand() % 255, rand() % 255);
$fgcolor = imagecolorallocate($img, rand() % 255, rand() % 255, rand() % 255);
imagestring($img, 5, rand() % 200, rand() % 100, $_SESSION['password'], $fgcolor);
header('Content-Type: image/png');
imagepng($img);
