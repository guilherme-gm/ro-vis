<?php
require_once "../RecordApi/RecordSingleApiRoute.php";
require_once "StateParser.php";

class StateHistory extends RecordSingleApiRoute {
	public function __construct() {
		parent::__construct("states", new StateParser(), "Id");
	}
}

(new StateHistory())->run();
