<?php
require_once "../RecordApi/RecordListApiRoute.php";

class ListItems extends RecordListApiRoute {
	public function __construct() {
		parent::__construct("items", ["Id", "Patch", "IdentifiedName"], "Id");
	}
}

(new ListItems())->run();
