SELECT `dc`.`id` AS `id`,
       `dc`.`name` AS `category_name`,
       `dc`.`status` AS `category_status`,
       `dci`.`sort` AS `sort`,
       `acc`.`id` AS `account_id`,
       `acc`.`number`,
       `acc`.`currency_code`,
       `acc`.`name`,
       `acc`.`bank_code`,
       `acc`.`bank_branch`,
       `acc`.`min_amount`,
       `acc`.`max_amount`,
       `acc`.`type`,
       `acc`.`qrcode`,
       `acc`.`status` AS `account_status`
FROM   `bank_deposit_category` AS `dc`
       JOIN `bank_deposit_category_item` AS `dci`
         ON `dc`.`id` = `dci`.`category_id`
       JOIN `bank_account` AS `acc`
         ON `dci`.account_id = `acc`.`id`

