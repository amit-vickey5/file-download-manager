package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/amit/file-download-manager/internal/util"
	"github.com/amit/file-download-manager/pkg/db"
	"github.com/amit/file-download-manager/pkg/logger"
	"strings"
	"time"
)

func generateSecretKey(username string) string {
	secretKey := "FILE_DOWNLOAD_MANAGER_"+strings.ToUpper(username)+"_"
	curTime := time.Now().Unix()
	timeKey := fmt.Sprintf("%v", curTime)
	secretKey = secretKey + timeKey[len(timeKey)-5:]
	return secretKey
}

func GetUserDetails(ctx context.Context, username string) (User, error) {
	logger.LogStatement("fetching details of user ::", username)
	var users []User
	whereAttributes := make(map[string]interface{})
	whereAttributes["username"] = username
	findErr := db.RepoClient.FindWhere(ctx, &users, whereAttributes, USER_TABLE_NAME)
	if findErr != nil {
		logStatement := "ERROR :: exception while reading user details for username :: "+username+" ::"
		logger.LogStatement(logStatement, findErr)
		return User{}, errors.New("internal_server_error")
	}
	if len(users) < 1 {
		errMsg := "no user found with username :: "+username
		logger.LogStatement(errMsg, nil)
		return User{}, errors.New("no_user_found_for_username")
	} else if len(users) > 1 {
		errMsg := "multiple users found with username :: "+username
		logger.LogStatement(errMsg, nil)
		return User{}, errors.New("multiple_users_found_for_username")
	}
	return users[0], nil
}

func AddUser(ctx context.Context, usr *User) error {
	dbUsr, dbUsrErr := GetUserDetails(ctx, usr.Username)
	if dbUsrErr != nil {
		if dbUsrErr.Error() == "internal_server_error" || dbUsrErr.Error() == "multiple_users_found_for_username" {
			return dbUsrErr
		}
	}
	if dbUsr != (User{}) {
		return errors.New("username_already_exists")
	}
	uniqueId := util.GenerateUniqueId()
	secretKey := generateSecretKey(usr.Username)
	usr.Id = USER_ID_PREFIX + uniqueId
	usr.SecretKey = secretKey
	usr.CreatedAt = time.Now().Unix()
	usr.SecretKeyUpdatedAt = time.Now().Unix()
	createErr := db.RepoClient.Create(ctx, USER_TABLE_NAME, usr)
	if createErr != nil {
		return createErr
	}
	return nil
}
