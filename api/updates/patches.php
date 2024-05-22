<?php
require_once "../ApiRoute.php";

class ListPatches extends ApiRoute {
	public function main() {
		$id = $_GET['id'];
		if (!isset($id)) {
			echo json_encode([
				'Error' => 'Missing ID',
			]);
			die();
		}

		$this->initDb();

		$data = $this->db->get(
			"updates",
			["id", "patches"],
			["id" => $id]
		);

		return [
			"id" => $data["id"],
			"patches" => json_decode($data["patches"], true),
		];
	}
}

(new ListPatches())->run();
