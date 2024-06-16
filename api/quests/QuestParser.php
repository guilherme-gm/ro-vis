<?php
require_once "../RecordApi/IRecordParser.php";

class QuestParser implements IRecordParser {
	public function getFields() {
		return [
			"Patch",
			"_FileVersion",
			"Id",
			"Title",
			"Description",
			"Summary",
			"OldImage",
			"IconName",
			"NpcSpr",
			"NpcNavi",
			"NpcPosX",
			"NpcPosY",
			"RewardEXP",
			"RewardJEXP",
			"RewardItemList",
			"CoolTimeQuest",
		];
	}

	public function db2record($dbVal, $prefix) {
		$quest = [];

		foreach ($this->getFields() as $fld) {
			$quest[$fld] = $dbVal[$prefix . "_" . $fld];
		}
		$quest["RewardItemList"] = json_decode($quest["RewardItemList"], true);

		return $quest;
	}
}
