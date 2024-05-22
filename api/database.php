<?php
require_once "config.php";
require_once "Medoo.php";

use Medoo\Medoo;

$database = new Medoo([
	// [required]
	'type' => 'mysql',
	'host' => $config["host"],
	'database' => $config["db"],
	'username' => $config["username"],
	'password' => $config["password"],

	// [optional]
	'charset' => 'utf8mb4',
	'collation' => 'utf8mb4_general_ci',
	'port' => $config["port"],

	// [optional]
	// Error mode
	// Error handling strategies when the error has occurred.
	// PDO::ERRMODE_SILENT (default) | PDO::ERRMODE_WARNING | PDO::ERRMODE_EXCEPTION
	// Read more from https://www.php.net/manual/en/pdo.error-handling.php.
	'error' => PDO::ERRMODE_EXCEPTION,
]);
