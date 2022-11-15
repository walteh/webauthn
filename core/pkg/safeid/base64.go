package safeid

import "encoding/base64"

// Decode a base64 encoded string into a SafeID
func ParseFromChallengeString(encoded string) (id SafeID, err error) {
	return ParseFromChallenge([]byte(encoded))
}

// Decode a base64 encoded string into a SafeID
func ParseFromChallenge(encoded []byte) (id SafeID, err error) {
	dst := make([]byte, EncodedSize)
	_, err = base64.RawURLEncoding.Decode(dst, encoded)
	if err != nil {
		return SafeID{}, err
	}
	return ParseBytes(dst)
}
