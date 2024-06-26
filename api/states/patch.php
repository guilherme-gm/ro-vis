<?php
require_once "../RecordApi/RecordPatchApiRoute.php";
require_once "StateParser.php";

class PatchStates extends RecordPatchApiRoute {
	public function __construct() {
		parent::__construct("states", new StateParser(), "Id");
	}
}

(new PatchStates())->run();
