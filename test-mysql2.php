<?php

$mysqli = mysqli_connect("127.0.0.1", "root", "", "testdb");
if (!$mysqli) {
    die("连接错误: " . mysqli_connect_error());
}

$mysqli->query("use testdb;");

$sql = "SELECT `UUID`,`PID`,`GID`,`CID`,`REFER`,`PAPID`,`DEV`,`VERSION`,`DOID`,`DSID`,`DEXT`,`DRID`,`DRNAME`,`DRLEVEL`,`DMONEY`,`UID`,`UNAME`,`MOID`,`PWAY`,`MONEY`,`SMONEY`,`STATE`,`REMARK`,`APID`,`BOID`,`CMONEY`,`BCARD`,`BCPWD`,`PAID`,`PTIME`,`CODE`,`HPID`,`IP` from sdk_order LIMIT 1";
$rs = $mysqli->query($sql);
while ($row = $rs->fetch_assoc()) {
    echo $row["UUID"], PHP_EOL;
}
$rs->free();
$mysqli->close();
