# pricing

## Design

### Long term design

- The plan service will host the management of multiple plans across countries.
- The customer service has the management of customer data.
- Customer service gets the pricing information from the plan service.
- When plan updates are pushed using plan service, once the prices are changed, customers need to be notified using the notification service.
  - This can be done using the country. For a specific country and plan for which the prices are changed, all customers needs to be notified/
  - The change needs to be communicated to the billing service in order to make any updates necessary.
- The service can be extended to have price changes start into effect after a certain period of time rather than effective immediately.
- For efficient retrieval of prices, there can be a cache on top of plan service, which can easily retrieve the prices of plans for customers for that country.


### Current design

- In the current implementation, only a single webservice is implemented which hosts both customer service and plan service.
- The information is saved in files for customers and plans. The plan data is in plans.json, supported plan information is in supported_plans.json and customer data is in customers.json
- The plan information is indexed by country to efficiently get information about particular plans.
- The data access layer is implemented by interfaces to easily change it to use DB in future without changing much code.
- The customer service talks to the plan service to get the price for the countryCode and planName pair.
- The webservice is written in Go and it exposes URLs to manage plans and customers.
- For changing plans, there is only one URL which changes plans by country.
- In future, the design can be extended for customer notification and other updates by doing that asynchronously.


### Assumptions for the implementation

- With each country, customer information is saved. In the current implementation, the customer data is saved in file but it can be saved in DB and indexed by countryCode to efficiently get customers for particular country.
- Once the price is changed for a particular country, the change is reflected immediately on the customer side as well. In actual system, there will be notification services that will be used to notify the customers for the price change and notify the billing service as well. But in the current implementation, there is no notification mechanism.
- In the current implementation, the price is not saved along with each customer. That is just redundant information as the price is saved with the plan information and with each customer, plan information is saved, using the plan service, the price can be fetched.
  - If the price would have been saved for each customer, we need to keep that in sync with the prices in the plan service.
- In the current implementation, the plan information is also saved in the file in a map like format with countryCode as key.
- The current implementation does not handle multiple updates going on in parallel as they are being done on a file. It can easily be added by adding a RWMutex mechanism. But if we change to use a DB, the database libraries internally provide that functionality.
- The price plan service can be hosted as a separate micro-service. but in the current implementation, both the plan service and customer service are running as a single webservice.
- Once the price is changed for a plan, the implementation keeps the last price saved for history purposes. In future, this can be separated out to keep a separate history record for audit purposes.
- The current implementation only exposes URL to change plan price by country, for a global change, the same URL can be called multiple times for every country. The main reason behind this implementation is that the pricing may not be globally be the same proportionate to the exchange rate.


## Compilation and Deployment

 - For compilation, installation Go 1.10.2 on your computer.
 - set GOPATH environment variable to the root directory of this repository and cd to the root directory of the repository.
 - Run command `go get github.com/gorilla/mux`
 - Run command to build `go build pricing`. It will create a binary called pricing in the root directory.
 - Run the binary while in the root repo directory, ./pricing.
 - This will start the webservice at port `8080` and now you should be able to call the REST APIs in that service.


## REST API documentation

###  `GET http://localhost:8080/pricing/api/v1/plans`

- Returns all the plans for all countries. It is in the following format. Each country has a field for the currency, the price of each plan, valid_from shows the timestamp since the plan price was valid and the old price. If old price is 0, then it is the first price.
```
{"IND":{"currency":"INR","plan_price":{"1S":{"price":500,"valid_from":0,"old_price":0},"2S":{"price":1000,"valid_from":0,"old_price":0},"4S":{"price":1500,"valid_from":0,"old_price":0}}}
,"UK":{"currency":"GBP","plan_price":{"1S":{"price":10,"valid_from":0,"old_price":0},"2S":{"price":15,"valid_from":0,"old_price":0}}}}
```
- Returns 200 for success.
- Returns 400 if there is any error internally.


### `GET http://localhost:8080/pricing/api/v1/plans/<countryCode>`

- Returns all the plans and their pricing for the specific country code. The country code is for example, USA/UK/IND etc. The output is in the following format-
```
{"currency":"INR","plan_price":{"1S":{"price":500,"valid_from":0,"old_price":0},"2S":{"price":1000,"valid_from":0,"old_price":0},"4S":{"price":1500,"valid_from":0,"old_price":0}}}
```
- Returns 200 if success
- Returns 400 and error message if country code is invalid.

### `GET http://localhost:8080/pricing/api/v1/plans/<countrycode>/<planname>`

- Returns the plan details for a specific plan in a specific country. The plan name could be 1S, 2S, 4S. The output is in following format-
```
{"plan_name":"1S","plan_info":"1 stream plan","price":500,"valid_from":0,"old_price":0}
```
- Returns 200 for success
- Returns 400 if country code or planname is invalid.

### `PUT http://localhost:8080/pricing/api/v1/plans/<countryCode>`

- This is used to change the prices of plans in specific country code. This accepts a request body in the following format
```
{
  "1S":15,
  "2S":20,
  "4S":30
}
```
This is a json with string key and int values, the string key is the plan name and the int value is the new price for that plan in that country.
If a country does not support a plan yet, this API can also be used to add a plan to that country. But that can only be from the list of supported plans, which is 1S, 2S or 4S.

- This returns 200 if the update is successful.
- It returns 400, if country code or planName is invalid or if update fails.

### `GET http://localhost:8080/pricing/api/v1/customers`

- Returns list of all customers. This returns the email, plan, country and price for every customer. Following is the response format-
```
[{"email":"a@abcusa.com","countryCode":"USA","planName":"1S","price":15},{"email":"b@abcusa.com","countryCode":"USA","planName":"2S","price":20}]
```
- Returns 200 if success
- Returns 400 if there is any error.

### `GET http://localhost:8080/pricing/api/v1/customers/<countryCode>`

- Returns list of all customers for a specific country. The format is the following
```
[{"email":"a@abcusa.com","countryCode":"USA","planName":"1S","price":15},{"email":"b@abcusa.com","countryCode":"USA","planName":"2S","price":20}]
```
- Returns 200 if success.
- Returns 400 if the countryCode is invalid.

