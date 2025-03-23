package ambidata

import (
	"encoding/json"
	"time"
)

type jsonRecvChannelAccessList []jsonRecvChannelAccess

func (j jsonRecvChannelAccessList) ToChannelAccessList() []ChannelAccess {
	l := make([]ChannelAccess, len(j))
	for i := range j {
		l[i] = j[i].ToChannelAccess()
	}
	return l
}

type jsonRecvChannelAccess struct {
	jsonRecvChannelInfo
	ReadKey  string `json:"readKey"`
	WriteKey string `json:"writeKey"`
}

func (j *jsonRecvChannelAccess) ToChannelAccess() ChannelAccess {
	return ChannelAccess{
		ChannelInfo: j.jsonRecvChannelInfo.ToChannelInfo(),
		ReadKey:     j.ReadKey,
		WriteKey:    j.WriteKey,
	}
}

type jsonRecvChannelAccessLv1 struct {
	Ch       string `json:"ch"`
	WriteKey string `json:"writeKey"`
}

func (j *jsonRecvChannelAccessLv1) ToChannelAccessLv1() ChannelAccessLv1 {
	return ChannelAccessLv1(*j)
}

type jsonRecvChannelInfo struct {
	Ch         string                      `json:"ch"`
	User       string                      `json:"user"`
	Created    jsonRecvTime                `json:"created"`
	Modified   jsonRecvTime                `json:"modified"`
	LastPost   jsonRecvTime                `json:"lastpost"`
	Charts     int                         `json:"charts"`
	DataPerDay int                         `json:"dataperday"`
	DCh        bool                        `json:"d_ch"`
	ChName     string                      `json:"chName"`
	ChDesc     string                      `json:"chDesc"`
	D1         jsonFieldInfo               `json:"d1"`
	D2         jsonFieldInfo               `json:"d2"`
	D3         jsonFieldInfo               `json:"d3"`
	D4         jsonFieldInfo               `json:"d4"`
	D5         jsonFieldInfo               `json:"d5"`
	D6         jsonFieldInfo               `json:"d6"`
	D7         jsonFieldInfo               `json:"d7"`
	D8         jsonFieldInfo               `json:"d8"`
	Loc        jsonMaybe[jsonRecvLocation] `json:"loc"`
	PhotoID    string                      `json:"photoid"`
	DevKeys    []string                    `json:"devkeys"`
	Bd         string                      `json:"bd"`
	LastData   jsonRecvLastData            `json:"lastdata"`
}

func (j *jsonRecvChannelInfo) ToChannelInfo() ChannelInfo {
	return ChannelInfo{
		Ch:         j.Ch,
		User:       j.User,
		Created:    time.Time(j.Created),
		Modified:   time.Time(j.Modified),
		LastPost:   time.Time(j.LastPost),
		Charts:     j.Charts,
		DataPerDay: j.DataPerDay,
		DCh:        j.DCh,
		ChName:     j.ChName,
		ChDesc:     j.ChDesc,
		D1:         FieldInfo(j.D1),
		D2:         FieldInfo(j.D2),
		D3:         FieldInfo(j.D3),
		D4:         FieldInfo(j.D4),
		D5:         FieldInfo(j.D5),
		D6:         FieldInfo(j.D6),
		D7:         FieldInfo(j.D7),
		D8:         FieldInfo(j.D8),
		Loc:        Maybe[Location]{V: Location(j.Loc.V), OK: j.Loc.OK},
		PhotoID:    j.PhotoID,
		DevKeys:    j.DevKeys,
		Bd:         j.Bd,
		LastData:   j.LastData.ToLastData(),
	}
}

type jsonFieldInfo struct {
	Name  string `json:"name,omitzero"`
	Color Color  `json:"color,omitzero"`
}

type jsonRecvLastData struct {
	jsonRecvData
	ID string `json:"_id"`
}

func (j *jsonRecvLastData) ToLastData() LastData {
	return LastData{
		Data: j.jsonRecvData.ToData(),
		ID:   j.ID,
	}
}

type jsonRecvDataList []jsonRecvData

func (j jsonRecvDataList) ToDataList() []Data {
	l := make([]Data, len(j))
	for i := range j {
		l[i] = j[i].ToData()
	}
	return l
}

type jsonRecvData struct {
	Created jsonRecvTime                `json:"created"`
	D1      jsonMaybe[float64]          `json:"d1"`
	D2      jsonMaybe[float64]          `json:"d2"`
	D3      jsonMaybe[float64]          `json:"d3"`
	D4      jsonMaybe[float64]          `json:"d4"`
	D5      jsonMaybe[float64]          `json:"d5"`
	D6      jsonMaybe[float64]          `json:"d6"`
	D7      jsonMaybe[float64]          `json:"d7"`
	D8      jsonMaybe[float64]          `json:"d8"`
	Loc     jsonMaybe[jsonRecvLocation] `json:"loc"`
	Cmnt    string                      `json:"cmnt"`
	Hide    bool                        `json:"hide"`
}

func (j *jsonRecvData) ToData() Data {
	return Data{
		Created: time.Time(j.Created),
		D1:      Maybe[float64](j.D1),
		D2:      Maybe[float64](j.D2),
		D3:      Maybe[float64](j.D3),
		D4:      Maybe[float64](j.D4),
		D5:      Maybe[float64](j.D5),
		D6:      Maybe[float64](j.D6),
		D7:      Maybe[float64](j.D7),
		D8:      Maybe[float64](j.D8),
		Loc:     Maybe[Location]{V: Location(j.Loc.V), OK: j.Loc.OK},
		Cmnt:    j.Cmnt,
		Hide:    j.Hide,
	}
}

type jsonSendDataListRequest struct {
	WriteKey string           `json:"writeKey"`
	Data     jsonSendDataList `json:"data"`
}

type jsonSendDataList []jsonSendData

func toJSONSendDataList(arr []Data) jsonSendDataList {
	l := make(jsonSendDataList, len(arr))
	for i := range arr {
		l[i] = toJSONSendData(arr[i])
	}
	return l
}

type jsonSendDataRequest struct {
	jsonSendData
	WriteKey string `json:"writeKey"`
}

type jsonSendData struct {
	Created time.Time          `json:"created,omitzero"`
	D1      jsonMaybe[float64] `json:"d1,omitzero"`
	D2      jsonMaybe[float64] `json:"d2,omitzero"`
	D3      jsonMaybe[float64] `json:"d3,omitzero"`
	D4      jsonMaybe[float64] `json:"d4,omitzero"`
	D5      jsonMaybe[float64] `json:"d5,omitzero"`
	D6      jsonMaybe[float64] `json:"d6,omitzero"`
	D7      jsonMaybe[float64] `json:"d7,omitzero"`
	D8      jsonMaybe[float64] `json:"d8,omitzero"`
	Lat     jsonMaybe[float64] `json:"lat,omitzero"`
	Lng     jsonMaybe[float64] `json:"lng,omitzero"`
	Cmnt    string             `json:"cmnt,omitzero"`
}

func toJSONSendData(data Data) jsonSendData {
	return jsonSendData{
		Created: data.Created,
		D1:      jsonMaybe[float64](data.D1),
		D2:      jsonMaybe[float64](data.D2),
		D3:      jsonMaybe[float64](data.D3),
		D4:      jsonMaybe[float64](data.D4),
		D5:      jsonMaybe[float64](data.D5),
		D6:      jsonMaybe[float64](data.D6),
		D7:      jsonMaybe[float64](data.D7),
		D8:      jsonMaybe[float64](data.D8),
		Lat:     jsonMaybe[float64]{V: data.Loc.V.Lat, OK: data.Loc.OK},
		Lng:     jsonMaybe[float64]{V: data.Loc.V.Lng, OK: data.Loc.OK},
		Cmnt:    data.Cmnt,
		// ignore hide
	}
}

type jsonRecvTime time.Time

func (j *jsonRecvTime) UnmarshalJSON(data []byte) error {
	err := (*time.Time)(j).UnmarshalJSON(data)
	if err == nil && (*time.Time)(j).Equal(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)) {
		*j = jsonRecvTime{}
	}
	return err
}

type jsonRecvLocation Location

func (j *jsonRecvLocation) UnmarshalJSON(data []byte) error {
	var loc [2]float64
	err := json.Unmarshal(data, &loc)
	j.Lat = loc[1]
	j.Lng = loc[0]
	return err
}

type jsonMaybe[T any] Maybe[T]

func (j jsonMaybe[T]) IsZero() bool {
	return !j.OK
}

func (j jsonMaybe[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.V)
}

func (j *jsonMaybe[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &j.V)
	j.OK = err == nil
	return err
}
