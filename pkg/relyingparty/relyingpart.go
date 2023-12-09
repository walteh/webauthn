package relyingparty

const (
	AppleJwtHeader          = "x-nugg-apple-token"
	AppsyncAuthHeaderPrefix = "Nugg06 "
)

type Provider interface {
	RPDisplayName() string
	RPID() string
	RPOrigin() string
}

type simpleRelyingParty struct {
	rpDisplayName string
	rpID          string
	rpOrigin      string
}

func NewSimpleRelyingParty(rpDisplayName, rpID, rpOrigin string) Provider {
	return &simpleRelyingParty{
		rpDisplayName: rpDisplayName,
		rpID:          rpID,
		rpOrigin:      rpOrigin,
	}
}

func (me *simpleRelyingParty) RPDisplayName() string { return me.rpDisplayName }
func (me *simpleRelyingParty) RPID() string          { return me.rpID }
func (me *simpleRelyingParty) RPOrigin() string      { return me.rpOrigin }
