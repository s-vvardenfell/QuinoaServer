# QuinoaServer<br>

[![build-test](https://github.com/s-vvardenfell/QuinoaServer/actions/workflows/build-test.yml/badge.svg)](https://github.com/s-vvardenfell/QuinoaServer/actions/workflows/build-test.yml) <br>

Server for Quinoa project that connects Parser, TgBot and Cache<br>

Config example:<br>
```yaml
server_addr: localhost
server_port: "port"
parser_port: "port"
redis_serv_port: "port"
with_reflection: false
exp_time: 60
logrus: 
  log_level: 4
  to_file: false
  to_json: false
  log_dir: "logs/logs.log"
```
