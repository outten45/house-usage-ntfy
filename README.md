# Development

With inotifywait installed, you can run the command:

```bash
while inotifywait -e close_write main.go; echo "------------------"; go run main.go; echo "--------------------"; end
```
