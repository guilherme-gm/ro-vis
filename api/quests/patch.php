<?php
require_once "../RecordApi/RecordPatchApiRoute.php";
require_once "QuestParser.php";

class PatchQuests extends RecordPatchApiRoute {
	public function __construct() {
		parent::__construct("quests", new QuestParser(), "Id");
	}
}

(new PatchQuests())->run();
