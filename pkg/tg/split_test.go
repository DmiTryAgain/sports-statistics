package tg

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

// buildFencedTable строит текст в духе buildTableByStat: таблица из rows строк внутри ```-блока
func buildFencedTable(rows int) string {
	var b strings.Builder
	b.WriteString("```\n")
	b.WriteString("Упражнение    Вес     Кол-во    Подходы\n")
	for i := range rows {
		fmt.Fprintf(&b, "подтягивания    %dкг    %d    %d\n", 20+i%50, 100+i, 5+i%10)
	}
	b.WriteString("```")
	return b.String()
}

// contentLines возвращает строки текста без маркеров ```
func contentLines(text string) []string {
	var res []string
	for _, line := range strings.Split(text, "\n") {
		if !strings.HasPrefix(line, "```") {
			res = append(res, line)
		}
	}
	return res
}

func assertChunksValid(t *testing.T, chunks []string, limit int) {
	t.Helper()
	if len(chunks) == 0 {
		t.Fatal("splitMessage returned no chunks")
	}
	for i, c := range chunks {
		if got := utf8.RuneCountInString(c); got > limit {
			t.Errorf("chunk %d has %d runes, want <= %d", i, got, limit)
		}
		if c == "" {
			t.Errorf("chunk %d is empty", i)
		}
		if strings.Count(c, "```")%2 != 0 {
			t.Errorf("chunk %d has unbalanced code fences:\n%s", i, c)
		}
	}
}

func TestSplitMessageShortText(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		limit int
	}{
		{"short text", "привет", 100},
		{"exactly at limit", strings.Repeat("а", 10), 10},
		{"short fenced table", buildFencedTable(3), 4096},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitMessage(tt.text, tt.limit)
			if len(got) != 1 || got[0] != tt.text {
				t.Errorf("splitMessage() = %d chunks %q, want the original text as a single chunk", len(got), got)
			}
		})
	}
}

func TestSplitMessagePlainText(t *testing.T) {
	lines := make([]string, 40)
	for i := range lines {
		lines[i] = fmt.Sprintf("строка номер %d", i)
	}
	text := strings.Join(lines, "\n")
	const limit = 100

	chunks := splitMessage(text, limit)

	assertChunksValid(t, chunks, limit)
	if len(chunks) < 2 {
		t.Fatalf("got %d chunks, want at least 2", len(chunks))
	}
	if joined := strings.Join(chunks, "\n"); joined != text {
		t.Errorf("joined chunks differ from original text:\ngot:\n%s\nwant:\n%s", joined, text)
	}
}

func TestSplitMessageFencedTable(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		limit int
	}{
		{"long table", buildFencedTable(200), 300},
		{"table with prefix line", "Нераспознанные периоды: вчера, сегодня\n" + buildFencedTable(100), 300},
		{"table near real limit", buildFencedTable(300), sendMessageLimit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := splitMessage(tt.text, tt.limit)

			assertChunksValid(t, chunks, tt.limit)
			if len(chunks) < 2 {
				t.Fatalf("got %d chunks, want at least 2", len(chunks))
			}

			// Ни одна строка данных не потеряна и не разрезана
			var got []string
			for _, c := range chunks {
				got = append(got, contentLines(c)...)
			}
			want := contentLines(tt.text)
			if len(got) != len(want) {
				t.Fatalf("got %d content lines, want %d", len(got), len(want))
			}
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("content line %d = %q, want %q", i, got[i], want[i])
				}
			}

			// Каждый кусок с содержимым таблицы — валидный markdown code-block
			for i, c := range chunks[1:] {
				if !strings.HasPrefix(c, "```\n") {
					t.Errorf("continuation chunk %d does not open a code fence:\n%s", i+1, c)
				}
			}
			if last := chunks[len(chunks)-1]; !strings.HasSuffix(last, "```") {
				t.Errorf("last chunk does not close the code fence:\n%s", last)
			}
		})
	}
}

func TestSplitMessageOverlongLine(t *testing.T) {
	text := strings.Repeat("я", 250)
	const limit = 100

	chunks := splitMessage(text, limit)

	assertChunksValid(t, chunks, limit)
	if joined := strings.Join(chunks, ""); joined != text {
		t.Errorf("joined chunks differ from original text: got %d runes, want %d", utf8.RuneCountInString(joined), utf8.RuneCountInString(text))
	}
}

func TestPaginateTableSinglePage(t *testing.T) {
	got := paginateTable("Заголовок", []string{"строка 1", "строка 2"}, 100)
	want := []string{"```\nЗаголовок\nстрока 1\nстрока 2\n```"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("paginateTable() = %q, want %q", got, want)
	}
}

func TestPaginateTableManyRows(t *testing.T) {
	header := "Упражнение    Кол-во"
	rows := make([]string, 100)
	for i := range rows {
		rows[i] = fmt.Sprintf("подтягивания    %d", i)
	}
	const limit = 200

	pages := paginateTable(header, rows, limit)

	if len(pages) < 2 {
		t.Fatalf("got %d pages, want at least 2", len(pages))
	}
	var got []string
	for i, p := range pages {
		if runes := utf8.RuneCountInString(p); runes > limit {
			t.Errorf("page %d has %d runes, want <= %d", i, runes, limit)
		}
		lines := strings.Split(p, "\n")
		if lines[0] != "```" || lines[1] != header || lines[len(lines)-1] != "```" {
			t.Errorf("page %d is not a fenced table with header:\n%s", i, p)
		}
		got = append(got, lines[2:len(lines)-1]...)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Errorf("rows across pages differ from input: got %d rows, want %d", len(got), len(rows))
	}
}

func TestPaginateTableOversizedRow(t *testing.T) {
	rows := []string{strings.Repeat("я", 300), "обычная строка"}

	pages := paginateTable("Заголовок", rows, 100)

	// Сверхдлинная строка не теряется: её страница выходит за лимит и будет дорезана splitMessage при отправке
	var got []string
	for _, p := range pages {
		lines := strings.Split(p, "\n")
		got = append(got, lines[2:len(lines)-1]...)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Errorf("rows across pages differ from input: got %q, want %q", got, rows)
	}
}
