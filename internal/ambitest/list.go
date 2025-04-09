package main

var TestList = []TestEntry{
	{"TestManagerGetChannelList", TestManagerGetChannelList},
	{"TestManagerGetDeviceChannel", TestManagerGetDeviceChannel},
	{"TestManagerGetDeviceChannelLv1", TestManagerGetDeviceChannelLv1},
	{"TestManagerDeleteData", TestManagerDeleteData},
	{"TestSenderSend", TestSenderSend},
	{"TestSenderSendBulk", TestSenderSendBulk},
	{"TestSenderSendBulkTooLarge", TestSenderSendBulkTooLarge},
	{"TestSenderSetCmnt", TestSenderSetCmnt},
	{"TestSenderSetHide", TestSenderSetHide},
	{"TestFetcherGetChannel", TestFetcherGetChannel},
	{"TestFetcherFetchRange", TestFetcherFetchRange},
	{"TestFetcherFetchPeriod", TestFetcherFetchPeriod},
}
