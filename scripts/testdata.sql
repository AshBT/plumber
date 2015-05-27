-- MySQL dump 10.13  Distrib 5.6.19, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: memex_ht
-- ------------------------------------------------------
-- Server version	5.6.19-0ubuntu0.14.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `images`
--

DROP TABLE IF EXISTS `images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `images` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Auto incremented row identifier, unique to table.',
  `sources_id` int(11) unsigned NOT NULL COMMENT 'ID of the source the ad came from.',
  `ads_id` int(11) unsigned DEFAULT NULL COMMENT 'ID of the ad in the ads table that these images are from.',
  `url` varchar(2083) CHARACTER SET utf8mb4 NOT NULL COMMENT 'URL of the source image.',
  `location` varchar(128) DEFAULT NULL COMMENT 'Location of our cached copy of the image.',
  `importtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when image was imported into database.',
  `modtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Timestamp of the most recent modification.',
  PRIMARY KEY (`id`),
  UNIQUE KEY `url_unique` (`url`(191)),
  KEY `timestamp` (`importtime`),
  KEY `url` (`url`(128)),
  KEY `ads_id` (`ads_id`),
  KEY `location` (`location`(64)),
  KEY `sources_id` (`sources_id`),
  KEY `modtime` (`modtime`)
) ENGINE=InnoDB AUTO_INCREMENT=71324962 DEFAULT CHARSET=ascii;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `images`
--
-- WHERE:  ads_id in (1, 84, 5823, 23, 561, 903542)

LOCK TABLES `images` WRITE;
/*!40000 ALTER TABLE `images` DISABLE KEYS */;
INSERT INTO `images` VALUES (6245046,4,1,'http://www.foobar.com/p/8a32ad3567cb05f05d41bf18967eb7f0.jpg','https://www.foobar.com/63266a92559ce8535b80a71ad81db53eaf908c76.jpg','2014-06-10 16:54:14','2014-07-15 21:12:26'),(6245047,4,1,'http://www.foobar.com/p/992882507616e6125b1bc98927392b1e.jpg','https://www.foobar.com/68b019ba27a82a6e0738d060f83fa72db7a25567.jpg','2014-06-10 16:54:14','2014-07-15 21:13:10'),(6245048,4,1,'http://www.foobar.com/p/47d54e1275de633f76574b225cf39461.jpg','https://www.foobar.com/9dcfd2b6459a44c44381bd55be9a858b3301dcd9.jpg','2014-06-10 16:54:14','2014-07-15 21:13:10'),(6245049,4,1,'http://www.foobar.com/p/6cd94b9fb05ad53e906c5743292ffaf7.jpg','https://www.foobar.com/e6088c37fadb0626deb421494843fa912be600f1.jpg','2014-06-10 16:54:14','2014-07-15 21:13:10'),(6245114,4,23,'http://www.foobar.com/p/fec585148b3361016e6fdac8f50e861d.jpg','https://www.foobar.com/0c2d999edbcbddad1349e62e44788fae7fa3561f.jpg','2014-06-10 16:54:15','2014-07-16 00:29:39'),(6245115,4,23,'http://www.foobar.com/p/a7843bfe3bd3416d53fe6254564c1be9.jpg','https://www.foobar.com/a96e55bab4d75da84ca375b5a980522afc9edba2.jpg','2014-06-10 16:54:15','2014-07-16 00:29:39'),(6245116,4,23,'http://www.foobar.com/p/6e4b559f8b23c94ddb767f4bea00859f.jpg','https://www.foobar.com/a72c789777ec09890291831eb7a25976ce3ed81f.jpg','2014-06-10 16:54:15','2014-07-16 00:29:39'),(6245374,4,84,'http://www.foobar.com/p/28c34dd695c5fabc3568ccf617422499.jpg','https://www.foobar.com/f27ecae2c2f5b1719b8a2102bda970fd1b5b5033.jpg','2014-06-10 16:54:19','2014-07-16 00:29:43'),(6245375,4,84,'http://www.foobar.com/p/bea0c5dc4086b8c8f73accd1ced0d55e.jpg','https://www.foobar.com/8b2ea33d585e2b6b803a60aecccb768d89925baf.jpg','2014-06-10 16:54:19','2014-07-16 00:29:44'),(6245376,4,84,'http://www.foobar.com/p/90835d7ed5eb4c2989f038a0ddc4dc74.jpg','https://www.foobar.com/e041c37f7d0eb2f8c2c3916d7d9ba1be4f3b6983.jpg','2014-06-10 16:54:19','2014-07-16 00:29:44'),(6245377,4,84,'http://www.foobar.com/p/8893aa34a5d813f935f9d9328db81ae5.jpg','https://www.foobar.com/cfde19817425d39bd31f0f9578d9194c574f973e.jpg','2014-06-10 16:54:19','2014-07-16 00:29:44'),(6245378,4,84,'http://www.foobar.com/p/fdf33908941faa7c8c2bd6b45da040f0.jpg','https://www.foobar.com/ddcb6036094affc288f966aab1650074510ca0e9.jpg','2014-06-10 16:54:19','2014-07-16 00:29:44'),(6247010,4,561,'http://www.foobar.com/p/5eca2ebdd24ae3f5905c939712285870.jpg','https://www.foobar.com/a0353b184faddbb2f3e8f02d9592f8ef893e7967.jpg','2014-06-10 16:54:49','2014-07-16 00:30:15'),(6247011,4,561,'http://www.foobar.com/p/1d6ae58869732bc7e4af3d65afc699a9.jpg','https://www.foobar.com/5a25e401f828aa029fb245b49742bbbb788dc613.jpg','2014-06-10 16:54:49','2014-07-16 00:30:15'),(5879551,2,903542,'http://www.foobar.com/00M0M_hCu6gBYFwvv_600x450.jpg',NULL,'2014-06-10 14:12:33','2014-06-15 18:48:12');
/*!40000 ALTER TABLE `images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ads_attributes`
--

DROP TABLE IF EXISTS `ads_attributes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ads_attributes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID of the attribute entry',
  `ads_id` int(11) unsigned NOT NULL COMMENT 'Parent ID for this attribute.',
  `attribute` varchar(32) NOT NULL COMMENT 'Attribute name (Age, location, etc.)',
  `value` varchar(2500) CHARACTER SET utf8 NOT NULL COMMENT 'Value of the attribute.',
  `extracted` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT 'If no the value was from the structure of the website. If yes we used an algorithm on the text to get the value and it may be less accurate.',
  `extractedraw` varchar(512) CHARACTER SET utf8 DEFAULT NULL COMMENT 'Raw text of the data if extracted.',
  `modtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Timestamp of the most recent modification.',
  PRIMARY KEY (`id`),
  KEY `ads_id` (`ads_id`),
  KEY `attribute` (`attribute`(4)),
  KEY `extracted` (`extracted`),
  KEY `value` (`value`(16)),
  KEY `modtime` (`modtime`)
) ENGINE=InnoDB AUTO_INCREMENT=501604310 DEFAULT CHARSET=ascii;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ads_attributes`
--
-- WHERE:  ads_id in (1, 84, 5823, 23, 561, 903542)

LOCK TABLES `ads_attributes` WRITE;
/*!40000 ALTER TABLE `ads_attributes` DISABLE KEYS */;
INSERT INTO `ads_attributes` VALUES (126114125,1,'phone','5555559574',1,' 5555559574','2014-07-15 21:05:35'),(126114136,1,'height','165',0,NULL,'2014-07-15 21:05:36'),(358026106,23,'username','janedoe92',0,NULL,'2014-10-27 17:25:51'),(358026109,23,'reviewsite1','http://www.foobar.com/reviews/show.asp?id=208894',0,NULL,'2014-10-27 17:25:51'),(358026117,23,'weight','115',0,NULL,'2014-10-27 17:25:51'),(358026122,23,'longitude','0',0,NULL,'2014-10-27 17:25:51'),(358026127,23,'email','janedoe@foobar.com',1,NULL,'2014-10-27 17:25:51'),(358026131,23,'build','Fail',0,NULL,'2014-10-27 17:25:51'),(358026137,23,'latitude','0',0,NULL,'2014-10-27 17:25:51'),(358026140,23,'availability','none',0,NULL,'2014-10-27 17:25:51'),(358026145,23,'ethnicity','ethnicity',1,NULL,'2014-10-27 17:25:51'),(299941033,84,'username','johndoe',0,NULL,'2014-10-22 21:29:49'),(299941052,84,'eyes','Brown',0,NULL,'2014-10-22 21:29:50'),(299941063,84,'weight','140',0,NULL,'2014-10-22 21:29:50'),(299941077,84,'cup','stanley',0,NULL,'2014-10-22 21:29:50'),(299941098,84,'bust','boom',0,NULL,'2014-10-22 21:29:50'),(299941110,84,'height','5\'6\'\'',0,NULL,'2014-10-22 21:29:50'),(299941124,84,'hair','Bald',0,NULL,'2014-10-22 21:29:50'),(299941135,84,'email','johndoe@foobar.com',1,NULL,'2014-10-22 21:29:50'),(299941150,84,'build','Success',0,NULL,'2014-10-22 21:29:50'),(299941156,84,'availability','sometimes',0,NULL,'2014-10-22 21:29:50'),(299941161,84,'ethnicity','ethnicity',1,NULL,'2014-10-22 21:29:50'),(358026133,561,'username','janedoe1776',0,NULL,'2014-10-27 17:25:51'),(358026136,561,'eyes','Brown',0,NULL,'2014-10-27 17:25:51'),(358026141,561,'longitude','0',0,NULL,'2014-10-27 17:25:51'),(358026149,561,'reviewsite1','http://www.foobar.com/site_listing/reviewed_seal.gif',0,NULL,'2014-10-27 17:25:51'),(358026152,561,'weight','165',0,NULL,'2014-10-27 17:25:51'),(358026156,561,'cup','mushroom',0,NULL,'2014-10-27 17:25:51'),(358026162,561,'waist','32',0,NULL,'2014-10-27 17:25:51'),(358026169,561,'bust','boom',0,NULL,'2014-10-27 17:25:51'),(358026176,561,'height','5\'4\'\'',0,NULL,'2014-10-27 17:25:51'),(358026181,561,'hair','Brown',0,NULL,'2014-10-27 17:25:51'),(358026194,561,'email','janedoe1776@foobar.com',1,NULL,'2014-10-27 17:25:51'),(358026203,561,'build','Error',0,NULL,'2014-10-27 17:25:51'),(358026213,561,'latitude','0',0,NULL,'2014-10-27 17:25:51'),(358026218,561,'availability','none',0,NULL,'2014-10-27 17:25:51'),(358026224,561,'hips','do not lie',0,NULL,'2014-10-27 17:25:51'),(358026231,561,'phone','5555552667',1,'       555-555-2667','2014-10-27 17:25:51'),(98116507,5823,'phone','5555557218',1,' 555-555-7218','2014-07-07 16:55:16');
/*!40000 ALTER TABLE `ads_attributes` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-05-27  2:16:08
