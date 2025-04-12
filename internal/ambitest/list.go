package main

var TestList = []TestEntry{
	{"TestManagerGetChannelList", TestManagerGetChannelList},
	{"TestManagerGetDeviceChannel", TestManagerGetDeviceChannel},
	{"TestManagerGetDeviceChannelLv1", TestManagerGetDeviceChannelLv1},
	{"TestManagerDeleteData", TestManagerDeleteData},
	{"TestSenderSend", TestSenderSend},
	{"TestSenderSendTimePrecision", TestSenderSendTimePrecision},
	{"TestSenderSendCmntSize", TestSenderSendCmntSize},
	{"TestSenderSendBulk", TestSenderSendBulk},
	{"TestSenderSendBulkTooLarge", TestSenderSendBulkTooLarge},
	{"TestSenderSetCmnt", TestSenderSetCmnt},
	{"TestSenderSetHide", TestSenderSetHide},
	{"TestSenderSetHideMultiple", TestSenderSetHideMultiple},
	{"TestFetcherGetChannel", TestFetcherGetChannel},
	{"TestFetcherFetchRange", TestFetcherFetchRange},
	{"TestFetcherFetchPeriod", TestFetcherFetchPeriod},
}
