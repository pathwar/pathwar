<?php
  session_start();
  require_once ("LocaleManager.php");
  $currentLocale = new LocaleManager();
  ?>
  <!DOCTYPE html>
  <html>
  <head>
      <meta charset="UTF-8">
      <title><?php echo _('Captcha Luigi'); ?></title>
      <link rel="stylesheet" href="/style.css">
  </head>
  <body>
      <header>
          <h1><?php echo _('Captcha Luigi'); ?></h1>
      </header>
      <nav class="topright">
          <?php echo $currentLocale->get_locale_form(); ?>
      </nav>
      <section class="content">
      <article>

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

  echo '<p>' . _('Score:') . ' ', $_SESSION['score'], '/100';

  if ($_SESSION['score'] >= 100) {
    echo "<p>" . _('The passphrase is:') . " {passphrase}</p>";
  }
?>
  <p><img src="./captcha.php" /></p>
</article>
<article>
  <form method="POST">
    <input name="password" type="text" />
    <input type="submit" />
  </form>
</article>
</section>
</body>
</html>
