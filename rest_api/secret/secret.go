package secret

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")

func GetSigningKey() []byte {
	return mySigningKey
}
