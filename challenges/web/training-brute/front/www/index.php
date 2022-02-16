<DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
  <html xmlns="http://www.w3.org/1999/xhtml" xml:lang="fr" >
    <head>
      <title>Training Bruteforce</title>
      <meta http-equiv="Content-Type" content="text/html; charset=UTF8" />
      <link rel="stylesheet" media="screen" type="text/css" title="design" href="design.css" />
    </head>
    <body>
      <div id="main">
    <h1>Bruteforce</h1>

    <p>
      The goal of this training is to learn how to write a simple brute force script.
    </p>

        <h2>You shall not pass!</h2>
        <p>
          <form name="bruteforce" action="index.php" method="get">
            Password <input type="text" name="pass" />
            <input type="submit" value="go" />
          </form>
        </p>
        <?php if (isset($_GET['pass'])) : ?>
        <p>
          <center>
            <?php if ($_GET['pass'] == md5("__RANDNUM0__")) : ?>
            Wisely done! The passphrase is <b>__PASSPHRASE0__</b>
            <?php else : ?>
            <img width='600px' src="./gandalf.jpg" alt="u_shall_not_pass" />
            <br />
        Invalid password, young mage.
            <?php endif ?>
          </center>
        </p>
        <?php endif ?>


        <h2>Tools</h2>
        <ol>
      <li><b>Network request tool.</b> For this training, you'll have to make several HTTP requests. curl can be handy, but you can use any command or tool you wish.</li>
      <li><b>Programmatic iteration.</b> At the core of a brute-force attack is repetition (think battering ram). Here's how to perform a loop in bash:

      <pre class="code_block"><?php echo htmlentities(file_get_contents('source_0.txt')); ?></pre>
          <br />
          Warning: <i>echo $i</i> is not the same as <i>echo -n $i</i>.</li>

      <li>You'll also need to generate a MD5 hash from a number (man md5sum, <a href="http://www.php.net/manual/en/function.md5.php">php.net</a>).</li>
        </ol>

        <h2>The Form</h2>
        <p>
      <pre class="code_block"><?php echo htmlentities(file_get_contents('source_1.txt')); ?></pre>
          <br />
      Take note of the important bits:
          <ul>
            <li>The action field, which indicates the script called when the form is submitted (here, index.php).</li>
            <li>
          The method field, which indicates the type of the HTTP query to be made. In this case, we're using GET, so we must include all of our query parameters in the requested URL. Here's how we format and encode them:
          <br />
              <b>/action?param_1=value_1&amp;param_2=value_2&amp;param_3=value_3</b>
            </li>
          </ul>
      So, whenever the form is submitted, the page <a href="/index.php?pass=input_value">/index.php?pass=input_value</a>
      gets called. (Note that you can request this link directly without submitting the form.)
        </p>

        <h2>The hash</h2>
        <p>
      Here's the PHP code that checks the password:
      <pre class="code_block"><?php echo htmlentities(file_get_contents('source_2.txt')); ?></pre>
          <br />
      You now have to find the hash using your brute force script. <i>For even the very wise cannot see all ends...</i>
        </p>

      </div>
    </body>
  </html>
