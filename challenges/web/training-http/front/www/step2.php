<?php header('X-Pathwar-Header: ValidationCodeWillBeGivenAtANextStep'); ?>
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
                    <li><a href="/?step=4">HTTP Cookies</a></li>
                </ul>
                <br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br /><br />
            </div>
        </div>
        <div class="content">
            <h1>HTTP Headers</h1>
            <h2>Request Headers</h2>
            <p>When a client send a request to a server, it can send headers to give more information about the request (e.g. languages supported by the client, hostname of the requested website, stored cookies).</p>
            <p>You can easily change the request headers sent by your browser with the help of an extension or a proxy that intercepts the requests it makes.</p>
            <?php
            if ($_SERVER['HTTP_USER_AGENT'] == 'Sup3rH4x0RBr0ws3r') {
                echo "W00t! Seems you are a real hacker!";
            }
            else {
                echo 'Try to change your User-Agent to <b>Sup3rH4x0RBr0ws3r</b> to get the passphrase !';
            }
            ?>
            <h2>Response Headers</h2>
            <p>When a web server sends a response to a client, it can send along a few response headers, containing information about the response (like its size),
                or instructions for the browser, e.g. if, or for how long, it should cache the page.</p>
            <p>This page has sent you a custom header containing a passphrase. See if you can find it! (Tip: You can use the debug console of your browser.)</p>
        </div>
    </div>
</div>
</html>
