use goweb;

CREATE TABLE IF NOT EXISTS `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(60) NOT NULL UNIQUE,
    `pwd` text NOT NULL,
    `time` int(11) DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=MyISAM  DEFAULT CHARSET=utf8 ;

