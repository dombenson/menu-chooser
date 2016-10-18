CREATE TABLE `users` (
  `userid` INT NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(64) NOT NULL,
  `pass` VARCHAR(64) NULL,
  `name` VARCHAR(64) NULL,
  PRIMARY KEY (`userid`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB;

CREATE TABLE `events` (
  `eventid` INT NOT NULL AUTO_INCREMENT,
  `userid` INT NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `location` VARCHAR(64) NOT NULL,
  `date` DATETIME NOT NULL,
  PRIMARY KEY (`eventid`),
  FOREIGN KEY (`userid`) REFERENCES users(userid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `courses` (
  `courseid` INT NOT NULL AUTO_INCREMENT,
  `eventid` INT NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `order` TINYINT NOT NULL,
  PRIMARY KEY (`courseid`),
  FOREIGN KEY (`eventid`) REFERENCES events(eventid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `options` (
  `optionid` INT NOT NULL AUTO_INCREMENT,
  `courseid` INT NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `description` TEXT NULL,
  PRIMARY KEY (`optionid`),
  FOREIGN KEY (`courseid`) REFERENCES courses(courseid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `attendees` (
  `attendeeid` INT NOT NULL AUTO_INCREMENT,
  `eventid` INT NOT NULL,
  `email` VARCHAR(64) NOT NULL,
  `name` VARCHAR(64) NOT NULL,
  `loginkey` VARCHAR(64) NOT NULL,
  PRIMARY KEY (`attendeeid`),
  FOREIGN KEY (`eventid`) REFERENCES events(eventid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `selections` (
  `selectionid` INT NOT NULL AUTO_INCREMENT,
  `attendeeid`  INT NOT NULL,
  `courseid`    INT NOT NULL,
  `optionid`    INT NOT NULL,
  PRIMARY KEY (`selectionid`),
  UNIQUE KEY `person_course` (`attendeeid`, `courseid`),
  FOREIGN KEY (`attendeeid`) REFERENCES attendees (attendeeid)
    ON DELETE CASCADE ON UPDATE CASCADE ,
  FOREIGN KEY (`courseid`) REFERENCES courses (courseid)
    ON DELETE CASCADE ON UPDATE CASCADE ,
  FOREIGN KEY (`optionid`) REFERENCES options (optionid)
    ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE =InnoDB;

CREATE TRIGGER `insert_selection_course` BEFORE INSERT ON `selections` FOR EACH ROW SET NEW.courseid = (SELECT courseid FROM options WHERE optionid = NEW.optionid);
CREATE TRIGGER `update_selection_course` BEFORE UPDATE ON `selections` FOR EACH ROW SET NEW.courseid = (SELECT courseid FROM options WHERE optionid = NEW.optionid);

CREATE TABLE `resets` (
  `resetid` INT NOT NULL AUTO_INCREMENT,
  `userid` INT NOT NULL,
  `date` DATETIME NOT NULL DEFAULT NOW(),
  `key` VARCHAR(64) NOT NULL,
  `used` BOOL NOT NULL,
  PRIMARY KEY (`resetid`),
  FOREIGN KEY (`userid`) REFERENCES users(userid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB;

