CREATE TABLE `user` (
    `id` char(14) COLLATE utf8_bin NOT NULL,
    `username` varchar(255) COLLATE utf8_bin NOT NULL,
    `email` varchar(255) COLLATE utf8_bin DEFAULT NULL,
    `is_active` int(1) COLLATE utf8_bin DEFAULT 1,
    `secret_key` varchar(255) COLLATE utf8_bin NOT NULL,
    `created_at` int(11) NOT NULL,
    `secret_key_updated_at` int(11) NOT NULL,
    `deleted_at` int(11) DEFAULT NULL,
    PRIMARY KEY (`id`)
 );


CREATE TABLE `download_request` (
    `id` varchar(255) NOT NULL,
    `requested_user_id` varchar(255) NOT NULL,
    `download_type` varchar(255) NOT NULL,
    `download_status` varchar(255) NOT NULL,
    `zip_file_name` varchar(255) DEFAULT NULL,
    `failure_reason` varchar(255) DEFAULT NULL,
    `requested_at` int(11) NOT NULL,
    `finished_at` int(11) DEFAULT NULL,
    `failed_at` int(11) DEFAULT NULL,
    PRIMARY KEY (`id`)
 );

CREATE TABLE `file_details` (
    `id` varchar(255) NOT NULL,
    `download_request_id` varchar(255) NOT NULL,
    `file_name` varchar(255) NOT NULL,
    `file_size` int(11) DEFAULT NULL,
    `file_type` varchar(255) NOT NULL,
    `status` varchar(255) NOT NULL,
    `created_at` int(11) NOT NULL,
    PRIMARY KEY (`id`)
);