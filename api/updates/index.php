<?php
require_once "../ApiRoute.php";

class ListUpdates extends ApiRoute {
	public function main() {
		$this->initDb();

		$start = $_GET['start'] ?? 0;

		$count = $this->db->count("updates");

		$data = $this->db->select(
			"updates",
			["id", "order", "updates", "tags"],
			[
				"ORDER" => "order",
				"LIMIT" => [$start, 100],
			]
		);

		$list = [];
		foreach ($data as $value) {
			array_push($list, [
				"id" => $value["id"],
				"order" => $value["order"],
				"updates" => json_decode($value["updates"], true),
				"tags" => explode(",", $value["tags"]),
			]);
		}

		return [
			"total" => $count,
			"list" => $list,
		];
	}
}

(new ListUpdates())->run();
