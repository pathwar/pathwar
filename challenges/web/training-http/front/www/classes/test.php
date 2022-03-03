<?php
$str = "lmkjljlj _('Hello') ";
include "LocaleManager.php";


echo LocaleManager::translate_tagged_parts ($str);