<?php
require_once "../RecordApi/RecordPatchApiRoute.php";
require_once "ItemParser.php";

class PatchItems extends RecordPatchApiRoute {
	public function __construct() {
		parent::__construct("items", new ItemParser(), "Id");
	}
}

(new PatchItems())->run();
