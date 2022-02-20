<?php

$header = ["X-Test-Trace: HelloWorld"];
$ch = curl_init();
curl_setopt_array($ch, [
    // CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
    CURLOPT_RETURNTRANSFER => 1,
]);
curl_setopt($ch, CURLOPT_URL, "http://im.qq.com/");
curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
// curl_setopt($ch, CURLOPT_VERBOSE, 1);
curl_exec($ch);
// $a = curl_getinfo($ch);
curl_close($ch);

$ch = curl_init();
curl_setopt_array($ch, [
    CURLOPT_RETURNTRANSFER => 1,
]);
curl_setopt($ch, CURLOPT_URL, "http://www.hao123.com/");
curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
curl_exec($ch);
curl_close($ch);

var_dump($_SERVER);