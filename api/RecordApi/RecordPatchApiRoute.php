<?php
require_once "../ApiRoute.php";
require_once "IRecordParser.php";

class RecordPatchApiRoute extends ApiRoute {
	protected $table = "";

	protected $fields = [];

	protected $idField = "";

	protected IRecordParser $parser;

	function __construct($table, IRecordParser $parser, $idField) {
		$this->table = $table;
		$this->fields = ['HistoryId', 'Patch', ...$parser->getFields()];
		$this->idField = $idField;
		$this->parser = $parser;
	}

	private function buildFieldList() {
		$fieldList = [
			"previous" => [],
			"current" => [],
			"latest.patch (lastUpdate)",
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

		if (!$_GET['patch']) {
			return ['error' => 'No patch'];
		}

		$start = $_GET['start'] ?? 0;
		$patch = $_GET['patch'];

		$count = $this->db->count($this->table . "_history", ["Patch" => $patch]);

		$list = $this->db->select(
			$this->table . "_history",
			[
				"[>]" . $this->table . "_history (previous)" => ["PreviousId" => "HistoryId"],
				"[>]" . $this->table . " (latest)" => [$this->idField => $this->idField],
			],
			$this->buildFieldList(),
			[
				$this->table . "_history.Patch" => $patch,
				"ORDER" => "cur_" . $this->idField,
				"LIMIT" => [$start, 100],
			]
		);

		$list = array_map(function ($val) {
			return [
				"previous" => $this->remap($val["previous"], "prev"),
				"current" => $this->remap($val["current"], "cur"),
				"lastUpdate" => $val["lastUpdate"],
			];
		}, $list);

		return [
			"total" => $count,
			"list" => $list,
		];
	}
}
