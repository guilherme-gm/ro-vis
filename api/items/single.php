<?php
require_once "../RecordApi/RecordSingleApiRoute.php";
require_once "ItemParser.php";

class ItemHistory extends RecordSingleApiRoute {
	public function __construct() {
		parent::__construct("items", new ItemParser(), "Id");
	}
}

(new ItemHistory())->run();
