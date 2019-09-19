<?php
if (!isset($_COOKIE['chocolate'])) {
	setcookie('chocolate', 'isbad');
}
?>
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
			<h1>HTTP Cookies</h1>
			<p>HTTP itself is a stateless protocol (meaning that each request is independant), which is problematic for most web applications as they need to keep track of users sessions (ie, is the user logged in, what products did he add to his cart,...).
		To bypass this limitations, HTTP cookies were created.</p>
			<h2>How HTTP cookies works</h2>	
			<p>A cookie is typically set by the server by sending a <b>Set-Cookie</b> header to the client per cookie.
			A response that set cookies looks like this :</p>
			<div class="code" >
				<pre>
GET / HTTP/1.1
Host: foo.com

HTTP/1.1 200 OK
Set-Cookie: bar=42; Path=/; Domain=.foo.com
Set-Cookie: session=12345; HTTPOnly
...
			</pre>
			</div>
			<p>Cookies are automatically sent by the client with the header 'Cookie' containing all the valid cookies for the current domain, separated by a semi-colon (you will have the occasion to study this behaviour latter in some levels :)). From the previous example, any subsequent request to foo.com will include the two cookies set by the server :</p>
		<div class="code">
			<pre>
GET / HTTP/1.1
Host: foo.com
Cookie: bar=42; session=12345

HTTP/1.1 200 OK
...
			</pre>
		</div>
		
		<h2>Cookies Manipulation</h2>
		<p>When you opened this page for the first time, the server sent a cookie named <b>chocolate</b> to your browser.
			As cookie as stored client-side, you are free to modify them as you will (delete them, change their values, attributes,...).
			There is a lots of way to do this, but the simplest is to use a browser extensions (for Chrome or Chromium, you can use EditThisCookie, for other browsers, check <a href="http://lmgtfy.com/?q=firefox+edit+cookies">here</a>)</p>
		<?php
			if (isset($_COOKIE['chocolate']) && $_COOKIE['chocolate'] == 'isgood') {
				echo 'Congrats! The passphrase is __PASSPHRASE5__';
			}
		else {
			echo 'Try to set the value of the cookie <b>chocolate</b> to <b>isgood</b> !';
		}
		?>
		</div>
	</div>
</div>
</html>
