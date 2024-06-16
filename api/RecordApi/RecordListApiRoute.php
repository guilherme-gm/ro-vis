<?php
require_once "../ApiRoute.php";

class RecordListApiRoute extends ApiRoute {
	protected $table = "";

	protected $fields = [];

	protected $idField = "";

	function __construct($table, $fields, $idField) {
		$this->table = $table;
		$this->fields = $fields;
		$this->idField = $idField;
	}

	public function main() {
		$this->initDb();

		$start = $_GET['start'] ?? 0;

		$count = $this->db->count($this->table);

		$data = $this->db->select(
			$this->table,
			$this->fields,
			[
				"ORDER" => $this->idField,
				"LIMIT" => [$start, 100],
			]
		);

		$list = [];
		foreach ($data as $value) {
			$obj = [];
			foreach ($this->fields as $fld) {
				$obj[$fld] = $value[$fld];
			}
			array_push($list, $obj);
		}

		return [
			"total" => $count,
			"list" => $list,
		];
	}
}
