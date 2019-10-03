<html>
<head>
  <title>Training Include</title>
  <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css">
</head>
<body>

    <div class="container">
      <div class="page-header">
        <h1>Training Include</h1>
      </div>
   <?php

   if (!isset($_GET['page'])) {
	 include('home.php');
   } else {
	 include($_GET['page'].'.php');
   }
?>
    </div>

</body>
<script src="//code.jquery.com/jquery-1.11.2.min.js"></script>
<script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/js/bootstrap.min.js"></script>
<script src="//markusslima.github.io/bootstrap-filestyle/js/bootstrap-filestyle.min.js"></script>
</html>
