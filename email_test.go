package sendios

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func TestSendEmail(t *testing.T) {
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

	result, err := c.SendEmail(EmailRequest{
		TypeID:   "11",
		Category: CategorySystem,
		Data: map[string]interface{}{
			"redirect_url": "https://localhost:3000",
			"user": map[string]interface{}{
				"email": "igor.chornij@gen.tech",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "SUCCESS", result.Meta.Status)
}
