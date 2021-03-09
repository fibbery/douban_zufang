package g

const (
	// 豆列详情页
	DouListUrl = "https://www.douban.com/doulist/%s/?start=%d&sort=time&playable=0&sub_type="

	//收藏该主题的豆列列表
	TopicCollectUrl = "https://m.douban.com/rexxar/api/v2/group/topic/%s/collections?start=0&count=100"

	// 文章主题页
	TopicUrl = "https://www.douban.com/group/topic/%s/"

	// 豆瓣中豆列固定分页记录，貌似是设置了朝25取整
	DouListPageSize = 25
)
