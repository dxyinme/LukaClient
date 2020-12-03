package db

var (

	// about create table
	CREATE_USERINFO_TABLE = "CREATE TABLE IF NOT EXISTS userInfo (UID CHAR(64) PRIMARY KEY NOT NULL, name CHAR(64) NOT NULL)"
	CREATE_MSG_TABLE = "CREATE TABLE IF NOT EXISTS msgTable (msgId CHAR(134) PRIMARY KEY NOT NULL," +
		"msgType INT NOT NULL," +
		"msgContentType INT NOT NULL," +
		"content TEXT NOT NULL," +
		"sendTime TEXT NOT NULL," +
		"msg_from CHAR(64) NOT NULL," +
		"msg_target CHAR(64) NOT NULL," +
		"groupName CHAR(64) NOT NULL)"

	// insert
	INSERT_USERINFO = "INSERT INTO userInfo (UID, name) values (?,?)"
	INSERT_MSG = "INSERT INTO msgTable (msgId,msgType,msgContentType,content,sendTime,msg_from,msg_target,groupName) " +
		"values (?,?,?,?,?,?,?,?)"

	// select
	SELECT_MSG_BY_GROUP_ALL = "SELECT * FROM msgTable WHERE groupName = ? ORDER BY msgId DESC"
	SELECT_MSG_BY_GROUP = SELECT_MSG_BY_GROUP_ALL + " LIMIT ? "
	SELECT_MSG_BY_FROM_ALL = "SELECT * FROM msgTable WHERE msg_from = ? ORDER BY msgId DESC"
	SELECT_MSG_BY_FROM = SELECT_MSG_BY_FROM_ALL + " LIMIT ? "
)
