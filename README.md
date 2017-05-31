# tcat
## Build
```
go build
mv ./tcat /usr/local/bin/ # or somewhere
```
## Configuration
tcat uses `$HOME/.tcat` for configuration.
```
tcat --configure
```
## Usage
```
echo -e "hello\ntypetalk" | tcat --topicId 123
```
