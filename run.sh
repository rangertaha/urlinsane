
for i in $(seq 0 5); do
    go run ./cmd typo -e tomlee@google.com
    go run ./cmd typo -n tomlee@google.com

done
