<?php header('X-Pathwar-Header: __PASSPHRASE3__'); ?>
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
			<h1>HTTP Headers</h1>
			<h2>Request Headers</h2>
			<p>When a client send a request to a server, it can send headers to give more information about the request (languages supported by the client, hostname of the requestd website, cookies,...).</p>
			<p>You can easily change the request headers sent by your browser with the help of an extension or a proxy that intercepts the requests made.</p>
			<?php
			if ($_SERVER['HTTP_USER_AGENT'] == 'Sup3rH4x0RBr0ws3r') {
				echo "W00t! Seems you are a real hacker ! Here's your passphrase: __PASSPHRASE2__";
			}
			else {
				echo 'Try to change your User-Agent to <b>Sup3rH4x0RBr0ws3r</b> to get the passphrase !';
			}
			?>
			<h2>Response Headers</h2>
			<p>When a web server sends a response to a client, it can send along a few response headers, containing informations about the response (like its size),
				or informations used by the browser to know for example if the page must be kept in cache.</p>
			<p>This page has sent you a custom header containing a passphrase, try to get it (you can use the debug console of your browser)</p>
		</div>
	</div>
</div>
</html>
