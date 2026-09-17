// Package dto mirrors the Kotlin data classes under Dto/*.kt one-to-one.
// JSON field names match the default Jackson serialization of the Kotlin
// camelCase property names exactly, so the wire format is unchanged.
package dto

// RegisterRequest mirrors Dto/AuthDto.kt RegisterRequest.
//
// NOTE: the Kotlin file imports jakarta.validation annotations (@Email,
// @NotBlank, @Size) but never actually applies them to the RegisterRequest
// fields, and there's no @Valid on the controller parameter either — so no
// validation is actually enforced in the original. We preserve that (i.e.
// we do NOT add validation here) to keep behavior identical.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Name         string `json:"name"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
