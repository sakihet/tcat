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
# post plain text
cat hello.go | tcat --topicId 123 --plain

# post text with fenced code blocks
echo -e "hello\ntypetalk" | tcat --topicId 123

# post code with syntax highlighting
cat hello.go | tcat --topicId 123 --syntax go
```
topic url: https://typetalk.in/topics/:topicId
