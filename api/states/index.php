<?php
require_once "../RecordApi/RecordListApiRoute.php";

class ListStates extends RecordListApiRoute {
	public function __construct() {
		parent::__construct("states", ["Id", "Patch", "Constant"], "Id");
	}
}

(new ListStates())->run();
