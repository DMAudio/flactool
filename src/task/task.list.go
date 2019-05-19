package task

import (
	"flag"
	"fmt"
	"p20190417/config"
	"p20190417/types"
	"strconv"
)

type Item struct {
	Operation string
	Arguments interface{}
}

type Collection struct {
	Owner string
	List  []*Item
}

var CollectionFile *string

func Init() *types.Exception {
	config.NewRegister(func() {
		CollectionFile = flag.String("task", "", "任务配置文件")
	})

	if err := GlobalHandler().Register("COMMON", Handler_COMMON); err != nil {
		return err
	}

	if err := GlobalArgFilter().Register("env", Filter_Env); err != nil {
		return err
	}

	if err := GlobalArgFilter().Register("fmtFName", Filter_FmtFileName); err != nil {
		return err
	}

	return nil
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
		collection := &Collection{
			Owner: CollectionOwnerParsed,
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

func LoadTaskList() ([]*Collection, *types.Exception) {
	Collections := make([]*Collection, 0)
	if fileParsed, err := config.ParseConfig(*CollectionFile); err != nil {
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
				Collections = append(Collections, CollectionParsed)
			}
		}
	}

	return Collections, nil
}

func ExecuteTasks() *types.Exception {
	if Collections, err := LoadTaskList(); err != nil {
		return err
	} else {
		UnhandledErr := map[string]*types.Exception{}
		for cIndex, collection := range Collections {
			UnhandledErrInCollection := map[string]*types.Exception{}
			for tIndex, task := range collection.List {
				if err := GlobalHandler().Execute(collection.Owner, task.Operation, task.Arguments); err != nil {
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
	}
	return nil
}
