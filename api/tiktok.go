package api

type TikTokJSONResponse struct {
	ItemInfo TikTokItemInfo
}

type TikTokHTMLResponse struct {
	Props TikTokHTMLProps
}

type TikTokHTMLProps struct {
	PageProps TikTokResponse
}

type TikTokResponse struct {
	ServerCode int
	StatusCode int
	ItemInfo   TikTokItemInfo
}

type TikTokItemInfo struct {
	ItemStruct TikTokItemStruct
}

type TikTokItemStruct struct {
	Desc           string
	Author         TikTokAuthor
	Music          TikTokMusic
	Video          TikTokVideo
	StickersOnItem []TikTokSticker
}

type TikTokAuthor struct {
	UniqueId string
	Nickname string
}

type TikTokMusic struct {
	Id         string
	Title      string
	AuthorName string
}

type TikTokVideo struct {
	Id           string
	DownloadAddr string
	ShareCover   []string
	Bitrate      int
	CodecType    string
}

type TikTokSticker struct {
	StickerText []string
	StickerType int
}
