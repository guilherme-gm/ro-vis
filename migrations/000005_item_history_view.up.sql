-- sqlc doesn't handle left joins very well, so we make it into a view so
-- it is forced to create a new type. Not ideal, but should not cause huge performance issues

CREATE VIEW previous_item_history_vw AS (
	SELECT prev.* FROM item_history curr LEFT JOIN item_history prev ON prev.history_id = curr.previous_history_id
);
