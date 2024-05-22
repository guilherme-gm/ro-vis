<?php
require_once "database.php";

$handle = fopen("./files/ro-vis.updates-nl.json", "r");
if ($handle) {
	$arr = [];
	while (($line = fgets($handle)) !== false) {
		$obj = json_decode($line, true);
		array_push($arr, [
			"id" => $obj["_id"],
			"order" => $obj["order"],
			"updates" => json_encode($obj["updates"]),
			"patches" => json_encode($obj["patches"]),
			"tags" => "",
		]);

		if (count($arr) == 50) {
			$database->insert("updates", $arr);
			$arr = [];
		}
	}

	if (count($arr) > 0) {
		$database->insert("updates", $arr);
	}

	fclose($handle);
} else {
	echo "File not found";
}

echo "Done";
