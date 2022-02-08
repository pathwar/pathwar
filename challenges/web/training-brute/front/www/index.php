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
      The goal of this training is to learn how to script a simple bruteforce.
    </p>

        <h2>You shall not pass</h2>
        <p>
          <form name="bruteforce" action="index.php" method="get">
            Password <input type="text" name="pass" />
            <input type="submit" value="go" />
          </form>
        </p>
        <?php if (isset($_GET['pass'])) : ?>
        <p>
          <center>
            <?php if ($_GET['pass'] == md5("763")) : ?>
            Congratz! The passphrase is <b>__PASSPHRASE0__</b>
            <?php else : ?>
            <img width='600px' src="./gandalf.jpg" alt="u_shall_not_pass" />
            <br />
        Invalid password.
            <?php endif ?>
          </center>
        </p>
        <?php endif ?>


        <h2>Tools</h2>
        <p>
      For this training, you'll have to perform several HTTP requests, curl can be handy, but you are free
      to use whatever you want. Here's how to perform a loop in bash:

      <pre class="code_block"><?php echo htmlentities(file_get_contents('source_0.txt')); ?></pre>
          <br />
          Warning: <i>echo $i</i> is not the same as <i>echo -n $i</i>.

      You'll also need to generate a MD5 hash from a number (man md5sum, <a href="http://www.php.net/manual/en/function.md5.php">php.net</a>).
        </p>

        <h2>The Form</h2>
        <p>
      <pre class="code_block"><?php echo htmlentities(file_get_contents('source_1.txt')); ?></pre>
          <br />
      The interesting parts are:
          <ul>
            <li>The action field, which indicates the resource that is to be called whenever the form is submitted, here, index.php</li>
            <li>
          The method field, which indicates the type of the HTTP query is to be made, here, GET. What it means in practice is that all
          parameters are to be encoded in the called URL in this fashion:
          <br />
              <b>/action?param_1=value_1&amp;param_2=value_2&amp;param_3=value_3</b>
            </li>
          </ul>
      All this to say that whenever the form is submitted, the page <a href="/index.php?pass=input_value">/index.php?pass=input_value</a>
      gets called. You can directly use this link without submitting the form.
        </p>

        <h2>The hash</h2>
        <p>
      Here the PHP code that checks the password:
      <pre class="code_block"><?php echo htmlentities(file_get_contents('source_2.txt')); ?></pre>
          <br />
      You now have to find the hash using a bruteforce.
        </p>

      </div>
    </body>
  </html>
