<?php

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

session_start();

require_once ("LocaleManager.php");
$currentLocale = new LocaleManager();
?>
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title><?php echo _('Captcha Mario'); ?></title>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <header>
        <h1><?php echo _('Captcha Mario'); ?></h1>
    </header>
    <nav class="topright">
        <?php echo $currentLocale->get_locale_form(); ?>
    </nav>
    <section class="content">

<?php
  if (!isset($_SESSION['score'])) {
    $_SESSION['score'] = 0;
  }

  if (isset($_POST['password']) && isset($_SESSION['password'])) {
    if ($_POST['password'] == $_SESSION['password']) {
      $_SESSION['score']++;
      unset($_SESSION['password']);
      echo "<p>" . _('You win !') . "</p>";
    } else {
      echo "<p>" . _('You lose !') ."</p>";
    }
  }

  echo '<p>' . _('Score:') . ' ' . $_SESSION['score'], '/100';

  if ($_SESSION['score'] >= 100) {
    echo "<p>" . _('The passphrase is:') ." {passphrase}</p>";
  }

$_SESSION['password'] = make_password();

?>
<div>
<?php echo "<p>Captcha: ".$_SESSION['password']."</p>"; ?>
</div>
<div>
  <form method="POST">
    <input name="password" type="text" />
    <input type="submit" />
  </form>
</div>
