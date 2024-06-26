<?php
require_once "../RecordApi/IRecordParser.php";

class StateParser implements IRecordParser {
	public function getFields() {
		return [
			"Patch",
			"_FileVersion",
			"Id",
			"Constant",
			"Description",
			"HasTimeLimit",
			"TimeStrLineNum",
			"HasEffectImage",
			"IconImage",
			"IconPriority",
		];
	}

	public function db2record($dbVal, $prefix) {
		$state = [];

		foreach ($this->getFields() as $fld) {
			$state[$fld] = $dbVal[$prefix . "_" . $fld];
		}
		$state["Description"] = json_decode($state["Description"], true);

		return $state;
	}
}
