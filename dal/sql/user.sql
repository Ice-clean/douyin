/*
SQLyog Ultimate v13.1.1 (64 bit)
MySQL - 8.0.27 : Database - douyin
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`douyin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `douyin`;

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `password` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `name` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `create_time` bigint DEFAULT NULL,
  `last_login` bigint DEFAULT NULL,
  `freeze` tinyint DEFAULT '0',
  `age` int DEFAULT '20',
  `personal_signature` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '你们喜欢的话题，就是我们采访的内容',
  `site` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '杭州',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `user` */

insert  into `user`(`id`,`username`,`password`,`name`,`create_time`,`last_login`,`freeze`,`age`,`personal_signature`,`site`) values 
(2,'aaaaaa','111111','张三',1653052639043,1653053036270,0,20,'你们喜欢的话题，就是我们采访的内容','杭州'),
(3,'bbbbb','22222','张三',1653054724764,1653054724764,0,20,'你们喜欢的话题，就是我们采访的内容','杭州'),
(4,'ccccc','33333','张三',1653055529799,1653055529799,0,20,'你们喜欢的话题，就是我们采访的内容','杭州');

/*Table structure for table `user_follow` */

DROP TABLE IF EXISTS `user_follow`;

CREATE TABLE `user_follow` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `follow_user_id` bigint DEFAULT NULL,
  `followed_user_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `user_follow` */

insert  into `user_follow`(`id`,`follow_user_id`,`followed_user_id`) values 
(1,3,2),
(2,2,3);

/*Table structure for table `user_follow_count` */

DROP TABLE IF EXISTS `user_follow_count`;

CREATE TABLE `user_follow_count` (
  `id` bigint NOT NULL,
  `follow_count` bigint DEFAULT '0',
  `follower_count` bigint DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

/*Data for the table `user_follow_count` */

insert  into `user_follow_count`(`id`,`follow_count`,`follower_count`) values 
(2,1,1),
(3,1,1),
(4,0,0);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
