package novel

import (
	"../../DO"
	"../../NovelSpider"
	"../../service"
)

func init() {
	func() {
		var d DO.PageInfoDO
		service.RegisterAuthNoNeedPath("/novel/list")
		service.RegisterListJSONServiceV2("/novel/list", &d,
			func(ctx service.ServiceCtx) (interface{}, int, error) {
				return NovelSpider.ListSummary(d.Page(), d.Size())
			}, "小说列表")
	}()

	func() {
		var d DO.DetailDO
		service.RegisterAuthNoNeedPath("/novel/search")
		service.RegisterJSONServiceV2(
			"/novel/search/:key", &d,
			func(ctx service.ServiceCtx) (interface{}, error) {
				key := ctx.Ctx.Param("key")
				return NovelSpider.SearchSummary(key)
			}, "小说搜索")
	}()

	func() {
		var d DO.DetailPageDO
		service.RegisterAuthNoNeedPath("/novel/chapters")
		service.RegisterListJSONServiceV2("/novel/chapters/:id", &d,
			func(ctx service.ServiceCtx) (interface{}, int, error) {
				id := ctx.Ctx.Param("id")
				return NovelSpider.CatelogListOfNovel(id, d.Page(), d.Size())
			}, "小说目录")
	}()

	func() {
		var d DO.EmptyDO
		service.RegisterAuthNoNeedPath("/chapter/detail")
		service.RegisterJSONServiceV2("/chapter/detail/:novelId/:id", &d,
			func(ctx service.ServiceCtx) (interface{}, error) {
				id := ctx.Ctx.Param("id")
				novelId := ctx.Ctx.Param("novelId")
				return NovelSpider.ChapterDetail(novelId, id)
			}, "小说章节详情")
	}()
}
