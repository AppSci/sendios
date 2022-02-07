package sendios

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMappingError(t *testing.T) {
	data := `{"_meta":{"status":"ERROR","time":4,"count":1},"data":{"error":"Not found"}}`

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal([]byte(data), &resp))

	assert.Equal(t, "ERROR", resp.Meta.Status)
	assert.Equal(t, 4, resp.Meta.Time)
	assert.Equal(t, 1, resp.Meta.Count)

	assert.Equal(t, "Not found", resp.Data.Error)
}

func TestLoadingFromEnv(t *testing.T) {
	require.NoError(t, os.Setenv("SENDIOS_CONFIG", "client_id=134983&client_token=xxx&project=21515"))

	c, err := NewFromEnv()
	require.NoError(t, err)
	require.NotNil(t, c)
	assert.Equal(t, 134983, c.Config.ClientID)
	assert.Equal(t, "xxx", c.Config.ClientToken)
	assert.Equal(t, 21515, c.Config.Project)
}

func TestMapLastPaymentResponse(t *testing.T) {
	data := `{"_meta":{"status":"SUCCESS","time":17,"count":3},"data":{"message":"done","date":"2021-07-01 12:34:58.000000","status":true}}`

	var resp LastPaymentResponse
	require.NoError(t, json.Unmarshal([]byte(data), &resp))

	assert.Equal(t, "SUCCESS", resp.Meta.Status)
	assert.Equal(t, 17, resp.Meta.Time)
	assert.Equal(t, 3, resp.Meta.Count)

	assert.Equal(t, "2021-07-01 12:34:58.000000", resp.Data.Date)
	assert.Equal(t, "done", resp.Data.Message)
	assert.True(t, resp.Data.Status)
}
