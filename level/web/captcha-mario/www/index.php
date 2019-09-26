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


  if (!isset($_SESSION['score'])) {
    $_SESSION['score'] = 0;
  }

  if (isset($_POST['password']) && isset($_SESSION['password'])) {
    if ($_POST['password'] == $_SESSION['password']) {
      $_SESSION['score']++;
      unset($_SESSION['password']);
      echo "<p>You win !</p>";
    } else {
      echo "<p>You lose !</p>";
    }
  }

  echo '<p>Score: ', $_SESSION['score'], '/1000';

  if ($_SESSION['score'] >= 1000) {
    echo "<p>The passphrase is: ", file_get_contents('/tmp/passphrase.txt'), "</p>";
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
