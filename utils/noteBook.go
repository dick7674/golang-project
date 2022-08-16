package utils

import (
	"strings"
	"time"
)

var notes = make(map[string]Note)

var curNote = make(map[string]string)

func GetCurNote(userId string, userName string) string {
	cur, res := curNote[userId]
	if res {
		return "@" + userName + " \n当已打开的笔记【" + cur + "】，请输入：@机器人 /编写笔记 笔记内容"
	} else {
		return "@" + userName + " \n当前暂无打开的笔记"
	}
}

func CreateNote(title string, userId string, userName string) string {
	_, res := notes[title]
	if res {
		return "@" + userName + " \n标题为【" + title + "】的笔记已经存在！"
	} else {
		curNote[userId] = title
		return "@" + userName + " \n创建【" + title + "】笔记成功，正在编辑中，请输入：@机器人 /编写笔记 笔记内容"
	}
}

func UpdateNote(title string, userId string, userName string) string {
	_, res := notes[title]
	if res {
		curNote[userId] = title
		return "@" + userName + " \n正在编辑【" + title + "】，请输入：@机器人 /编写笔记 笔记内容"
	} else {
		return "@" + userName + " \n标题为【" + title + "】的笔记不存在！"
	}
}

func AddNote(content string, userId string, userName string) string {
	if curNote[userId] == "" {
		return "当前暂无正在编辑的笔记"
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	note := Note{content, t, userId, userName}
	notes[curNote[userId]] = note
	res := "@" + userName + " \n笔记【" + curNote[userId] + "】保存成功"
	curNote[userId] = ""
	return res
}

func GetAllNote() string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	keys := make([]string, len(notes))
	for k := range notes {
		keys[j] = k
		j++
	}
	return "已创建的笔记列表：\n" + strings.Join(keys, "\n")
}

func GetNote(title string) string {
	note, res := notes[title]
	if res {
		return "笔记标题：" + title + "\n" +
			"笔记内容：" + note.Content + "\n" +
			"笔记保存时间：" + note.Days + "\n" +
			"作者：" + note.AuthorName
	} else {
		return "标题为【" + title + "】的笔记不存在！"
	}
}

func RemoveNote(title string, userId string, userName string, isAdmin bool) string {

	note, res := notes[title]
	if res {
		if note.AuthorId == userId || isAdmin {
			delete(notes, title)
			return "@" + userName + " \n标题为【" + title + "】的笔记已删除！"
		} else {
			return "@" + userName + " \n您无法删除【" + title + "】！"
		}
	} else {
		return "@" + userName + " \n标题为【" + title + "】的笔记不存在！"
	}
}

type Note struct {
	Content    string `json:"content"`    //笔记内容
	Days       string `json:"days"`       //日期，例如2022-01-01 00:00:00
	AuthorId   string `json:"authorId"`   //笔记作者ID
	AuthorName string `json:"authorName"` //笔记作者昵称
}
