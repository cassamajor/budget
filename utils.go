package budget

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"math"
	"math/big"
	"net/http"
)

// Balance is a custom type for a float64 value.
type Balance float64

// Float64 returns the Balance value as a float64.
func (b Balance) Float64() float64 {
	return float64(b)
}

// String returns the Balance value as a string with a dollar sign and comma.
func (b Balance) String() string {
	return AddComma(b.Float64())
}

// AddComma adds a comma to a float64 value and prepends a dollar sign
func AddComma(n float64) string {
	p := message.NewPrinter(language.English)
	value := p.Sprintf("%.2f", math.Abs(n))

	// Otherwise, the negative sign will be placed after the dollar sign
	if n < 0 {
		return "-$" + value
	}
	return "$" + value
}

// NewRat converts an `int` into an `int64` and is fed into `big.NewRat` to avoid floating point errors.
func NewRat(a int) Balance {
	balance, _ := big.NewRat(int64(a), 1000).Float64()
	return Balance(balance)
}

// GetURL builds and executes a request to the specified YNAB API endpoint.
func GetURL(apiToken, url string) ([]byte, error) {
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set the API token in the request header
	bearer := fmt.Sprintf("Bearer %s", apiToken)
	req.Header.Set("Authorization", bearer)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading stream of data:", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}

	return body, nil
}
