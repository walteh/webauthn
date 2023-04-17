package types

import (
	"git.nugg.xyz/webauthn/pkg/hex"
)

type VerifyAuenticatorDataArgs struct {
	Data                           hex.Hash
	AppId                          string
	RelyingPartyID                 string
	RequireUserVerification        bool
	RequireUserPresence            bool
	LastSignCount                  uint64
	OptionalAttestedCredentialData AttestedCredentialData
	UseSavedAttestedCredentialData bool
}

// Authenticators respond to Relying Party requests by returning an object derived from the
// AuthenticatorResponse interface. See §5.2. Authenticator Responses
// https://www.w3.org/TR/webauthn/#iface-authenticatorresponse
type AuthenticatorResponse struct {
	// From the spec https://www.w3.org/TR/webauthn/#dom-authenticatorresponse-clientdatajson
	// This attribute contains a JSON serialization of the client data passed to the authenticator
	// by the client in its call to either create() or get().
	UTF8ClientDataJSON string `json:"clientDataJSON"`
}

// AuthenticatorData From §6.1 of the spec.
// The authenticator data structure encodes contextual bindings made by the authenticator. These bindings
// are controlled by the authenticator itself, and derive their trust from the WebAuthn Relying Party's
// assessment of the security properties of the authenticator. In one extreme case, the authenticator
// may be embedded in the client, and its bindings may be no more trustworthy than the client data.
// At the other extreme, the authenticator may be a discrete entity with high-security hardware and
// software, connected to the client over a secure channel. In both cases, the Relying Party receives
// the authenticator data in the same format, and uses its knowledge of the authenticator to make
// trust decisions.
//
// The authenticator data, at least during attestation, contains the Public Key that the RP stores
// and will associate with the user attempting to register.
type AuthenticatorData struct {
	RPIDHash hex.Hash               `json:"rpid"`
	Flags    AuthenticatorFlags     `json:"flags"`
	Counter  uint64                 `json:"sign_count"`
	AttData  AttestedCredentialData `json:"att_data"`
	ExtData  hex.Hash               `json:"ext_data"`
}

type AttestedCredentialData struct {
	AAGUID       hex.Hash `json:"aaguid"`
	CredentialID hex.Hash `json:"credential_id"`
	// The raw credential public key bytes received from the attestation data
	CredentialPublicKey hex.Hash `json:"public_key"`
}

// AuthenticatorAttachment https://www.w3.org/TR/webauthn/#dom-authenticatorselectioncriteria-authenticatorattachment
type AuthenticatorAttachment string

const (
	// Platform - A platform authenticator is attached using a client device-specific transport, called
	// platform attachment, and is usually not removable from the client device. A public key credential
	//  bound to a platform authenticator is called a platform credential.
	Platform AuthenticatorAttachment = "platform"
	// CrossPlatform A roaming authenticator is attached using cross-platform transports, called
	// cross-platform attachment. Authenticators of this class are removable from, and can "roam"
	// among, client devices. A public key credential bound to a roaming authenticator is called a
	// roaming credential.
	CrossPlatform AuthenticatorAttachment = "cross-platform"
)

// ResidentKeyRequirement https://www.w3.org/TR/webauthn/#dom-authenticatorselectioncriteria-residentkey
type ResidentKeyRequirement string

const (
	// ResidentKeyRequirementDiscouraged indicates to the client we do not want a discoverable credential. This is the default.
	ResidentKeyRequirementDiscouraged ResidentKeyRequirement = "discouraged"

	// ResidentKeyRequirementPreferred indicates to the client we would prefer a discoverable credential.
	ResidentKeyRequirementPreferred ResidentKeyRequirement = "preferred"

	// ResidentKeyRequirementRequired indicates to the client we require a discoverable credential and that it should
	// fail if the credential does not support this feature.
	ResidentKeyRequirementRequired ResidentKeyRequirement = "required"
)

// Authenticators may implement various transports for communicating with clients. This enumeration defines
// hints as to how clients might communicate with a particular authenticator in order to obtain an assertion
// for a specific credential. Note that these hints represent the WebAuthn Relying Party's best belief as to
// how an authenticator may be reached. A Relying Party may obtain a list of transports hints from some
// attestation statement formats or via some out-of-band mechanism; it is outside the scope of this
// specification to define that mechanism.
// See §5.10.4. Authenticator Transport https://www.w3.org/TR/webauthn/#transport
type AuthenticatorTransport string

const (
	// USB The authenticator should transport information over USB
	USB AuthenticatorTransport = "usb"
	// NFC The authenticator should transport information over Near Field Communication Protocol
	NFC AuthenticatorTransport = "nfc"
	// BLE The authenticator should transport information over Bluetooth
	BLE AuthenticatorTransport = "ble"
	// Internal the client should use an internal source like a TPM or SE
	Internal AuthenticatorTransport = "internal"
)

// A WebAuthn Relying Party may require user verification for some of its operations but not for others,
// and may use this type to express its needs.
// See §5.10.6. User Verification Requirement Enumeration https://www.w3.org/TR/webauthn/#userVerificationRequirement
type UserVerificationRequirement string

const (
	// VerificationRequired User verification is required to create/release a credential
	VerificationRequired UserVerificationRequirement = "required"
	// VerificationPreferred User verification is preferred to create/release a credential
	VerificationPreferred UserVerificationRequirement = "preferred" // This is the default
	// VerificationDiscouraged The authenticator should not verify the user for the credential
	VerificationDiscouraged UserVerificationRequirement = "discouraged"
)

// AuthenticatorFlags A byte of information returned during during ceremonies in the
// authenticatorData that contains bits that give us information about the
// whether the user was present and/or verified during authentication, and whether
// there is attestation or extension data present. Bit 0 is the least significant bit.
type AuthenticatorFlags byte

// The bits that do not have flags are reserved for future use.
const (
	// FlagUserPresent Bit 00000001 in the byte sequence. Tells us if user is present
	FlagUserPresent AuthenticatorFlags = 1 << iota // Referred to as UP
	_                                              // Reserved
	// FlagUserVerified Bit 00000100 in the byte sequence. Tells us if user is verified
	// by the authenticator using a biometric or PIN
	FlagUserVerified // Referred to as UV
	_                // Reserved
	_                // Reserved
	_                // Reserved
	// FlagAttestedCredentialData Bit 01000000 in the byte sequence. Indicates whether
	// the authenticator added attested credential data.
	FlagAttestedCredentialData // Referred to as AT
	// FlagHasExtension Bit 10000000 in the byte sequence. Indicates if the authenticator data has extensions.
	FlagHasExtensions //  Referred to as ED
)

// UserPresent returns if the UP flag was set
func (flag AuthenticatorFlags) UserPresent() bool {
	return (flag & FlagUserPresent) == FlagUserPresent
}

// UserVerified returns if the UV flag was set
func (flag AuthenticatorFlags) UserVerified() bool {
	return (flag & FlagUserVerified) == FlagUserVerified
}

// HasAttestedCredentialData returns if the AT flag was set
func (flag AuthenticatorFlags) HasAttestedCredentialData() bool {
	return (flag & FlagAttestedCredentialData) == FlagAttestedCredentialData
}

// HasExtensions returns if the ED flag was set
func (flag AuthenticatorFlags) HasExtensions() bool {
	return (flag & FlagHasExtensions) == FlagHasExtensions
}
