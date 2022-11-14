package types

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/google/uuid"
	"github.com/mattgillard/user-pii-demo/common/crypto"
	"log"
	"os"
)

type User struct {
	Id       string `json:"id" dynamodbav:"id"`
	Name     string `json:"name" dynamodbav:"name"`
	Address  string `json:"address" dynamodbav:"address"`
	Status   bool   `json:"status" dynamodbav:"status"`
	Passport string `json:"passport" dynamodbav:"passport"`
}

type UpdateUser struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Status  bool   `json:"status" validate:"required"`
}

type CreateUser struct {
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Passport string `json:"passport" validate:"required"`
}

var sm secretsmanager.Client
var km kms.Client
var UserServiceKey string
var PassportServiceKey string
var prefix string

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	sm = *secretsmanager.NewFromConfig(sdkConfig)
	km = *kms.NewFromConfig(sdkConfig)
	UserServiceKey = os.Getenv("ServiceKey")
	PassportServiceKey = os.Getenv("PassportServiceKey")
	prefix = "keys/user/"
}

func (u *User) getKey() string {
	smresult, err := sm.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(prefix + u.Id),
	})
	if err != nil {
		return ""
		//	decryptedField = "**KEY DELETED**"
		//	log.Println(err.Error())
		//	err = nil
	}
	return *smresult.SecretString
}

func (u *User) Encode() error {
	key := uuid.NewString()
	id := uuid.NewString()
	result, err := sm.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name: aws.String(prefix + id),
		// descriptions are optional
		Description: aws.String("Encryption Key for " + id),
		// You must provide either SecretString or SecretBytes.
		// Both is considered invalid.
		SecretString: aws.String(key),
	})

	if err != nil {
		return err
	}

	log.Printf("New secret key created with ARN = %s", *result.ARN)

	u.Address, err = encryptField(u.Address, UserServiceKey, key)
	if err != nil {
		return err
	}
	u.Passport, err = encryptField(u.Passport, PassportServiceKey, key)
	if err != nil {
		return err
	}
	u.Id = id
	return nil
}

func (u *User) Decode() error {
	var err error
	userKey := u.getKey()
	log.Printf("Decoded user encryption key. Result: %s", userKey)
	u.Passport, err = decryptField(u.Passport, PassportServiceKey, userKey)
	if err != nil {
		return err
	}
	u.Address, err = decryptField(u.Address, UserServiceKey, userKey)
	if err != nil {
		return err
	}
	log.Printf("Decoded user. Result: %v", u)
	return nil
}

func encryptField(fieldName string, kmsAlias string, newKey string) (string, error) {
	ct := crypto.Encrypt([]byte(fieldName), newKey)
	keyInput := &kms.EncryptInput{
		KeyId:     &kmsAlias,
		Plaintext: ct,
	}
	keyOutput, err := km.Encrypt(context.TODO(), keyInput)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(keyOutput.CiphertextBlob), nil
}

func decryptField(fieldName string, kmsAlias string, key string) (string, error) {

	if key == "" {
		return "**KEY DELETED**", nil
	}
	decoded, err := base64.RawStdEncoding.DecodeString(fieldName)
	if err != nil {
		return "", err
	}
	DecryptInput := &kms.DecryptInput{
		KeyId:          &kmsAlias,
		CiphertextBlob: decoded,
	}
	decrypted, err := km.Decrypt(context.TODO(), DecryptInput)
	if err != nil {
		// If error from KMS - access to key was denied so just return ciphertext
		return fieldName, nil
	}

	fieldDecrypted := string(crypto.Decrypt(decrypted.Plaintext, key))
	return fieldDecrypted, nil
}
