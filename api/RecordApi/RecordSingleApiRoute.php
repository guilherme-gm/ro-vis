<?php
require_once "../ApiRoute.php";
require_once "IRecordParser.php";

class RecordSingleApiRoute extends ApiRoute {
	protected $table = "";

	protected $fields = [];

	protected $idField = "";

	protected IRecordParser $parser;

	function __construct($table, IRecordParser $parser, $idField) {
		$this->table = $table;
		$this->fields = [...$parser->getFields()];
		$this->idField = $idField;
		$this->parser = $parser;
	}

	private function buildFieldList() {
		$fieldList = [
			"previous" => [],
			"current" => [],
		];

		foreach ($this->fields as $fld) {
			array_push($fieldList["previous"], "previous.$fld (prev_$fld)");
		}

		foreach ($this->fields as $fld) {
			array_push($fieldList["current"], $this->table . "_history.$fld (cur_$fld)");
		}

		return $fieldList;
	}

	private function remap($val, $prefix) {
		if ($val[$prefix . "_" . $this->idField] === null) {
			return null;
		}

		return $this->parser->db2record($val, $prefix);
	}

	public function main() {
		$this->initDb();

		if (!isset($_GET['id'])) {
			return ['error' => 'No id'];
		}

		$start = $_GET['start'] ?? 0;
		$id = $_GET['id'];

		$count = $this->db->count($this->table . "_history", [$this->idField => $id]);

		$list = $this->db->select(
			$this->table . "_history",
			["[>]" . $this->table . "_history (previous)" => ["PreviousId" => "HistoryId"]],
			$this->buildFieldList(),
			[
				$this->table . "_history." . $this->idField => $id,
				"ORDER" => "cur_" . $this->idField,
				"LIMIT" => [$start, 100],
			]
		);

		$list = array_map(function ($val) {
			return [
				"previous" => $this->remap($val["previous"], "prev"),
				"current" => $this->remap($val["current"], "cur"),
			];
		}, $list);

		$current = $this->remap($this->db->first($this->table, [$this->idField => $id]));

		if (count($list) > 0) {
			$lastHistory = $list[count($list) - 1];
			array_push($list, [
				"previous" => $lastHistory["current"],
				"current" => $current,
			]);
		}

		$count++;

		if (count($list) > 0) {
			return [
				"total" => $count,
				"list" => $list,
			];
		}
	}
}
