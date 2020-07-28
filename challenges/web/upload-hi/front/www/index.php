<?php

$target = './';

//----------------------
// UPLOAD SCRIPT
//----------------------

if($_POST['posted'])
{
  ob_start();
  if(move_uploaded_file($_FILES['file']['tmp_name'],$target.$_FILES
  ['file']['name']))
  {
    chmod($target.$_FILES['file']['name'], 0644);
    echo '<br>';
    echo '<p align="center">';
    echo '<b>Image uploaded !</b>';
    echo '<hr>';
    echo '<p align="center">';
    echo '<b>File :</b> <a href="'.$target.$_FILES['file']['name'].'">'.$_FILES['file']['name'].'</a><br>';
    echo '<b>Size :</b> '.$_FILES['file']['size'].' Octets<br><br>';
    echo '<img src='.$target.$_FILES['file']['name'].' alt="Preview" style="max-height: 420px; max-width: 420px"">';
    echo '<hr>';
    echo '<br>';
    echo "\r\n";
  }
  else
  {
    echo '<br>';
    echo '<p align="center">';
    echo '<b>Upload error !</b><br><br><b>Error : '.$_FILES['file']['error'].'</b>';
    echo '<br>';
    echo "\r\n";
  }
  $php_output = ob_get_contents();
  ob_end_clean();
}
?>

<html>
<head>
  <title>Image Upload</title>
  <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
</head>
<body>
  <div class="col-sm-4 col-sm-offset-4 text-center">
    <div>
      <?php echo($php_output);?>
    </div>
    <form enctype="multipart/form-data" action="<?php echo $PHP_SELF; ?>" method="POST">
      <br>
      <input type="hidden" name="posted" value="1">
      <input name="file" type="file" class="filestyle" data-iconName="glyphicon-inbox" data-buttonText="Choose image file">
      <br>
      <button type="submit" class="btn btn-default">Upload</button>
    </form>
  </div>
</body>
<script src="//code.jquery.com/jquery-1.11.2.min.js"></script>
<script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
<script src="//markusslima.github.io/bootstrap-filestyle/js/bootstrap-filestyle.min.js"></script>
</html>
