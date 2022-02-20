<?php

$host = '127.0.0.1';
$port = 6379;
$auth = '';

// $redis = new Redis();
// $redis->connect($host, $port);
// $redis->auth($auth);
// $redis->ping();
// $val = $redis->get('xyz');
// var_dump($val);
//
// $addr = '127.0.0.1:6379';
// $redis = new RedisCluster(null, [$addr], 1, 1, false, 'password');
// for ($i = 0; $i < 100; $i++) {
//     $redis->get("test_{$i}");
// }

if (is_file('./vendor-php/autoload.php')) {
    require './vendor-php/autoload.php';
    try {
        $client = new Predis\Client();
        $client->set('foo', 'bar');
        $value = $client->get('foo');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client([
            'scheme' => 'tcp',
            'host' => $host,
            'port' => $port,
            'password' => $auth,
        ]);
        $client->get('predis-test-ok');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client(['scheme' => 'unix', 'path' => '/path/to/redis.sock']);
        $client->get('predis-test-unix1');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client('unix:/path/to/redis.sock');
        $client->get('predis-test-unix2');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client(['scheme' => 'tls', 'ssl' => ['cafile' => 'private.pem', 'verify_peer' => true]]);
        $client->get('predis-test-tls1');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client('tls://127.0.0.1?ssl[cafile]=private.pem&ssl[verify_peer]=1');
        $client->get('predis-test-tls2');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client([
            'tcp://127.0.0.1:6379?alias=first-node', ['host' => '127.0.0.1', 'port' => 7381, 'alias' => 'second-node'],
        ], [
            'cluster' => 'predis',
        ]);
        $client->get('predis-test-cluster');
    } catch (Exception $e) {}

    try {
        $client = new Predis\Client(['tcp://127.0.0.1:6379'], [
            'cluster' => 'redis',
            'parameters' => ['password' => 'password'],
        ]);
        $client->set('predis-test-cluster-ok', "abc");
        echo $client->get('predis-test-cluster-ok'), PHP_EOL;
    } catch (Exception $e) {
        throw $e;
    }
}
