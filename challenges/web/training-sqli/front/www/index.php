<?php
//load classes
require_once ("classes/TplBlock.php");
$tpl = new TplBlock();

require_once ("classes/LocaleManager.php");
$currentLocale = new LocaleManager();
$tpl->addVars(array("langForm" => $currentLocale->get_locale_form()));

$step = isset($_GET['step'])? $_GET['step'] : 0;

//fill the left menus links
$navMenus = array(
	array("href"	=> "/", "caption"	=> _('Home')),
	array("href"	=> "/?step=1", "caption"	=> _('First injection')),
	array("href"	=> "/?step=2", "caption"	=> _('Second injection'))
);
foreach($navMenus as $menu){
	$tplNav = new TplBlock("nav");
	$tplNav->addVars($menu);
	$tpl->addSubBlock($tplNav);
}



switch($step){

	case "1":

		$tpl->addVars( array("pageTitle" => htmlentities( _('First injection')) ) );

		//Analyse post actions (win or not)
		$win = false;
		$loginFormPostMessage = "";
		
		if (isset($_POST['password'])) {

			require_once('db.php');

			$password = $_POST['password'];
			$query="SELECT * FROM users WHERE username='admin' AND password='$password'";
			$result = $mysqli->query($query) or die($mysqli->error);
			if ($result->num_rows > 0) { // login successful
			  $data = $result->fetch_row();
			  $win = true;
			  $adminUsername = $data[1];
			  $adminPassword = $data[2];
			}
			else {
				$loginFormPostMessage =  _('Either the login or the password is wrong!');
			}
		}

		//display the page
		if( $win )
		{
			$tplArticle = new TplBlock("article");
			$tplArticle->addVars(
				array("title"		=> _('Welcome') . " " . $adminUsername . "!"
					,"content"		=> _('In case you have forgotten it, your password is: ')
										. $adminPassword 
										._('(this is not the challenge passphrase, just the current mysql user password)')
					)			
			);
			$tpl->addSubBlock($tplArticle);

		}else{
			$tplForm = new TplBlock("loginForm");
			$tplForm->addVars(array("user" => "admin", "title"	=> _('Sign in'), "postMessage" => $loginFormPostMessage ));
			$tpl->addSubBlock($tplForm);

			$tplArticle = new TplBlock("article");
			$tplArticle->addVars(
				array("title"		=> _('You need to login as the admin user.')
					,"content"	=> _('But you don\'t know their password. (Luckily or not :)), the login form is vulnerable to a SQL injection!')
									. '<br/>' . _('Try to exploit it to connect as the admin! (<b>Hint</b>: Take a look at the OR keyword.)')
									. '<br/>' . _('The request you must exploit looks like this:')
					)			
			);
			$tplCode = new TplBlock("code");
			$tplCode->addVars( array( "content"	=> highlight_string( file_get_contents('step1.txt'),true ) ) );
			$tplArticle->addSubBlock($tplCode);
			$tpl->addSubBlock($tplArticle);
		}


		break;
	case "2":
		$tpl->addVars( array("pageTitle" => htmlentities( _('Second injection')) ) );


		//Analyse post actions (win or not)
		$win = false;
		$loginFormPostMessage = "";
		
		if (isset($_POST['password'])) {

			require_once('db.php');

			$password = $_POST['password'];
			$query="SELECT * FROM users2 WHERE username='user' AND password='$password'";
			$result = $mysqli->query($query) or die($mysqli->error);
			if ($result->num_rows > 0) { // login successful
				$data = $result->fetch_row();
				if($data[1] == 'admin'){
			  		
			  		$win = true;
			  		$adminUsername = $data[1];
				}else{
					$loginFormPostMessage =  _('Hey! What are you doing here? You are not the admin!');
				}
			} else {
				$loginFormPostMessage =  _('Either the login or the password is wrong!');
			}
		}
		//code to display the step2 page
		if( $win )
		{

			$tplArticle = new TplBlock("article");
			$tplArticle->addVars(
				array("title"		=> _('Welcome') . " " . $adminUsername . "!"
					,"content"		=> _('In case you have forgotten it, your password is: ')
										. '__PASSPHRASE0__'
					)			
			);
			$tpl->addSubBlock($tplArticle);

		}else{

			$tplForm = new TplBlock("loginForm");
			$tplForm->addVars(array("user" => "user", "title"	=> _('Sign in'), "postMessage" => $loginFormPostMessage ));
			$tpl->addSubBlock($tplForm);

			$tplArticle = new TplBlock("article");
			$tplArticle->addVars(
				array("title"		=> _('You need to login as the admin user again.')
					,"content"		=> _('This time, the admin has secured the website by removing their account from the database. Prove them wrong and login with their account !')
										. '<br/><b>' . _('TIP') . '</b> ' . _('take a look at the mysql documentation for UNION.')
										. _('The request looks like this:')
					)			
			);
			$tplCode = new TplBlock("code");
			$tplCode->addVars( array( "content"	=> highlight_string( file_get_contents('step2.txt'),true ) ) );
			$tplArticle->addSubBlock($tplCode);
			$tpl->addSubBlock($tplArticle);
		}

		break;
	default:

		//The default page. on a real site it should be a 404 page, here, it's the main page

		$tpl->addVars( array("pageTitle" => htmlentities( _('SQLi training')) ) );

		$articles = array(
						  array("title"		=> _('Welcome to the SQL Injection training grounds!')
							   ,"content"	=> _('In this training level, you\'ll get the hang of writing basic SQL injections.')
							   )
						 ,array("title"		=> _('What is a SQL injection?')
						 	   ,"content"	=> _('A SQL injection is a vulnerability commonly found on websites allowing an attacker to inject arbitrary content into a SQL request due to a lack of sanitization of user input.')
							   )
						 ,array("title"		=> _('What are the consequences?')
						 	   ,"content"	=> _('When you find a website vulnerable to a SQL injection, you can do pretty much whatever you want with its database. Depending on where the injection is, you can list passwords for all users, bypass a login form, or sometimes even execute arbitrary code on the server!')
						       )
						 ,array("title"		=> _('Some resources on the subject:')
							   ,"content"	=>'	
									<ul>
									<li><a href="https://en.wikipedia.org/wiki/SQL_injection">Wikipedia</a></li>
									<li><a href="https://www.owasp.org/index.php/SQL_Injection">OWASP</a></li>
									<li><a href="https://dev.mysql.com/doc/">' . _('MySQL Documentation') . '</a></li>
									</ul>'
								)
						);

		foreach( $articles as $article ){
			$tplArticle = new TplBlock("article");
			$tplArticle->addVars($article);
			$tpl->addSubBlock($tplArticle);
		}
		
		break;
}
echo $tpl->applyTplFile("templates/main.html");