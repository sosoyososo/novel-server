package NovelSpider

import "testing"

func TestSpliteCateLog(t *testing.T) {
	var list []CatelogInfo
	err := defaultDB.Model(CatelogInfo{}).Scan(&list).Error
	if nil != err {
		t.Error(err)
	}
	for _, v := range list {
		if !defaultDB.HasTable(cateTableNameWithNovelID(v.NovelID)) {
			err := defaultDB.Table(cateTableNameWithNovelID(v.NovelID)).CreateTable(&CatelogInfo{}).Error
			if nil != err {
				t.Error(err)
				break
			}
		}
		err := cateTableWithNovelID(v.NovelID).Create(&v).Error
		if nil != err {
			t.Error(err)
			break
		}
	}
}

func TestDropTable(t *testing.T) {
	names := []string{
		"catelog_infos_00",
		"catelog_infos_5a", "catelog_infos_94", "catelog_infos_d9",
		"catelog_infos_18", "catelog_infos_5c", "catelog_infos_97",
		"catelog_infos_dd", "catelog_infos_29", "catelog_infos_5d",
		"catelog_infos_e4", "catelog_infos_2f",
		"catelog_infos_61", "catelog_infos_a1", "catelog_infos_f0",
		"catelog_infos_3c", "catelog_infos_6a", "catelog_infos_a2",
		"catelog_infos_f9", "catelog_infos_45", "catelog_infos_73",
		"catelog_infos_a8", "catelog_infos_16",
		"catelog_infos_46", "catelog_infos_7d", "catelog_infos_c2",
		"catelog_infos_47", "catelog_infos_81", "catelog_infos_c5",
		"catelog_infos_57", "catelog_infos_8b", "catelog_infos_d4",
	}
	for _, v := range names {
		if defaultDB.HasTable(v) {
			defaultDB.DropTable(v)
		}
	}
}
