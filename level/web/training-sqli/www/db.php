<?php
mysql_connect(getenv('MYSQL_PORT_3306_TCP_ADDR'), 'training_sqli') or die(mysql_error());
mysql_select_db('training_sqli');
?>