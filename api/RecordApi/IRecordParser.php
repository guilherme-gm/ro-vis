<?php

interface IRecordParser {
	public function getFields();

	public function db2record($val, $prefix);
}
