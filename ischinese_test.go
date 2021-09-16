package ischinese

import (
	"testing"
)

func TestParseUnicodeString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    rune
		wantErr bool
	}{
		{
			s:    "U+3469",
			want: '\u3469',
		},
		{
			s:    "U+2966A",
			want: '\U0002966A',
		},
		{
			s:    "U+295E1",
			want: '\U000295E1',
		},
		{
			s:       "FFFFFFFFF",
			wantErr: true,
		},
		{
			s:       "G",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUnicodeString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUnicodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseUnicodeString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildDictionary(t *testing.T) {
	simplifiedDict := make(map[rune]struct{})
	traditionalDict := make(map[rune]struct{})
	err := buildDictionary(simplifiedDict, traditionalDict)
	if err != nil {
		t.Error(err)
	}
}

func TestIsPureSimplifiedChinese(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			s:    "",
			want: true,
		},
		{
			s:    "hello world",
			want: false,
		},
		{
			s:    "こんにちは世界",
			want: false,
		},
		{
			s:    "你很機車哎",
			want: false,
		},
		// source: 「明朝那些事儿」
		{
			s:    "在军队中，汤和算是个奇特的人，他在朱元璋刚参军时，已经是千户，但他却很尊敬朱元璋，在军营里，人们可以看到一个奇特的现象，官职高得多的汤和总是走在士兵朱元璋的后边，并且毫不在意他人的眼神，更奇特的是朱元璋似乎认为这是理所应当的事情，也没有推托过。",
			want: true,
		},
		{
			s:    "朱元璋奉命带兵攻击郭子兴的老家，定远，从这一点可以看出他的岳父实在存心不良，当时的定远有重兵看守，估计郭子兴让他去就是不想再看到活着的朱元璋，但朱元璋就是朱元璋，他找到了元军的一个缝隙，攻克了定远，然后在元军回援前撤出，此后，连续攻击怀远、安奉、含山、虹县，四战四胜，锐不可当！",
			want: true,
		},
		{
			s:    "【厉害的陈友谅】",
			want: true,
		},
		{
			s:    "喜欢锻炼的人，身体应该比较好，天天锻炼的人（比如运动员），就不一定好，旅游也是如此。",
			want: true,
		},
		{
			s:    "一、历史证明，叛徒是没有好下场的。同志瞧不起的人，敌人也瞧不起。",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPureSimplifiedChinese(tt.s); got != tt.want {
				t.Errorf("IsPureSimplifiedChinese() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPureTraditionalChinese(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		//{
		//	s:    "了", // TODO: both in simplified and traditional in real world, but not in traditionalDict
		//	want: true,
		//},
		{
			s:    "",
			want: true,
		},
		{
			s:    "hello world",
			want: false,
		},
		{
			s:    "2021",
			want: false,
		},
		{
			s:    "こんにちは世界",
			want: false,
		},
		{
			s:    "你很機車哎",
			want: true,
		},
		{
			s:    "【厉害的陈友谅】",
			want: false,
		},
		{
			s:    "喜欢锻炼的人，身体应该比较好，天天锻炼的人（比如运动员），就不一定好，旅游也是如此。",
			want: false,
		},
		// source: https://zh.wikipedia.org/wiki/%E5%B0%84%E9%B5%B0%E8%8B%B1%E9%9B%84%E5%82%B3
		{
			s:    "《射鵰英雄傳》小說前後一共有三個版本：連載版（舊版）、流行版（新版）、世紀修訂版（新修版）。",
			want: true,
		},
		{
			s:    "然而連載《射鵰英雄傳》期間，因為金庸在長城電影公司擔任編劇和導演，瑣事繁多，精力時有不殆，所以小說中有很多情節他本人並不是非常滿意。",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPureTraditionalChinese(tt.s); got != tt.want {
				t.Errorf("IsPureTraditionalChinese() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsChinese(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			s:    "",
			want: true,
		},
		{
			s:    "机车abc",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChinese(tt.s); got != tt.want {
				t.Errorf("IsChinese() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSimplifiedChinese(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			s:    "你很機車哎", // TODO:
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSimplifiedChinese(tt.s); got != tt.want {
				t.Errorf("IsSimplifiedChinese() = %v, want %v", got, tt.want)
			}
		})
	}
}
