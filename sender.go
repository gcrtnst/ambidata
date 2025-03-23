package ambidata

type Sender struct {
	Ch       string
	WriteKey string
	Config   *Config
}

func NewSender(ch string, writeKey string) *Sender {
	return &Sender{Ch: ch, WriteKey: writeKey}
}

func NewSenderFromChannelAccess(ca *ChannelAccess) *Sender {
	return NewSender(ca.Ch, ca.WriteKey)
}

func NewSenderFromChannelAccessLv1(ca1 *ChannelAccessLv1) *Sender {
	return NewSender(ca1.Ch, ca1.WriteKey)
}
