<?php
require_once "../RecordApi/RecordListApiRoute.php";

class ListQuests extends RecordListApiRoute {
	public function __construct() {
		parent::__construct("quests", ["Id", "Patch", "Title"], "Id");
	}
}

(new ListQuests())->run();
