package b_rawEvent

import (
	"errors"
	"fmt"
	"github.com/Sapomie/wayne-data/internal/model"
	"github.com/Sapomie/wayne-data/internal/model/cons"
	"strings"
)

const (
	//Category = "@c"
	Project     = "@p"
	Stuff       = "@s"
	Tag         = "@t"
	NoneProject = "、none"
)

type propertyPair struct {
	Key   string
	Value string
}

//从comment中获取自定义property的信息，并且添加到event的property id中
func (svc RawEventService) processCommentProperty(raw *RawEvent, taskId int) (stuffIds, tagIds, remark string, projectId int, updateInfos []string, err error) {

	pairs, remark, err := unpackEventComment(raw.Comment)
	if err != nil {
		info := fmt.Sprintf("unpack tag info error coment: %v,date: %v ", raw.Comment, raw.StartTime)
		err = errors.New(info)
		return
	}

	for _, pair := range pairs {
		switch pair.Key {
		case Stuff:
			stuff, addingInfo, err := model.InsertAndGetStuff(svc.db, &model.Stuff{Name: pair.Value})
			if err != nil {
				return "", "", "", 0, nil, err
			}
			if addingInfo != "" {
				updateInfos = append(updateInfos, addingInfo)
			}
			if stuffIds == "" {
				stuffIds = fmt.Sprint(stuff.Id)
			} else {
				stuffIds = stuffIds + "," + fmt.Sprint(stuff.Id)
			}
		case Tag:
			tag, addingInfo, err := model.InsertAndGetTag(svc.db, pair.Value)
			if err != nil {
				return "", "", "", 0, nil, err
			}
			if addingInfo != "" {
				updateInfos = append(updateInfos, addingInfo)
			}
			if tagIds == "" {
				tagIds = fmt.Sprint(tag.Id)
			} else {
				tagIds = tagIds + "," + fmt.Sprint(tag.Id)
			}
		}
	}

	//add stuff mv
	if raw.TaskName == cons.Movie {
		stuff, addingInfo, err := model.InsertAndGetStuff(svc.db, &model.Stuff{Name: cons.StuMovie})
		if err != nil {
			return "", "", "", 0, nil, err
		}
		if addingInfo != "" {
			updateInfos = append(updateInfos, addingInfo)
		}
		if stuffIds == "" {
			stuffIds = fmt.Sprint(stuff.Id)
		} else {
			stuffIds = stuffIds + "," + fmt.Sprint(stuff.Id)
		}
	}

	// project
	if cons.IsProjectTask(raw.TaskName) && !isNoneProject(raw.Comment) {
		strs := strings.Split(raw.Comment, "，")
		project, addingInfo, err := model.InsertAndGetProject(svc.db, &model.Project{Name: strs[0], TaskId: taskId})
		if err != nil {
			return "", "", "", 0, nil, err
		}
		if addingInfo != "" {
			updateInfos = append(updateInfos, addingInfo)
		}
		projectId = project.Id
	}

	return
}

func unpackEventComment(comment string) ([]*propertyPair, string, error) {

	type indexPair struct {
		start int
		end   int
	}
	indexPairs := make([]*indexPair, 0)
	indexStart := 0

OUT:
	for i1, s1 := range comment[indexStart:] {
		if string(s1) == "@" {
			indexPair := new(indexPair)
			indexPair.start = i1 + indexStart
			for i, s2 := range comment[indexPair.start:] {
				if string(s2) == "，" {
					indexPair.end = indexPair.start + i
					indexStart = indexPair.end
					indexPairs = append(indexPairs, indexPair)
					goto OUT
				}
			}
		}
	}

	var propertyString []string //由"@"开头，"，"作为结尾
	for _, pair := range indexPairs {
		str := comment[pair.start:pair.end]
		propertyString = append(propertyString, str)
	}

	var remark string
	if len(indexPairs) >= 1 {
		remarkIndex := indexPairs[len(indexPairs)-1].end + 3
		remark = comment[remarkIndex:]
	} else {
		remark = comment
	}

	tagPairs := make([]*propertyPair, 0)
	for _, str := range propertyString {
		ss := strings.Split(str, "：")
		if len(ss) < 2 {
			return nil, "", errors.New("wrong length of property")
		}
		tagPair := &propertyPair{
			Key:   ss[0],
			Value: ss[1],
		}
		tagPairs = append(tagPairs, tagPair)
	}

	return tagPairs, remark, nil
}

func isNoneProject(comment string) bool {
	if strings.Contains(comment, NoneProject) {
		return true
	}
	return false
}
