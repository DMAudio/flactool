package task

import (
	"fmt"
	"gitlab.com/KTGWKenta/DADP.FlacTool/config"
	"gitlab.com/MGEs/Com.Base/types"
	"strconv"
	"strings"
)

type Item struct {
	Operation string
	Arguments interface{}
}

type Collection struct {
	Owner string
	When  string
	List  []*Item
}

type List struct {
	Collections []*Collection
}

func parseTaskItem(TaskItemRaw interface{}) (*Item, interface{}) {
	if TaskItemParsed, ok := TaskItemRaw.(map[interface{}]interface{}); !ok {
		return nil, fmt.Errorf("格式不正确")
	} else if TaskItemOperationRaw, ok := TaskItemParsed["OPER"]; !ok {
		return nil, fmt.Errorf("找不到主键:OPER")
	} else if TaskItemOperationParsed, ok := TaskItemOperationRaw.(string); !ok {
		return nil, fmt.Errorf("主键 OPER 应为String型")
	} else if TaskItemArgumentsRaw, ok := TaskItemParsed["ARGS"]; !ok {
		return &Item{
			Operation: TaskItemOperationParsed,
			Arguments: nil,
		}, nil
	} else {
		return &Item{
			Operation: TaskItemOperationParsed,
			Arguments: TaskItemArgumentsRaw,
		}, nil
	}

}

func parseTaskCollection(TaskCollectionRaw interface{}) (*Collection, interface{}) {
	if CollectionParsed, ok := TaskCollectionRaw.(map[interface{}]interface{}); !ok {
		return nil, fmt.Errorf("格式不正确")
	} else if CollectionOwnerRaw, ok := CollectionParsed["OWNER"]; !ok {
		return nil, fmt.Errorf("找不到主键:OWNER")
	} else if CollectionOwnerParsed, ok := CollectionOwnerRaw.(string); !ok {
		return nil, fmt.Errorf("主键 OWNER 应为String型")
	} else if CollectionListRaw, ok := CollectionParsed["LIST"]; !ok {
		return nil, fmt.Errorf("找不到主键:LIST")
	} else if CollectionListParsed, ok := CollectionListRaw.([]interface{}); !ok {
		return nil, fmt.Errorf("主键 LIST 格式不正确")
	} else {
		whenStr := ""
		if CollectionWhenRaw, ok := CollectionParsed["WHEN"]; ok {
			if CollectionWhenParsed, ok := CollectionWhenRaw.(string); !ok {
				return nil, fmt.Errorf("主键 WHEN 应为String型")
			} else {
				whenStr = strings.TrimSpace(CollectionWhenParsed)
			}
		}

		collection := &Collection{
			Owner: CollectionOwnerParsed,
			When:  whenStr,
			List:  []*Item{},
		}

		for TaskIndex, TaskItemRaw := range CollectionListParsed {
			if TaskParsed, err := parseTaskItem(TaskItemRaw); err != nil {
				return nil, types.NewException(TMConfig_UnableToParse_TaskItem, map[string]string{
					"task": strconv.Itoa(TaskIndex + 1),
				}, err)
			} else {
				collection.List = append(collection.List, TaskParsed)
			}
		}

		return collection, nil
	}

}

func LoadTaskList(filePath string) (*List, *types.Exception) {
	taskList := List{
		Collections: make([]*Collection, 0),
	}

	if fileParsed, err := config.ParseConfig(filePath); err != nil {
		return nil, types.NewException(TMConfig_UnableToRead_TaskFile, nil, err)

	} else if TaskCollectionsRaw, ok := fileParsed["tasks"]; !ok || TaskCollectionsRaw == nil {
		return nil, types.NewException(TMConfig_UnableToParse_TaskFile, map[string]string{
			"reason": "找不到主键 TASKS 或 主键 TASKS 为空",
		}, nil)

	} else if TaskCollections, ok := TaskCollectionsRaw.([]interface{}); !ok || TaskCollections == nil {
		return nil, types.NewException(TMConfig_UnableToParse_TaskFile, map[string]string{
			"reason": "任务集列表格式不正确",
		}, nil)

	} else {
		for CollectionIndex, CollectionRaw := range TaskCollections {
			if CollectionParsed, err := parseTaskCollection(CollectionRaw); err != nil {
				return nil, types.NewException(TMConfig_UnableToParse_TaskCollection, map[string]string{
					"collection": strconv.Itoa(CollectionIndex + 1),
				}, err)
			} else {
				taskList.Collections = append(taskList.Collections, CollectionParsed)
			}
		}
	}

	return &taskList, nil
}

func (t *List) ExecuteTasks(env map[string]interface{}) *types.Exception {
	UnhandledErr := map[string]*types.Exception{}
	for cIndex, collection := range t.Collections {
		UnhandledErrInCollection := map[string]*types.Exception{}

		if collection.When != "" {
			if whenStrParsed, _, err := GlobalArgFiller().FillArgs(collection.When, nil); err != nil {
				UnhandledErrInCollection["_when_"] = types.NewException(TMTask_FailedTo_Parse_WhenPattern, nil, err)
			} else if whenStrSplit := strings.SplitN(whenStrParsed, "->", 2); len(whenStrSplit) != 2 {
				UnhandledErrInCollection["_when_"] = types.Exception_Mismatched_Format(
					"LeftPattern->RightPattern",
					whenStrParsed,
				)
			} else if strings.TrimSpace(whenStrSplit[0]) != strings.TrimSpace(whenStrSplit[1]) {
				continue
			}
		}

		for tIndex, task := range collection.List {
			if err := GlobalHandler().Execute(collection.Owner, task.Operation, env, task.Arguments); err != nil {
				UnhandledErrInCollection["T"+strconv.Itoa(tIndex)] = err
			}
		}
		if len(UnhandledErrInCollection) > 0 {
			UnhandledErr["C"+strconv.Itoa(cIndex)] = types.NewException(TMTask_Collection_UnhandledThrowable, map[string]string{
				"collection": strconv.Itoa(cIndex),
				"handler":    collection.Owner,
			}, UnhandledErrInCollection)
		}
	}
	if len(UnhandledErr) > 0 {
		return types.NewException(TMTask_UnhandledThrowable, nil, UnhandledErr)
	}

	return nil
}
