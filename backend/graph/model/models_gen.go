// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

// 認証成功時のペイロード
type AuthPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	// 認証されたユーザー情報
	User *User `json:"user"`
}

// ユーザー作成時の入力データ
type CreateUserInput struct {
	UniqueName  string `json:"uniqueName"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type Mutation struct {
}

type Query struct {
}

// ユーザー更新時の入力データ
type UpdateUserInput struct {
	ID          string  `json:"id"`
	UniqueName  *string `json:"uniqueName,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
}

type User struct {
	ID          string `json:"id"`
	UniqueName  string `json:"uniqueName"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}
