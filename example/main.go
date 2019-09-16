package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tranduythanh/payment"
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/payment/checkout", checkout)
	e.GET("/payment/callback", callback)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func hello(c echo.Context) error {
	return c.HTML(http.StatusOK, `
<!DOCTYPE html>
<html>
<body>

<h2>Demo</h2>
<p><a href="/payment/checkout">Checkout</a></p>

</body>
</html>
	`)
}

func checkout(c echo.Context) error {
	url, err := payment.NewSandboxDomestic("http://6787db72.ngrok.io/payment/callback").BuildCheckoutURL(&payment.CheckoutParams{
		Amount:      100000,
		OrderInfo:   "ahihi",
		MerchTxnRef: "behihi",
		TicketNo:    "::1",
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(url)

	return c.Redirect(http.StatusPermanentRedirect, url)
}

func callback(c echo.Context) error {
	var v = &payment.DomesticResponse{}

	err := c.Bind(v)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	v.Norm()

	return c.JSON(http.StatusOK, &v)
}
