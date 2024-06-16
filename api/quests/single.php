<?php
require_once "../RecordApi/RecordSingleApiRoute.php";
require_once "QuestParser.php";

class QuestHistory extends RecordSingleApiRoute {
	public function __construct() {
		parent::__construct("quests", new QuestParser(), "Id");
	}
}

(new QuestHistory())->run();
