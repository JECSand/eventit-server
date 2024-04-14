package auth

import (
	"github.com/JECSand/eventit-server/src/shared/enums"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func TestSession_GetToken(t *testing.T) {
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string   // The name of the test
		want    int      // What out instance we want our function to return.
		wantErr bool     // whether we want an error.
		input   *Session // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			215,
			false,
			&Session{ProfileId: "000000000000000000000001", Role: enums.ROOT},
		},
		{
			"missing profileId",
			0,
			true,
			&Session{Role: enums.ADMIN},
		},
		{
			"missing role",
			0,
			true,
			&Session{ProfileId: "000000000000000000000001"},
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.SetDefault("auth_jwt_secret", "random")
			viper.SetDefault("auth_jwt_expiry", "15m")
			viper.SetDefault("auth_jwt_refresh_expiry", "1h")
			got, err := tt.input.GetToken()
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want { // Asserting whether we get the correct wanted value
				t.Errorf("GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_LoadSession(t *testing.T) {
	viper.SetDefault("auth_jwt_secret", "random")
	viper.SetDefault("auth_jwt_expiry", "15m")
	viper.SetDefault("auth_jwt_refresh_expiry", "1h")
	session := &Session{ProfileId: "000000000000000000000001", Role: enums.ROOT}
	tokenString, err := session.GetToken()
	if err != nil {
		panic(err)
	}
	// Defining our test slice. Each unit test should have the following properties:
	tests := []struct {
		name    string   // The name of the test
		want    *Session // What out instance we want our function to return.
		wantErr bool     // whether we want an error.
		input   string   // The input of the test
	}{
		// Here we're declaring each unit test input and output data as defined before
		{
			"success",
			&Session{ProfileId: "000000000000000000000001", Role: enums.ROOT},
			false,
			tokenString,
		},
		{
			"empty token string",
			&Session{},
			true,
			"",
		},
		{
			"invalid token string",
			&Session{},
			true,
			"invalidTokenaaaaaaaa!!!!",
		},
	}
	// Iterating over the previous test slice
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadSession(tt.input)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) { // Asserting whether we get the correct wanted value
				t.Errorf("LoadSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
