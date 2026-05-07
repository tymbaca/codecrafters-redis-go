[tester::#XZ2] Running tests for Stage #XZ2 (AOF Persistence - Replay a single command)
[tester::#XZ2] Creating append-only directory "blueberry":
[tester::#XZ2] [blueberry]   - strawberry.aof.manifest
[tester::#XZ2] [blueberry]   - pineapple.aof.1.incr.aof
[tester::#XZ2] Writing 1 command to append-only file "pineapple.aof.1.incr.aof"
[tester::#XZ2] [pineapple.aof.1.incr.aof] *3\r\n$3\r\nSET\r\n$5\r\napple\r\n$3\r\n329\r\n
[tester::#XZ2] Creating manifest file "strawberry.aof.manifest"
[tester::#XZ2] [strawberry.aof.manifest] file pineapple.aof.1.incr.aof seq 1 type i
[tester::#XZ2] $ ./your_program.sh --dir /tmp/aof-8712 --appendonly yes --appenddirname blueberry --appendfilename strawbe
rry.aof --appendfsync always
[your_program] {"time":"2026-05-07T15:12:05.42623831Z","level":"DEBUG","msg":"no records in manifest, creating first AOF f
ile"}
[tester::#XZ2] Checking if the command in append-only file was replayed
[tester::#XZ2] [client] $ redis-cli GET apple
[tester::#XZ2] Expected bulk string, found null bulk string ($-1\r\n)
[tester::#XZ2] Test failed
[your_program] {"time":"2026-05-07T15:12:05.524471021Z","level":"ERROR","msg":"handle conn failed","err":"read value heade
r: read tcp 127.0.0.1:6379->127.0.0.1:34714: use of closed network connection"}
