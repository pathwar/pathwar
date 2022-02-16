<!DOCTYPE html>
<html lang="en" style="background: black">
<head>
  <title>1337Chan</title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
</head>
<center><p>Welcome to my new imageboard! It's not finished yet, but here's a preview to give you a sense of what it'll look like.</p></center>
<center><p>For now, only the admins can upload pictures, but soon you'll be able to, too!</p></center>

<?php
include('db.php');
//<a href="/admin.php">Admin<a/> //admin not finished yet
$query='SELECT * from  posts LIMIT 10;';
$res = $mysqli->query($query) or die($mysqli->error);
while ($row=mysqli_fetch_assoc($res)) {
	echo '<div class="post">';
	echo '<span class="info">Post '.$row['id'].' | Author : '.$row['author'].'</span>';
	$query='SELECT * FROM images where id='.$row['image_id'].';';
	$res2=$mysqli->query($query) or die($mysqli->error);
	echo '<img style="float: left;" src="get.php?file='.mysqli_fetch_assoc($res2)['path'].'">';
	$query='SELECT * FROM comments where post_id='.$row['id'].';';
	$res2=$mysqli->query($query) or die($mysqli->error);
	while ($row2=mysqli_fetch_assoc($res2)) {
		echo '<span class="comment" style="float: left; margin-top: 10px; margin-left: 10px;">';
		echo '<blockquote>'.htmlentities($row2['content']).'</blockquote>';
		echo '</span><br /><br /><br /><br /><br />';
	}
	echo '</div><br /><br /><br /><br /><br /><br /><br /><br /><br /><br />';
}
?>
</html>
