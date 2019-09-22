package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/tranduythanh/payment"
)

var domesticPayment *payment.OnePayDomestic
var internationalPayment *payment.OnePayInternational

func main() {
	domesticPayment = payment.NewSandboxDomestic("https://a98f94fa.ngrok.io/payment/callback/domestic")
	internationalPayment = payment.NewSandboxInternational("https://a98f94fa.ngrok.io/payment/callback/international")

	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/payment/checkout/domestic", checkoutDomestic)
	e.GET("/payment/checkout/international", checkoutInternational)

	e.GET("/payment/callback/domestic", callbackDomestic)
	e.GET("/payment/callback/international", callbackInternational)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.HTML(http.StatusOK, `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

    <title>Hello, world!</title>
  </head>
  <body>

	<h2>Checkout</h2>
	<a class="btn btn-primary" href="/payment/checkout/domestic" role="button">Domestic</a>
	<a class="btn btn-primary" href="/payment/checkout/international" role="button">International</a>

    <!-- Optional JavaScript -->
    <!-- jQuery first, then Popper.js, then Bootstrap JS -->
    <script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
  </body>
</html>
	`)
}

func checkoutDomestic(c echo.Context) error {
	timeStr := fmt.Sprintf("%d", time.Now().UnixNano())
	url, err := domesticPayment.BuildCheckoutURL(&payment.CheckoutParams{
		Amount:      100000,
		OrderInfo:   timeStr,
		MerchTxnRef: timeStr,
		TicketNo:    "::1",
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(url)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func checkoutInternational(c echo.Context) error {
	timeStr := fmt.Sprintf("%d", time.Now().UnixNano())
	url, err := internationalPayment.BuildCheckoutURL(&payment.CheckoutParams{
		Amount:      100000,
		OrderInfo:   timeStr,
		MerchTxnRef: timeStr,
		TicketNo:    "::1",
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(url)

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func callbackDomestic(c echo.Context) error {
	v, err := domesticPayment.HandleCallback(c.QueryParams())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &v)
}

func callbackInternational(c echo.Context) error {
	v, err := internationalPayment.HandleCallback(c.QueryParams())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &v)
}
