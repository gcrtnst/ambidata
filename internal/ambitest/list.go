package main

var TestList = []TestEntry{
	{"TestManagerGetChannelList", TestManagerGetChannelList},
	{"TestManagerGetDeviceChannel", TestManagerGetDeviceChannel},
	{"TestManagerGetDeviceChannelLv1", TestManagerGetDeviceChannelLv1},
	{"TestManagerDeleteData", TestManagerDeleteData},
	{"TestSenderSend", TestSenderSend},
	{"TestSenderSendBulk", TestSenderSendBulk},
	{"TestSenderSetCmnt", TestSenderSetCmnt},
	{"TestSenderSetHide", TestSenderSetHide},
	{"TestFetcherGetChannel", TestFetcherGetChannel},
}
