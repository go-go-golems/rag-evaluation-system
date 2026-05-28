---
Title: Dump Schema Inspection
Ticket: RAGEVAL-002
Status: active
Topics: [rag, ingestion, database, corpus]
DocType: reference
Intent: evidence
---

# Dump inspection
path: /home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2
size_bytes: 44889956

## Selected CREATE TABLE blocks

### search_products
CREATE TABLE `search_products` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int DEFAULT NULL,
  `parent_id` int DEFAULT NULL,
  `title` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `botanical_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `product_id` (`product_id`) USING BTREE,
  UNIQUE KEY `idx_product_id` (`product_id`),
  KEY `parent_id` (`parent_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type` (`type`),
  FULLTEXT KEY `title` (`title`),
  FULLTEXT KEY `title_2` (`title`),
  FULLTEXT KEY `title_botanical_name` (`title`,`botanical_name`)
) ENGINE=InnoDB AUTO_INCREMENT=15907977 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

### wp_postmeta
CREATE TABLE `wp_postmeta` (
  `meta_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `post_id` bigint unsigned NOT NULL DEFAULT '0',
  `meta_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `meta_value` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`meta_id`),
  KEY `post_id` (`post_id`),
  KEY `meta_key` (`meta_key`(191)),
  KEY `meta_value` (`meta_value`(100)),
  KEY `bk1` (`post_id`,`meta_key`),
  FULLTEXT KEY `wp_postmeta_fulltext` (`meta_value`)
) ENGINE=InnoDB AUTO_INCREMENT=41396773 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

### wp_posts
CREATE TABLE `wp_posts` (
  `ID` bigint unsigned NOT NULL AUTO_INCREMENT,
  `post_author` bigint unsigned NOT NULL DEFAULT '0',
  `post_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `post_date_gmt` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `post_content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `post_title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `post_excerpt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `post_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'publish',
  `comment_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'open',
  `ping_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'open',
  `post_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `post_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `to_ping` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `pinged` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `post_modified` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `post_modified_gmt` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `post_content_filtered` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `post_parent` bigint unsigned NOT NULL DEFAULT '0',
  `guid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `menu_order` int NOT NULL DEFAULT '0',
  `post_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'post',
  `post_mime_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `comment_count` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`ID`),
  KEY `type_status_date` (`post_type`,`post_status`,`post_date`,`ID`),
  KEY `post_parent` (`post_parent`),
  KEY `post_author` (`post_author`),
  KEY `post_name` (`post_name`(191)),
  KEY `guid` (`guid`(191)),
  KEY `type_status_author` (`post_type`,`post_status`,`post_author`)
) ENGINE=InnoDB AUTO_INCREMENT=878756 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

### wp_term_relationships
CREATE TABLE `wp_term_relationships` (
  `object_id` bigint unsigned NOT NULL DEFAULT '0',
  `term_taxonomy_id` bigint unsigned NOT NULL DEFAULT '0',
  `term_order` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`object_id`,`term_taxonomy_id`),
  KEY `term_taxonomy_id` (`term_taxonomy_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

### wp_term_taxonomy
CREATE TABLE `wp_term_taxonomy` (
  `term_taxonomy_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `term_id` bigint unsigned NOT NULL DEFAULT '0',
  `taxonomy` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `description` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent` bigint unsigned NOT NULL DEFAULT '0',
  `count` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`term_taxonomy_id`),
  UNIQUE KEY `term_id_taxonomy` (`term_id`,`taxonomy`),
  KEY `taxonomy` (`taxonomy`)
) ENGINE=InnoDB AUTO_INCREMENT=735 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

### wp_terms
CREATE TABLE `wp_terms` (
  `term_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `slug` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `term_group` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`term_id`),
  KEY `slug` (`slug`(191)),
  KEY `name` (`name`(191))
) ENGINE=InnoDB AUTO_INCREMENT=734 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

### wp_wc_product_meta_lookup
CREATE TABLE `wp_wc_product_meta_lookup` (
  `product_id` bigint NOT NULL,
  `sku` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `virtual` tinyint(1) DEFAULT '0',
  `downloadable` tinyint(1) DEFAULT '0',
  `min_price` decimal(19,4) DEFAULT NULL,
  `max_price` decimal(19,4) DEFAULT NULL,
  `onsale` tinyint(1) DEFAULT '0',
  `stock_quantity` double DEFAULT NULL,
  `stock_status` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'instock',
  `rating_count` bigint DEFAULT '0',
  `average_rating` decimal(3,2) DEFAULT '0.00',
  `total_sales` bigint DEFAULT '0',
  `tax_status` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'taxable',
  `tax_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `global_unique_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  PRIMARY KEY (`product_id`),
  KEY `virtual` (`virtual`),
  KEY `downloadable` (`downloadable`),
  KEY `stock_status` (`stock_status`),
  KEY `stock_quantity` (`stock_quantity`),
  KEY `onsale` (`onsale`),
  KEY `min_max_price` (`min_price`,`max_price`),
  KEY `sku` (`sku`(50))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

## INSERT statement counts
search_products	1
wp_actionscheduler_actions	2
wp_actionscheduler_claims	1
wp_actionscheduler_groups	1
wp_actionscheduler_logs	1
wp_commentmeta	1
wp_comments	16
wp_options	4
wp_postmeta	107
wp_posts	31
wp_redirects	1
wp_term_relationships	2
wp_term_taxonomy	2
wp_termmeta	1
wp_terms	1
wp_usermeta	19
wp_users	1
wp_wc_admin_note_actions	1
wp_wc_admin_notes	1
wp_wc_category_lookup	1
wp_wc_order_addresses	2
wp_wc_order_operational_data	1
wp_wc_orders	2
wp_wc_orders_meta	11
wp_wc_product_attributes_lookup	3
wp_wc_product_download_directories	1
wp_wc_product_meta_lookup	2
wp_wc_tax_rate_classes	1
wp_woocommerce_api_keys	1
wp_woocommerce_attribute_taxonomies	1
wp_woocommerce_order_itemmeta	11
wp_woocommerce_order_items	2
wp_woocommerce_shipping_table_rates	1
wp_woocommerce_shipping_zone_locations	1
wp_woocommerce_shipping_zone_methods	1
wp_woocommerce_shipping_zone_shipping_methods	1
wp_woocommerce_shipping_zones	1
wp_woocommerce_tax_rates	1

## wp_posts post_type/status tuple counts
attachment	inherit	17038
attachment	private	37
custom_css	publish	1
export_template	publish	2
faq	publish	35
hb_pricing_table	publish	2
j79_wc_price_rule	publish	1
jetpack_migration	draft	2
mc4wp-form	publish	1
my_keywords	publish	6
nav_menu_item	draft	1
nav_menu_item	publish	16
nm-forms	publish	1
oembed_cache	publish	1
opanda-item	publish	6
page	draft	13
page	private	1
page	publish	120
page	trash	8
post	draft	1
post	publish	483
post	trash	5
product	draft	111
product	private	6
product	publish	2594
product	trash	26
product_variation	draft	2
product_variation	private	289
product_variation	publish	11913
product_variation	trash	20
question_answer	publish	3
shop_coupon	draft	126
shop_coupon	private	1
shop_coupon	publish	1752
shop_coupon	trash	16
shop_order_placehold	draft	5000
shop_webhook	publish	1
ttc_banner	publish	19
ttc_guide	publish	19
wafs	publish	1
wp_global_styles	publish	1
wpcf7_contact_form	publish	3
wpforms	publish	5

## wp_posts samples by post_type

### attachment
758	inherit	credit-cards	credit-cards
3517	inherit	soil-wallpaper-1366x768-edited-2	soil-wallpaper-1366x768-edited
3522	inherit	thuja-bokeh-edited-2	thuja-bokeh-edited-2
3557	inherit	1000x1000-product_image_test	1000x1000-product_image_test
3597	inherit	potted-clementine-mandarin-tree	Potted Clementine Mandarin Tree

### custom_css
95951	publish	highendwp-child	HighendWP-child

### export_template
525288	publish	quick-customer	Quick Customer
525757	publish	order-export	Order export

### faq
4116	publish	hardiness-zone	What is my Hardiness Zone?
4120	publish	forms-payment-accepted	What forms of payment are accepted?
4121	publish	shipping-handling	How much is shipping and handling?
4122	publish	send-gift	Can I send a gift?
4123	publish	discounts-bulk-wholesale-orders	Are there discounts for bulk and wholesale orders?

### hb_pricing_table
642	publish	pricing-1	Pricing 1
1004	publish	pricing-3-items	Pricing - 3 Items

### j79_wc_price_rule
13739	publish	test	Test

### jetpack_migration
30336	draft		widget_image
30337	draft		sidebars_widgets

### mc4wp-form
10023	publish	default-sign-up-form	Default sign-up form

### my_keywords
10095	publish	starter-fertilizer	Starter Fertilizer 2
12491	publish	guarantee	
30654	publish	bio-tone	Starter Fertilizer
32540	publish	citrus-tone	Citrus-tone 4 Pound Bag
38285	publish	tree-staking-kit	Tree Staking Kit

### nav_menu_item
3240	publish	3240	
7706	draft		Gardenias
398491	publish	about-the-tree-center	About The Tree Center
398492	publish	faq	FAQ
398493	publish	398493	

### nm-forms
9715	publish	application-submission	Application Submission

### oembed_cache
69805	publish	6c04211216c9f533bdc976768a5f04a4	

### opanda-item
10443	publish	opanda_default_signin_locker	Sign-In Locker (default)
10444	publish	opanda_default_social_locker	Social Locker (default)
10451	publish	crape-myrtle-guide	Crape Myrtle Guide
10467	publish	nellie-stevens-holly	Nellie Stevens Holly
10486	publish	general-homepage-share	General Homepage Share

### page
119	publish	home-2	Home
2708	publish	sitemap	Sitemap
3121	publish	login	Login
3217	publish	cart	My Cart
3218	publish	checkout	Checkout

### post
4332	publish	the-best-privacy-trees	Best Privacy Trees For Your Backyard
4355	publish	shade-trees-fast-growing-trees-quickly-provide-shade	Best Shade Trees – Fast Growing Trees That Will Quickly Provide You Shade
4498	publish	japanese-maple-trees-everything-you-wanted-to-know	Japanese Maple Trees - Everything You Wanted To Know
4884	publish	buy-weeping-willow-trees	Everything To Know Before You Buy Weeping Willow Trees
5048	publish	dogwood-tree-facts	Dogwood Tree Facts

### product
3699	publish	thuja-green-giant	Thuja Green Giant
3700	trash	willow-hybrid__trashed	Willow Hybrid
3701	publish	leyland-cypress	Leyland Cypress
3702	publish	american-holly	American Holly
3703	publish	italian-cypress	Italian Cypress

### product_variation
4992	publish	product-3702-variation-3	American Holly - #3 Gallon
5100	publish	product-3850-variation	Red Knock Out® Rose - Tree Form - #3 Gallon
5109	publish	product-3849-variation	Frost Proof Gardenia - #1 Gallon
5854	publish	5854	Carolina Sapphire Arizona Cypress - #3 Gallon
5858	publish	5858	Emerald Green Arborvitae - 1-2 Foot

### question_answer
10809	publish	how-fast-does-this-tree-grow-per-year	How fast does this tree grow per year?
10969	publish	how-far-apart-do-i-need-to-plant-these-for-a-privacy-screen	How far apart do I need to plant these for a privacy screen?
10998	publish	whats-the-difference-between-a-2-gallon-and-5-gallon	What's the difference between a 2 gallon and 5 gallon

### shop_coupon
30608	publish	10114218mccomb	10114218mccomb
30610	publish	10109477zhao	10109477zhao
31091	publish	10111378beebe	10111378beebe
31092	publish	10116523centofante	10116523centofante
32456	publish	gratitude10	gratitude10

### shop_order_placehold
873393	draft		
873394	draft		
873395	draft		
873396	draft		
873397	draft		

### shop_webhook
432813	publish	zendesk-webhook-on-new-order	Zendesk webhook on new order

### ttc_banner
801335	publish	ttc-banner-4	Your landscape, delivered.
812766	publish	ttc-banner-3	Privacy Tree Sale
820249	publish	ttc-banner-2	Got A Big Project?
823094	publish	ttc-banner	Free Shipping Today Only!
829618	publish	ttc-banner-5	Today Only - Free Shipping!

### ttc_guide
398454	publish	plant-ball-burlap-trees	How To Plant Ball and Burlap Trees
398536	publish	plant-bamboo-trees	How To Plant Bamboo Trees
398553	publish	plant-peach-nectarine-trees	How To Plant Peach and Nectarine Trees
405420	publish	plant-deciduous-trees	How to Plant Deciduous Trees
405431	publish	plant-bare-root-trees	How to Plant Bare Root Trees

### wafs
8527	publish	free-shipping	Free Shipping

### wp_global_styles
645017	publish	wp-global-styles-ttc	Custom Styles

### wpcf7_contact_form
199	publish	contact-form-1-2	Contact form 1
2888	publish	faq-form	FAQ Form
9716	publish	application-submission	Application Submission

### wpforms
52189	publish	existing-order-ticket	Shipped/Delivered Orders
53058	publish	simple-contact-form	General Questions
79541	publish	simple-contact-form-2	Price Match Inquiry
456975	publish	shipped-delivered-order	Not Shipped Orders
457026	publish	not-shipped-orders	General Question
