
![yourITcity](https://storage.gra5.cloud.ovh.net/v1/AUTH_b49a48a77ecf40ff9e51a30fdf28633c/yic/yourITcity%201%20transparent%20300%20x%20324.png?temp_url_expires=1570545027&temp_url_sig=6a6d20d61404e29f9e132fb49ab09cf6abd4e778)

# Email sender service

This service wait on the RabbitMQ __queue__ called `email`.

For each message received to sent the it to smpt server.
The messages must have Content type `text/plain` or `text/html` and got `To` and `Subject` header fields. The body of the message will be sent as the content of the email.

Configuration by Environment variables

``` sh
    RABBITMQ_URI = "amqps://guest:guest@localhost:5672"
    EMAIL_FROM="no-reply@youritcity.io"
    MJ_APIKEY_PUBLIC="MJ_APIKEY_PUBLIC"
    MJ_APIKEY_PRIVATE="MJ_APIKEY_PRIVATE"

```
## Mailjet Account

You' ll need to have an Mailjet account to send emails and generate API keys
Go to [Mailjet](https://dev.mailjet.com)
## Compile & run

It use [go modules](https://blog.golang.org/using-go-moduleshttps://blog.golang.org/using-go-modules) to handler dependancy.

```sh
    go build
    ./yic_email
```

## Use fake email sender (print email to stdout)

```sh
    go run cmd/dev-email/main.go amqps://guest:guest@localhost:5672
```

## Send email

```sh
    go run cmd/send/main.go guest:guest@localhost:5672 email@address.com ["Subject with space"]
```

Then you can write the content of the email and press Ctrl+D to send or

```sh
    go run cmd/send/main.go guest:guest@localhost:5672 email@address.com ["Subject with space"] < content_file
```
