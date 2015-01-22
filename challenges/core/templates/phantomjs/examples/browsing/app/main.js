var page = require('webpage').create();

page.onConsoleMessage = function(msg) {
    console.log(msg);
};

page.open("http://radio.m1ch3l.biz/", function(status) {
    if (status === "success") {

	// inject local file
	console.log(
	    page.injectJs("inject.js") ?
		"Injected inject.js!" :
		"Failed to inject inject.js, file not found ?"
	);

	// include remote file
        page.includeJs("http://ajax.googleapis.com/ajax/libs/jquery/1.6.1/jquery.min.js", function() {
            page.evaluate(function() {
		var identifier = 'audio source' 
                console.log('jQuery("'+identifier+'") -> ' + jQuery(identifier));
            });
            // phantom.exit();
        });

	setInterval(function() {
	    page.evaluate(function() {
		console.log('Ping from page');
	    });
	}, 2000);
    }
});
