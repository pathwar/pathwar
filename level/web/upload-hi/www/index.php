<?php

$target = "./";
$max_size = 100000;

$nom_file = $_FILES['file']['name'];
$size = $_FILES['file']['size'];
$tmp = $_FILES['file']['tmp_name'];

//----------------------
// UPLOAD SCRIPT
//----------------------

if($_POST['posted'])
{
  if($_FILES['file']['name'] &&
  strpos($_FILES['file']['name'], '.') == strrpos($_FILES['file']['name'], '.'))
  {
    if(move_uploaded_file($_FILES['file']['tmp_name'],$target.$_FILES
    ['file']['name']))
    {
      chmod($target.$_FILES['file']['name'], 0644);
      echo '<br>';
      echo '<p align="center">';
      echo '<b>Image uploaded !</b>';
      echo '<hr>';
      echo '<p align="center">';
      echo '<b>File :</b> '.$_FILES['file']['name'].'</br>';
      echo '<b>Size :</b> '.$_FILES['file']['size'].' Octets</br>';
      echo '<hr>';
      echo '<br>';
    }
    else
    {
      echo '<br>';
      echo '<p align="center">';
      echo '<b>Upload error ! </b><br><br><b>'.$_FILES['file']['error'].'</b>';
      echo '<br>';
    }
  }
  else
  {
    echo '<br>';
    echo '<p align="center">';
    echo '<b>The file is not in the jpg file format !</b>';
    echo '<br>';
  }
}
else
{
  echo '<br>';
  echo '<p align="center">';
  echo '<b>Form field is empty !</b>';
  echo '</font><br>';
}

?>
<html>
<head>
  <title>Image Upload</title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
</head>
<body>
  <div class="col-sm-4 col-sm-offset-4 text-center">
    <form enctype="multipart/form-data" action="<?php echo $PHP_SELF; ?>" method="POST">
      <br>
      <input type="hidden" name="posted" value="1">
      <input name="file" type="file" class="filestyle" data-iconName="glyphicon-inbox" data-buttonText="Choose image file">
      <br>
      <button type="submit" class="btn btn-default">Upload</button>
    </form>
  </div>
</body>
<script src="https://code.jquery.com/jquery-1.11.2.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
<script src="http://markusslima.github.io/bootstrap-filestyle/js/bootstrap-filestyle.min.js"></script>
</html>
