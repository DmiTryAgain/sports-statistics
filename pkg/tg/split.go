package tg

import (
	"strings"
	"unicode/utf8"
)

// sendMessageLimit Лимит Telegram на длину текста одного сообщения
const sendMessageLimit = 4096

const fenceMarker = "```"

// splitMessage Разбивает текст на куски не длиннее limit рун по границам строк.
// Если разрез попадает внутрь блока ```, кусок закрывается маркером ```, а следующий открывается им заново,
// чтобы каждый кусок оставался валидным markdown.
func splitMessage(text string, limit int) []string {
	if utf8.RuneCountInString(text) <= limit {
		return []string{text}
	}

	// Любому куску может понадобиться закрытие блока "\n```", поэтому контент набирается с этим запасом
	fenceLen := utf8.RuneCountInString(fenceMarker)
	contentLimit := limit - fenceLen - 1

	lines := strings.Split(text, "\n")

	var (
		chunks  []string
		curLen  int // Длина собранного куска в рунах, каждая строка учитывается с завершающим "\n"
		inFence bool
	)
	curLines := make([]string, 0, len(lines))

	flush := func() {
		if len(curLines) == 0 {
			return
		}
		if inFence {
			curLines = append(curLines, fenceMarker)
		}
		chunks = append(chunks, strings.Join(curLines, "\n"))
		curLines = curLines[:0]
		curLen = 0
		if inFence {
			curLines = append(curLines, fenceMarker)
			curLen = fenceLen + 1
		}
	}

	for _, line := range lines {
		lineLen := utf8.RuneCountInString(line)

		// Строка, закрывающая блок, всегда помещается: запас contentLimit ровно её размера
		closesFence := inFence && line == fenceMarker
		if !closesFence {
			if curLen+lineLen > contentLimit {
				flush()
			}

			// Одиночная строка длиннее лимита — режем по рунам
			if curLen+lineLen > contentLimit {
				r := []rune(line)
				for curLen+len(r) > contentLimit {
					capacity := contentLimit - curLen
					curLines = append(curLines, string(r[:capacity]))
					curLen += capacity + 1
					flush()
					r = r[capacity:]
				}
				line, lineLen = string(r), len(r)
			}
		}

		curLines = append(curLines, line)
		curLen += lineLen + 1

		if strings.HasPrefix(line, fenceMarker) {
			inFence = !inFence
		}
	}
	flush()

	return chunks
}

// paginateTable Собирает страницы таблицы: каждая — самостоятельный блок ``` с заголовком, не длиннее limit рун.
// Строка длиннее лимита попадает на страницу целиком, такие страницы дорежет splitMessage при отправке
func paginateTable(header string, rows []string, limit int) []string {
	// Постоянная часть страницы: обрамление "```\n" и "\n```" плюс заголовок с завершающим "\n"
	base := utf8.RuneCountInString(header) + 2*(utf8.RuneCountInString(fenceMarker)+1)

	var pages []string
	pageRows := make([]string, 0, len(rows))
	pageLen := base

	flush := func() {
		if len(pageRows) == 0 {
			return
		}
		pages = append(pages, fenceMarker+"\n"+header+"\n"+strings.Join(pageRows, "\n")+"\n"+fenceMarker)
		pageRows = pageRows[:0]
		pageLen = base
	}

	for _, row := range rows {
		rowLen := utf8.RuneCountInString(row) + 1
		if pageLen+rowLen > limit {
			flush()
		}
		pageRows = append(pageRows, row)
		pageLen += rowLen
	}
	flush()

	return pages
}
