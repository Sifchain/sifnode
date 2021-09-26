package types

// ibctransfer module event types

const (
	EventTypeConvertTransfer  = "converted_transfer"
	EventTypeConvertReceived  = "converted_received_packet"
	EventTypeConvertRefund    = "converted_refund"
	AttributeKeySentAmount    = "sent_amount"
	AttributeKeySentDenom     = "sent_denom"
	AttributeKeyPacketAmount  = "packet_amount"
	AttributeKeyPacketDenom   = "packet_denom"
	AttributeKeyConvertAmount = "converted_amount"
	AttributeKeyConvertDenom  = "converted_denom"
)
