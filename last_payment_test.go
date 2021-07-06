package sendios

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestSendLastPayment(t *testing.T) {
	require.NoError(t, godotenv.Load())

	require.NotEmpty(t, os.Getenv("SENDIOS_PROJECT"))
	require.NotEmpty(t, os.Getenv("SENDIOS_CLIENT_ID"))
	require.NotEmpty(t, os.Getenv("SENDIOS_TOKEN"))
	require.NotEmpty(t, os.Getenv("TEST_EMAIL"))

	project, err := strconv.Atoi(os.Getenv("SENDIOS_PROJECT"))
	require.NoError(t, err)
	client, err := strconv.Atoi(os.Getenv("SENDIOS_CLIENT_ID"))
	require.NoError(t, err)

	c := New(project, client, os.Getenv("SENDIOS_TOKEN"))
	c.DebugRequests = true

	payment, err := c.SendLastPaymentByEmail(os.Getenv("TEST_EMAIL"), Payment{
		StartDate:   fmt.Sprintf("%d", time.Now().AddDate(0, -1, 0).Unix()),
		ExpireDate:  fmt.Sprintf("%d", time.Now().AddDate(0, 1, 0).Unix()),
		TotalCount:  "1",
		PaymentType: "",
		Amount:      "",
	})
	require.NoError(t, err)
	require.NotNil(t, payment)
	assert.Equal(t, "done", payment.Data.Message)
	assert.True(t, payment.Data.Status)

	_, err = time.Parse("2006-01-2 15:04:05.000000", payment.Data.Date)
	require.NoError(t, err)
}
