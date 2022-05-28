package cq

func GetImageCQCode(imageId string) string {
	return "[CQ:image,file=" + imageId + "]"
}
