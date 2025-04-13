package main

var TestList = []TestEntry{
	{"TestManagerGetChannelList", TestManagerGetChannelList, false},
	{"TestManagerGetDeviceChannel", TestManagerGetDeviceChannel, false},
	{"TestManagerGetDeviceChannelLv1", TestManagerGetDeviceChannelLv1, false},
	{"TestManagerDeleteData", TestManagerDeleteData, false},
	{"TestSenderSend", TestSenderSend, false},
	{"TestSenderSendTimePrecision", TestSenderSendTimePrecision, true},
	{"TestSenderSendCmntSize", TestSenderSendCmntSize, true},
	{"TestSenderSendBulk", TestSenderSendBulk, false},
	{"TestSenderSendBulkTooLarge", TestSenderSendBulkTooLarge, true},
	{"TestSenderSetCmnt", TestSenderSetCmnt, false},
	{"TestSenderSetHide", TestSenderSetHide, false},
	{"TestSenderSetHideNonexistent", TestSenderSetHideNonexistent, true},
	{"TestSenderSetHideMultiple", TestSenderSetHideMultiple, true},
	{"TestFetcherGetChannel", TestFetcherGetChannel, false},
	{"TestFetcherFetchRange", TestFetcherFetchRange, false},
	{"TestFetcherFetchPeriod", TestFetcherFetchPeriod, false},
}
