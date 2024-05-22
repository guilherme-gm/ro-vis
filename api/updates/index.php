<?php
require_once "../ApiRoute.php";

class ListUpdates extends ApiRoute {
	public function main() {
		$this->initDb();

		$start = $_GET['start'] ?? 0;

		$data = $this->db->select(
			"updates",
			["id", "order", "updates", "tags"],
			[
				"ORDER" => "order",
				"LIMIT" => [$start, 100],
			]
		);

		$output = [];
		foreach ($data as $value) {
			array_push($output, [
				"id" => $value["id"],
				"order" => $value["order"],
				"updates" => json_decode($value["updates"], true),
				"tags" => explode(",", $value["tags"]),
			]);
		}

		return $output;
	}
}

(new ListUpdates())->run();
