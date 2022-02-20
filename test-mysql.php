<?php

$mysqli = mysqli_connect("127.0.0.1", "root", "", "testdb");
if (!$mysqli) {
    die("连接错误: " . mysqli_connect_error());
}

// $mysqli = new Mysqli("127.0.0.1", "root", "", "testdb");
// if (!$mysqli) {
//     die("连接错误: " . mysqli_connect_error());
// }

// $rs = mysqli_query($mysqli, "SELECT ID, SERVER, MEMO, RATE FROM table_test");
// $rs = $mysqli->query("SELECT ID, SERVER, MEMO, RATE FROM table_test");
// while ($row = $rs->fetch_assoc()) {
//     echo $row["ID"] . "," . $row["SERVER"] . "," . $row["MEMO"] . "," . $row["RATE"], PHP_EOL;
// }
// $rs->free();

$stmt = $mysqli->prepare("SELECT ID, SERVER, MEMO, RATE FROM table_test WHERE ID IN (?,?)");
$stmt->bind_param("ii", $id1, $id3);
$id1 = 1;
$id3 = 3;
// $stmt->execute();
mysqli_stmt_execute($stmt);
$stmt->bind_result($a1, $a2, $a3, $a4);
while ($stmt->fetch()) {
    echo $a1 . "," . $a2 . "," . $a3 . "," . $a4, PHP_EOL;
}
$mysqli->close();
