<?php
include('db.php');
if (isset($_POST['password'])) {
    $password = $_POST['password'];
    $query="SELECT * FROM users WHERE username='admin' AND password='$password'";
    $result = $mysqli->query($query) or die($mysqli->error);
    if ($result->num_rows > 0) { // login successful
    $data = $result->fetch_row();
    echo "Welcome $data[1] !<br />";
    echo "In case you have forgotten it, your password is : $data[2] (this is not the challenge passphrase, just the current mysql user password)";
    }
    else {
    echo '<div>Either the login or the password is wrong !</div>';
    }
}
else {
?>
<div class="container-fluid">
  <div class="row">
    <div class="span2">
      <div class="col-sm-3 col-md-2 sidebar">
    <ul class="nav nav-sidebar">
      <li><a href="/">Home</a></li>
      <li><a href="/?step=1">First injection</a></li>
      <li><a href="/?step=2">Second injection</a></li>
    </ul>
    <br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br />
      </div>
    </div>
    <div class="span10">
      You need to login as the admin user, but you don't know its password, but luckily (or not :)), the login form is vulnerable to a SQL injection ! <br />
      Try to exploit it to connect as the admin ! (Tips: have a look at the OR keyword)<br />
      The request looks like this :<br />
      <div class="code">
    <pre><?php echo htmlentities(file_get_contents('step1.txt'))?></pre>
      </div>
      <form method="POST">
    <input type="text" name="login" value="admin" disabled/><br />
    <input type="password" name="password" /> <br />
    <input type='submit'/> <br />
      </form>
    </div>
<?php
}
?>
