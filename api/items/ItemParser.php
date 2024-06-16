<?php
require_once "../RecordApi/IRecordParser.php";

class ItemParser implements IRecordParser {
	public function getFields() {
		return [
			"Patch",
			"_FileVersion",
			"Id",
			"IdentifiedName",
			"IdentifiedDescription",
			"IdentifiedSprite",
			"UnidentifiedName",
			"UnidentifiedDescription",
			"UnidentifiedSprite",
			"SlotCount",
			"IsBook",
			"CanUseBuyingStore",
			"CardPrefix",
			"CardPostfix",
			"CardIllustration",
			"ClassNum",
			"MoveInfo_canDrop",
			"MoveInfo_canTrade",
			"MoveInfo_canMoveToStorage",
			"MoveInfo_canMoveToCart",
			"MoveInfo_canSellToNpc",
			"MoveInfo_canMail",
			"MoveInfo_canAuction",
			"MoveInfo_canMoveToGuildStorage",
			"MoveInfo_commentName",
		];
	}

	public function db2record($dbVal, $prefix) {
		$item = [];

		foreach ($this->getFields() as $fld) {
			$item[$fld] = $dbVal[$prefix . "_" . $fld];
		}

		$item["MoveInfo"] = [
			"canDrop" => $item["MoveInfo_canDrop"],
			"canTrade" => $item["MoveInfo_canTrade"],
			"canMoveToStorage" => $item["MoveInfo_canMoveToStorage"],
			"canMoveToCart" => $item["MoveInfo_canMoveToCart"],
			"canSellToNpc" => $item["MoveInfo_canSellToNpc"],
			"canMail" => $item["MoveInfo_canMail"],
			"canAuction" => $item["MoveInfo_canAuction"],
			"canMoveToGuildStorage" => $item["MoveInfo_canMoveToGuildStorage"],
		];

		unset($item["MoveInfo_canDrop"]);
		unset($item["MoveInfo_canTrade"]);
		unset($item["MoveInfo_canMoveToStorage"]);
		unset($item["MoveInfo_canMoveToCart"]);
		unset($item["MoveInfo_canSellToNpc"]);
		unset($item["MoveInfo_canMail"]);
		unset($item["MoveInfo_canAuction"]);
		unset($item["MoveInfo_canMoveToGuildStorage"]);

		return $item;
	}
}
