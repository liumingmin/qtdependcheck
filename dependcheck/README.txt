dependcheck.exe {需要校验的构建结果目录} {MD5列表文件路径}

echo %errorlevel%
成功返回0， 入参错误返回-1 文件校验失败返回-2