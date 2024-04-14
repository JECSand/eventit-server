package auth

import (
	"context"
	"github.com/JECSand/eventit-server/src/shared/enums"
	"github.com/JECSand/eventit-server/src/shared/routers"
	"net/http"
)

type ctxKey int

const (
	ctxClaims ctxKey = iota
)

// ClaimsFromCtx retrieves the parsed AppClaims from request context.
func ClaimsFromCtx(ctx context.Context) AppClaims {
	return ctx.Value(ctxClaims).(AppClaims)
}

// Authenticator inputs the route handler function along with User roleType to verify User token and permissions
func Authenticator(roleType enums.Role, next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	var errorObject routers.JWTError
	decodedToken, err := DecodeJWT(r.Header.Get("Auth-Token"))
	if err != nil {
		errorObject.Message = err.Error()
		routers.RespondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
	ctx := context.WithValue(r.Context(), ctxClaims, decodedToken)
	if roleType == enums.ROOT && decodedToken.Role == enums.ROOT {
		next.ServeHTTP(w, r.WithContext(ctx))
	} else if roleType == enums.ADMIN && decodedToken.Role == enums.ADMIN || decodedToken.Role == enums.ROOT {
		next.ServeHTTP(w, r.WithContext(ctx))
	} else if roleType == enums.MEMBER {
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	errorObject.Message = "Invalid Token"
	routers.RespondWithError(w, http.StatusUnauthorized, errorObject)
	return
}

// VerifyRootMiddleWare is used to verify that the requester is a valid admin
func VerifyRootMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Authenticator(enums.ROOT, next, w, r)
		return
	}
}

// VerifyAdminMiddleWare is used to verify that the requester is a valid admin
func VerifyAdminMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Authenticator(enums.ADMIN, next, w, r)
		return
	}
}

// VerifyMemberMiddleWare is used to verify that a requester is authenticated
func VerifyMemberMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Authenticator(enums.MEMBER, next, w, r)
		return
	}
}
