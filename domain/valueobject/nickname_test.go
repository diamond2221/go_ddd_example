package valueobject

import (
	"testing"
)

func TestNewNickname(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError error
	}{
		// 合法的昵称
		{
			name:      "合法昵称：纯中文",
			input:     "张三李四",
			wantError: nil,
		},
		{
			name:      "合法昵称：纯英文",
			input:     "Alice",
			wantError: nil,
		},
		{
			name:      "合法昵称：纯数字",
			input:     "123456",
			wantError: nil,
		},
		{
			name:      "合法昵称：中文+数字",
			input:     "张三123",
			wantError: nil,
		},
		{
			name:      "合法昵称：英文+数字",
			input:     "Alice123",
			wantError: nil,
		},
		{
			name:      "合法昵称：中文+英文+数字",
			input:     "张三Alice123",
			wantError: nil,
		},
		{
			name:      "合法昵称：最短长度（3个字符）",
			input:     "张三李",
			wantError: nil,
		},
		{
			name:      "合法昵称：最长长度（16个字符）",
			input:     "这是十六个字符的昵称测试啊",
			wantError: nil,
		},

		// 非法的昵称：长度问题
		{
			name:      "非法昵称：太短（2个字符）",
			input:     "张三",
			wantError: ErrNicknameTooShort,
		},
		{
			name:      "非法昵称：太短（1个字符）",
			input:     "A",
			wantError: ErrNicknameTooShort,
		},
		{
			name:      "非法昵称：太长（17个字符）",
			input:     "这是超过十六个字符的昵称测试啊",
			wantError: ErrNicknameTooLong,
		},

		// 非法的昵称：字符格式问题
		{
			name:      "非法昵称：包含特殊字符@",
			input:     "张三@123",
			wantError: ErrNicknameInvalidFormat,
		},
		{
			name:      "非法昵称：包含空格",
			input:     "张三 李四",
			wantError: ErrNicknameInvalidFormat,
		},
		{
			name:      "非法昵称：包含下划线",
			input:     "zhang_san",
			wantError: ErrNicknameInvalidFormat,
		},
		{
			name:      "非法昵称：包含emoji",
			input:     "张三😀",
			wantError: ErrNicknameInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nickname, err := NewNickname(tt.input)

			// 检查错误是否符合预期
			if tt.wantError != nil {
				if err != tt.wantError {
					t.Errorf("NewNickname() error = %v, wantError %v", err, tt.wantError)
				}
				return
			}

			// 检查是否成功创建
			if err != nil {
				t.Errorf("NewNickname() unexpected error = %v", err)
				return
			}

			// 检查值是否正确
			if nickname.Value() != tt.input {
				t.Errorf("nickname.Value() = %v, want %v", nickname.Value(), tt.input)
			}
		})
	}
}

func TestNickname_Equals(t *testing.T) {
	nickname1, _ := NewNickname("张三123")
	nickname2, _ := NewNickname("张三123")
	nickname3, _ := NewNickname("李四456")

	if !nickname1.Equals(nickname2) {
		t.Error("相同值的昵称应该相等")
	}

	if nickname1.Equals(nickname3) {
		t.Error("不同值的昵称不应该相等")
	}
}

func TestNickname_Length(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantLength int
	}{
		{
			name:       "纯中文",
			input:      "张三李四",
			wantLength: 4,
		},
		{
			name:       "纯英文",
			input:      "Alice",
			wantLength: 5,
		},
		{
			name:       "中文+英文+数字",
			input:      "张三Alice123",
			wantLength: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nickname, _ := NewNickname(tt.input)
			if got := nickname.Length(); got != tt.wantLength {
				t.Errorf("nickname.Length() = %v, want %v", got, tt.wantLength)
			}
		})
	}
}
