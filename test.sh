#! /bin/sh
go run ./cmd \
  -auth tests/.auth.txt \
  -dry \
  -emailcolname "email address" \
  -selectors ./tests/selectors.txt \
  -subscribers tests/subscribers.csv \
  ./tests/newsletter.txt
