package sendios

import (
	"encoding/json"
	"fmt"
	"github.com/appsci/app-core/env"
	"github.com/appsci/app-core/pointer"
	guuid "github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
)

func TestGetUser(t *testing.T) {
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

	user, err := c.GetUserInfo(os.Getenv("TEST_EMAIL"))
	require.NoError(t, err)
	require.NotZero(t, user.Data.User.ID)
}

func TestGetUndefinedUser(t *testing.T) {
	require.NoError(t, godotenv.Load())

	require.NotEmpty(t, os.Getenv("SENDIOS_PROJECT"))
	require.NotEmpty(t, os.Getenv("SENDIOS_CLIENT_ID"))
	require.NotEmpty(t, os.Getenv("SENDIOS_TOKEN"))

	project, err := strconv.Atoi(os.Getenv("SENDIOS_PROJECT"))
	require.NoError(t, err)
	client, err := strconv.Atoi(os.Getenv("SENDIOS_CLIENT_ID"))
	require.NoError(t, err)

	c := New(project, client, os.Getenv("SENDIOS_TOKEN"))

	user, err := c.GetUserInfo("undefined@gmail.com")
	require.Error(t, err)
	require.Nil(t, user)
	require.Equal(t, fmt.Sprintf("get user error: Not found  user for project %s and email undefined@gmail.com", os.Getenv("SENDIOS_PROJECT")), err.Error())
}

func TestMapResponse(t *testing.T) {
	data := `{"_meta":{"status":"SUCCESS","time":33,"count":1},"data":{"user":{"id":1186147905,"project_id":19300,"project_title":"Boosters Project","email":"test@gmail.com","name":"Anton","gender":null,"country":"USA","language":"en","err_response":0,"last_online":null,"last_reaction":null,"last_mailed":"2021-06-18 15:26:59","last_request":null,"activation":null,"meta":{"profile":{"age":0,"ak":null,"photo":null,"partner_id":null}},"clicks":0,"sends":8,"created_at":"2021-04-08 10:04:35","sent_mails":[{"id":16072377551,"category_id":1,"type_id":5,"subject_id":74416,"template_id":16282,"split_group":0,"source_id":0,"server_id":129,"created_at":"2021-06-18 15:26:59","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail129","server_ip":"61.81.27.209","source_name":"system","type_sig":"email_confirm"},{"id":15639511862,"category_id":1,"type_id":5,"subject_id":74416,"template_id":16282,"split_group":0,"source_id":0,"server_id":128,"created_at":"2021-05-18 08:55:06","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail128","server_ip":"61.81.27.208","source_name":"system","type_sig":"email_confirm"},{"id":15470158441,"category_id":1,"type_id":7,"subject_id":72121,"template_id":17917,"split_group":0,"source_id":0,"server_id":127,"created_at":"2021-05-06 11:52:27","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail127","server_ip":"61.81.27.207","source_name":"system","type_sig":"purchase_completed"},{"id":15470073547,"category_id":1,"type_id":5,"subject_id":74416,"template_id":16282,"split_group":0,"source_id":0,"server_id":212,"created_at":"2021-05-06 11:44:09","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail212","server_ip":"204.74.252.42","source_name":"system","type_sig":"email_confirm"},{"id":15469727870,"category_id":1,"type_id":7,"subject_id":72121,"template_id":17917,"split_group":0,"source_id":0,"server_id":128,"created_at":"2021-05-06 11:11:26","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail128","server_ip":"61.81.27.208","source_name":"system","type_sig":"purchase_completed"},{"id":15469684197,"category_id":1,"type_id":5,"subject_id":74416,"template_id":16282,"split_group":0,"source_id":0,"server_id":204,"created_at":"2021-05-06 11:07:39","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail204","server_ip":"61.81.27.204","source_name":"system","type_sig":"email_confirm"},{"id":15469671046,"category_id":1,"type_id":7,"subject_id":72121,"template_id":17917,"split_group":0,"source_id":0,"server_id":127,"created_at":"2021-05-06 11:06:20","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail127","server_ip":"61.81.27.207","source_name":"system","type_sig":"purchase_completed"},{"id":15469665293,"category_id":1,"type_id":5,"subject_id":74416,"template_id":16282,"split_group":0,"source_id":0,"server_id":128,"created_at":"2021-05-06 11:05:47","mail_group_id":11,"pre_header_id":null,"category_sig":"System","server_name":"mail128","server_ip":"61.81.27.208","source_name":"system","type_sig":"email_confirm"}],"unsubscribe":[],"unsubscribe_types":[],"unsub_promo":[],"webpush":[],"last_payment":[],"channel_id":null,"subchannel_id":null}}}`

	var resp UserResponse
	require.NoError(t, json.Unmarshal([]byte(data), &resp))

	assert.Equal(t, "SUCCESS", resp.Meta.Status)
	assert.Equal(t, 33, resp.Meta.Time)
	assert.Equal(t, 1, resp.Meta.Count)

	assert.NotZero(t, resp.Data.User.ID)
}

func TestMapUser(t *testing.T) {
	data, err := json.Marshal(CreateUserRequest{VIP: pointer.Int(1)})
	require.NoError(t, err)
	assert.Equal(t, `{"vip":1}`, string(data))

	data, err = json.Marshal(CreateUserRequest{VIP: pointer.Int(0)})
	require.NoError(t, err)
	assert.Equal(t, `{"vip":0}`, string(data))
}

func TestUnsubscribe(t *testing.T) {
	require.NoError(t, env.Load(".env"))

	client, err := NewFromEnv()
	require.NoError(t, err)
	require.NotNil(t, client)

	email := fmt.Sprintf("%s-testx@mail.ua", guuid.NewString()[:4])

	_, err = client.SendEmail(EmailRequest{
		TypeID:   1,
		Category: CategorySystem,
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"email": email,
			},
		},
	})
	require.NoError(t, err)

	unsubscribeUserResponse, err := client.UnsubscribeUser(email)
	require.NoError(t, err)
	require.NotNil(t, unsubscribeUserResponse)
	require.True(t, unsubscribeUserResponse.Data.Unsub)
}
