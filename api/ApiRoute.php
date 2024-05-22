<?php
require_once "Medoo.php";

use Medoo\Medoo;

abstract class ApiRoute {
	protected $db = null;

	protected function initDb() {
		require_once "config.php";

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

		$this->db = $database;
	}

	/**
	 *  An example CORS-compliant method.  It will allow any GET, POST, or OPTIONS requests from any
	 *  origin.
	 *
	 *  In a production environment, you probably want to be more restrictive, but this gives you
	 *  the general idea of what is involved.  For the nitty-gritty low-down, read:
	 *
	 *  - https://developer.mozilla.org/en/HTTP_access_control
	 *  - https://fetch.spec.whatwg.org/#http-cors-protocol
	 *
	 */
	private function cors() {
		// Allow from any origin
		if (isset($_SERVER['HTTP_ORIGIN'])) {
			// Decide if the origin in $_SERVER['HTTP_ORIGIN'] is one
			// you want to allow, and if so:
			header("Access-Control-Allow-Origin: {$_SERVER['HTTP_ORIGIN']}");
			header('Access-Control-Allow-Credentials: true');
			header('Access-Control-Max-Age: 86400');    // cache for 1 day
		}

		// Access-Control headers are received during OPTIONS requests
		if ($_SERVER['REQUEST_METHOD'] == 'OPTIONS') {

			if (isset($_SERVER['HTTP_ACCESS_CONTROL_REQUEST_METHOD']))
				// may also be using PUT, PATCH, HEAD etc
				header("Access-Control-Allow-Methods: GET, POST, OPTIONS");

			if (isset($_SERVER['HTTP_ACCESS_CONTROL_REQUEST_HEADERS']))
				header("Access-Control-Allow-Headers: {$_SERVER['HTTP_ACCESS_CONTROL_REQUEST_HEADERS']}");

			exit(0);
		}
	}

	public function run() {
		$this->cors();
		header('Content-Type: application/json');

		$res = $this->main();

		echo json_encode($res);
	}

	protected abstract function main();
}
