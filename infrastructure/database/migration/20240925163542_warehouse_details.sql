-- +goose Up
-- +goose StatementBegin
CREATE TABLE `warehouse_details` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `warehouse_id` int(11) NOT NULL,
  `qty` int(11) DEFAULT NULL,
  `item_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `warehouse_details_warehouse_id_IDX` (`warehouse_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `warehouse_details`;
-- +goose StatementEnd
