# tcat
## Configuration
```
export TYPETALK_API_CLIENT_ID=""
export TYPETALK_API_CLIENT_SECRET=""
```
## Build
```
go build
mv ./tcat /usr/local/bin/ # or somewhere
```
## Usage
```
echo -e "hello\ntypetalk" | tcat --topicId 123
```
