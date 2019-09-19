<?php include("header.php"); ?>
<div class="container-fluid">
	<div class="row">
		<div class="span2">
			<div class="col-sm-3 col-md-2 sidebar">				
				<ul class="nav nav-sidebar">
					<li><a href="/">Home</a></li>
					<li><a href="/?step=1">HTTP Methods</a></li>
					<li><a href="/?step=2">HTTP Headers</a></li>
					<li><a href="/?step=3">HTTP Redirection</a></li>
					<li><a href="/?step=4">HTTP Cookie</a></li>
				</ul>
				<br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br />
			</div>
		</div>
		<div class="content">
			<h2>HTTP Redirection</h2>
			<p>As seen previously, web servers can sends HTTP headers to the client.
				A commonly used header is <b>Location</b>. This header is used to redirect the client to another page or website (for example, it is commonly used on some websites to redirect the user to its profile page upon a succesful login). <br />
				The following link will send an HTTP redirection to this page to your browser that will be automatically followed.<br />
				<a href="/redir.php">Redirection</a><br />
				HTTP client are not required to follow redirections. For example, curl does not follow them by default.<br />
				Request the redirection page with curl to get the passphrase.
			</p>
		</div>
	</div>
</div>
</html>
