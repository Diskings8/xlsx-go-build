# xlsx-go-build
use this item to quickly build game_xlsx to go_file

拼装后的结构体声明如下：
//package data
// type Rap struct {
//     Id              int     	`col:"id             " client:"id      "`	 //      显示顺序1
//     BadgeId         int     	`col:"badgeId        " client:"badgeId "`	 //       徽章编号
//     RuneType        int     	`col:"runeType       " client:"runeType"`	 //  可镶嵌符文类型类型
//     SkillId         int     	`col:"skillId        " client:"skillId "`	 // 普攻，对应skill
//     RuneMax         int     	`col:"runeMax        " client:"runeMax "`	 //     符文等级上限
//     LightMax        int     	`col:"lightMax       " client:"lightMax"`	 //     徽章升阶上限
//     AddHp           int     	`col:"addHp          " client:"addHp   "`	 // 自己给自己加血修正（百分比）
//     BeAddHp         int     	`col:"beAddHp        "`	 // 别人给自己加血修正（百分比）
// }


拼装后的结构体数据如下：
// package data
// import "errors" 
// func (rap *Rap)Get(id int) (*Rap, error){
//     switch id {
//         case 1 : 
//             return &Rap{ Id:1, BadgeId:50100, RuneType:1, SkillId:1002, RuneMax:50, LightMax:40, AddHp:90, BeAddHp:50 }, nil
//         case 0 : 
//             return &Rap{ Id:0, BadgeId:50101, RuneType:2, SkillId:1001, RuneMax:50, LightMax:40, AddHp:75, BeAddHp:100 }, nil
//         case 2 : 
//             return &Rap{ Id:2, BadgeId:50102, RuneType:3, SkillId:1003, RuneMax:50, LightMax:40, AddHp:95, BeAddHp:75 }, nil
//         default: return nil, errors.New("no data")
//     }
// }

windows版本直接调用build_objs_file.bat
暂不支持Unix版本