<?php

require_once ("classes/LocaleManager.php");
$currentLocale = new LocaleManager();

$path = isset($_GET["path"]) ? $_GET["path"] : "/" ;


switch ($path){
    case "you-lose":
        $content = LocaleManager::translate_tagged_parts( file_get_contents("templates/you-lose.html") );
        break;
    case "you-win":
        $content = LocaleManager::translate_tagged_parts( file_get_contents("templates/you-win.html") );
        break;
    case "/":
        $content = LocaleManager::translate_tagged_parts( file_get_contents("templates/index.html") );
        break;

    default:
        header("HTTP/1.0 404 Not Found");
        $content = "<html><head><title>" . _('404 page not found') . "</title></head><body><h1>" . _('404 page not found') . "</h1></body></html>";
        break;


}

echo str_replace( "{{langForm}}", $currentLocale->get_locale_form(), $content );