# Communications Platform Coding Challenge

## Task

### Overview
In order to prevent downtime during an email service provider outage, you’re tasked with
creating a service that provides an abstraction between two different email service providers.
This way, if one of the services goes down, you can quickly failover to a different provider
without affecting your customers.

### Specifications:
Please create an HTTP service that accepts POST requests with JSON data to a ‘/email’
endpoint with the following parameters:
* ‘to’ The email address to send to
* ‘to_name’ The name to accompany the email
* ‘from’ The email address in the from and reply fields
* ‘from_name’ the name to accompany the from/reply emails
* ‘subject’ The subject line of the email
* ‘body’ the HTML body of the email

Example Request Payload:
`````
{
“to”: “fake@example.com”,
“to_name”: “Ms. Fake”,
“from”: “noreply@uber.com”,
“from_name”: “Uber”,
“subject”: “A Message from Uber”,
“body”: “<h1>Your Bill</h1><p>$10</p>”
}
`````

Your service should then do a bit of data processing on the request:
* Do the appropriate validations on the input fields (NOTE: all fields are required).
* Convert the ‘body’ HTML to a plain text version to send along to the email provider. You
can simply remove the HTML tags. Or if you’d like, you can do something smarter.

Once the data has been processed and meets the validation requirements, it should send the
email by making an HTTP request (don’t use SMTP) to one of the following two services:
- Sendgrid www.sendgrid.com (API Documentation: https://sendgrid.com/docs/)
- Postmark www.postmark.com (API Documentation: https://postmarkapp.com/developer)
  
Both services are free to try and are pretty painless to sign up for, so please register your own
  test accounts on each.
  Your service should send emails using one of the two options by default, but a simple
  configuration change and/or a redeploy of the service should switch it over to the other provider.

### Implementation Requirements:
* Please do not use the client libraries provided by Postmark or Sendgrid. In both cases,
  you’re making a simple post request. Do this with a lower level package or your
  language’s builtin commands.
* Stay away from big frameworks like Rails or Django. Feel free to use a microframework
  like Flask, Sinatra, Express, Silex, etc.
* This is a simple exercise, but organize, design, document and test your code as if it were
  going into production.
* Most of our back end is in Ruby, or Node.JS. If you know one of these languages,
  please use it. If not, that’s totally OK; write up your service in whatever language you’re
  most comfortable in.
* When you’re finished, post your code to github / gitlab / bitbucket and send us a link so
  we can check it out. Please don’t commit your email provider API keys in the repository.
  * Please include a README file in your repository with the following information:
  * How to install your application
  * Which language and/or microframework you chose and why
  * Tradeoffs you might have made, anything you left out, or what you might do
  differently if you were to spend additional time on the project
  * Anything else you wish to include
    
_____________

**How to install**

It is installed like any other Go application.
You must add environment variables `postmark-key`(with postmark server key for example '12345678-abcd-ef90-1234-1234567890as') and `sendgrip-key`(with sendgrip api key for example 'SG.1234567890asdfghjklo12.1234567890qwertyuioasdfghjklzxcvbnm12345678')

**Which language and/or microframework you chose and why**

I choose Go, because it is small, fast, and really easy to do micro-services like this one. Also, I've been working with this language the last year, so I did the project quite fast and easy.

**Tradeoffs you might have made, anything you left out, or what you might do differently if you were to spend additional time on the project**

I left out some tests because they require time, but I do the toughest ones.

If I'd spend additional time on the project I would add logs and metrics. For instance, in Mercadolibre we use Datadog for metrics, and ELK stack for logging. It is very useful. With them you can configure Opsgenie alerts.

Also, I'd improve error management and fix the different 'fixme' or 'todo' I left in the code. Last but not least, add some documentation.

**Anything else you wish to include**

Although the document said 'Your service should send emails using one of the two options by default, but a simple
configuration change and/or a redeploy of the service should switch it over to the other provider.', I did a little different.
If the default provider fails, it will try to send it by another provider. But if you don't like this approach it can be easily changed in main.go
