# GraphQL API

This is a solution for an [example test assignment](https://gist.github.com/alexesDev/4d9766bce86940106a1d97ee010c0c83).

To run the server use the following commands:
```
docker-compose up -d
dbmate up
cat ./db/seed.sql | psql postgres://staging:staging@localhost:7232
go run .
```

To use [SMS service](https://sms.ru/) specify the value for `SMS_API_ID` variable in the `.env` file.