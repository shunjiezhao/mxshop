package token

import (
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzKN3cFz2skKkStIJGnDA
In3MKwgJyyzivbeZUJhW8FCO2/w03rvggluxzWQX5bDgTiAAHscOcdyUTYcqIdjd
0YaUYX1ayi3FSUBnxpf/4eTYq4J9YDdI0yyzfOPJEkZmjXSSpbJ8A2IXZ9G5xTON
JT9Prkm9bWQbINove3v4G0JHnxj5d7nUPy8H7hzcfPyYPgPA3ryMA3B02ciuMX+A
M+qHsZpzPLTSMxjc10KTbRsViKylAJcguLr5Pw7GjXA8a0BbQnT2FJWcabtBTw/P
BycL3UNtzckQvcJoUfsIlaK9GD0AY3lnXcReVt9CRDJlp0GtgtnXzLLUphDPl/u1
/QIDAQAB
-----END PUBLIC KEY-----`

func TestJWTTokenVerifier_Verify(t *testing.T) {
	pem, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("can not parse public key : %v", err)
	}
	v := &JWTTokenVerifier{
		PublicKey: pem,
	}
	// 表单驱动测试
	cases := []struct {
		name       string
		token      string
		now        time.Time
		wantName   string
		wantUserId int32
		wantErr    bool
	}{
		{
			name:       "right",
			token:      "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuaWNrbmFtZSI6InpzaiIsInVzZXJfaWQiOjEyMywiZXhwIjoxNTE2MjM5MDI0LCJpYXQiOjE1MTYyMzkwMjIsImlzcyI6InpzaiIsIm5iZiI6MTUxNjIzOTAyMn0.kwYn68uic1mrH9fzMYgOJzi4C7MpojSHB-S1Rnj3tK22jhYbuDNwTQY-xIwan8WrwLaPBXeQVbwhxUaTMPhDF4pbyO0dHJ71SjLIjSpBcWJJe3wFMF7w9-aYkQdVIfauvQ4xGMcJaHR3-zCy4N2EkEZkMh9KpxmqEh84mfUt6kXUrYo2NPAHO63g_k3tBQOu3TiC_YeHJFQIwzo-O7ycXUQPAuADmQrSOG323yNM-XyZV0zvUaIriA7GcXQ-ckeU6c_cRlWZ4OjQu6M81GQKIEZrYjh9rmWHy54p3iOw7IZzQgQZ0BiFV-zccK-mD4DSUhDglGPGHyBCkZiJ6YaSPQ",
			now:        time.Unix(1516239022, 0),
			wantName:   "zsj",
			wantUserId: 123,
		},
		{
			name:    "wrong_exp",
			token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuaWNrbmFtZSI6InpzaiIsInVzZXJfaWQiOjEyMywiZXhwIjoxNTE2MjM5MDI0LCJpYXQiOjE1MTYyMzkwMjIsImlzcyI6InpzaiIsIm5iZiI6MTUxNjIzOTAyMn0.kwYn68uic1mrH9fzMYgOJzi4C7MpojSHB-S1Rnj3tK22jhYbuDNwTQY-xIwan8WrwLaPBXeQVbwhxUaTMPhDF4pbyO0dHJ71SjLIjSpBcWJJe3wFMF7w9-aYkQdVIfauvQ4xGMcJaHR3-zCy4N2EkEZkMh9KpxmqEh84mfUt6kXUrYo2NPAHO63g_k3tBQOu3TiC_YeHJFQIwzo-O7ycXUQPAuADmQrSOG323yNM-XyZV0zvUaIriA7GcXQ-ckeU6c_cRlWZ4OjQu6M81GQKIEZrYjh9rmWHy54p3iOw7IZzQgQZ0BiFV-zccK-mD4DSUhDglGPGHyBCkZiJ6YaSPQ",
			now:     time.Unix(1666263846, 0),
			wantErr: true,
		},
		{
			name:    "wrong_token",
			token:   "e.asdfsadf.NO63dFEJWO6b5pIOzq0Qh1vpPXvtROGtpeWRo5vpvwEuwkrqlPMbAaufBXh5_xitiuTCBGUkEY1nbHlhhc95aHuHFcHrl4hJ08tdS_JjA-SjAC71bLyACGM02o3CRhwnQd_Met0QTCAiBu6kIDG3qKF64PDMkwi7SGO7GkQs0xGkkAtWwzqH8WLfl7sj0d0ejTA0pakATcP0Lu_trvalGHy1xLlW6nfsQZrNK9MKGbAYxWZcfrhZS7MeBJqZcNQZuMdMtcwVfbh1EIf3kTTQ0bIaKvA9lkU1Q6xhUQGphXCUu9-87UIPRFvuyYB4JnpWu1_8pSAJ5-h9lcO4Vv-b1g",
			now:     time.Unix(1666263702, 0),
			wantErr: true,
		},
		{
			name:    "wrong_sign",
			token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuaWNrbmFtZSI6InpzaiIsInVzZXJfaWQiOjEyMywiZXhwIjoxNTE2MjM5MDI0LCJpYXQiOjE1MTYyMzkwMjIsImlzcyI6InpzaiIsIm5iZiI6MTUxNjIzOTAyMn0.kwYn68uic1mrH9fzMYgOJzi4C7MpojSHB-S1Rnj3tK22jhYbuDNwTQY-xIwan8WrwLaPBXeQVbwhxUaTMPhDF4pbyO0dHJ71SjLIjSpBcWJJe3wFMF7w9-aYkQdVIfauvQ4xGMcJaHR3-zCy4N2EkEZkMh9KpxmqEh84mfUt6kXUrYo2NPAHO63g_k3tBQOu3TiC_YeHJFQIwzo-O7ycXUQPAuADmQrSOG323yNM-XyZV0zvUaIriA7GcXQ-ckeU6c_cRlWZ4OjQu6M81GQKIEZrYjh9rmWHy54p3iOw7IZzQgQZ0BiFV-zccK-mD4DSUhDglGPGHyBCkZiJ6Yasda",
			now:     time.Unix(1666256501, 0),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			claim, err := v.Verify(c.token)
			if c.wantErr {
				if err == nil {
					t.Errorf("want eror; got no errorv")
				}
				return
			}
			if err != nil {
				t.Errorf("Verify failed : %v", err)
			}
			if claim.UserId != c.wantUserId || claim.Nickname != c.wantName {
				t.Errorf("%s: work Verify token error\n want:%q %d \b but:%q %d", c.name,
					c.wantName, c.wantUserId,
					claim.Nickname, claim.UserId)
			}
		})

	}
}

const privateKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAhaU6EAscXZqGNXaV1644zOrPNdVI2S2buIPWP7CVC6qnyMH2
kdIfNHuLR9plxcAmtyjfwjYM+/gGj/V80M7p1TZaLGy9J18rfRPylNBi/uzRE0ty
P9aPWg30FYHC/NhVE7cvkldZ9O695lIp7HATAc9AOiVQ2L/o81/DuYiBJev4V6Xb
f339NzukEsQ7lH8lELExQgRzHjjz8gDyYIZC9Gj1yCZDyZOIKxkIMHVxLbiyhGRX
V+WPZinqWj4HzT4+CpxBG7aXxWPNiAX7VAIj4aUWfPzF3DK2NNTkp42bdJMiVn1x
SvwJ963BYTcVnUqQvoBwnbkJOy4EdadR6sMuUQIDAQABAoIBAAqr28wGJxiuSwMf
2e0j4zMxmKQSyMNWqfV8yXHCdtQ8pzMXkcmo/obhSojNPo0gv7amU9+dE0JXVqbO
Ek5WB0PYGDEt6ZioN7/ABQGu9gim/jbNMm04g6/MJ25EMz6fQ4PUf77urKcrHQy0
CNbxSAM/+j/qVN9Jy58LSYSaCr3tUPdx7sZlPZRk1EHUM6zAP2X6bApAt4EsrnDU
fuZQQuKurbPsXnL4q6tJ7rsLIoe90A5yUeYJvTS3aivP5WYNHIe3jUf+S2YKNz1L
NOsCEriBAaGkrqsJy8dM/VC70njBgX2FuFSvbd+KPl+ecQLAk1oOvI/5ZhyzrDKY
R/Ejn3ECgYEAyvv7bUpYCLxoAbXRfBCejIriSHPQNrSg4D6wZFu+TTMRD6Oeetd7
FeFRlzHHz9jG8C1rlEiw4R9nGqnah4OuSggQgH1kZdTEj+p+7iXUD5ov5jbtfnJo
naZ8BrUq81I43MP4zeJrAwiBD7aDp9XPeQfbQpEOOdQ+JVwsHx0iMSUCgYEAqI0U
hyQTBihsxZo5EiTUNi6LDlfgsSRDzlr3uoulgqYEqqzRour/BoTL5GYmupSrtxNL
6Nt3P9dKRvIV2rHV4Xo/+ZV8vJBRGe1xBuh4D1ve8pov8Uc6nyGPSP99QmaMc7Pr
n20UNDxp/tIjxJpyCYtlNfCWv6UTihAdMj6bbr0CgYBhF67oVAtQAm7tgn61jW2J
ZFVguqT5xeS93r6ZAplAbBDZHjaMI84oZSKV46Xj8ZkXAWLYBv00ccTrqBtzfrU7
jCf4jgIcA24SOOSGHWoWHHaU8+kd9rO71Qq2Wqo0wTuZvdOhB5CQXtz9GxxWh5s6
FVv3t1LKro8bZ79jEphsUQKBgDKs+BcJiY64aLugergxynvf0n8lfLDFrn1EbGbx
xXlaYNzPyNeqv7I+Cu9IpyxBtr78Vj5Ufa38FKDv+BIglaWNE97+StqGqVuaP/lL
u40img1mvjNUrxNZC7Nu3UIxgtjmp1jveruZzmSG2aoqpU6pUmy9QRWtlApWffC1
UhYhAoGAK8PVsp0IFSpMJXTZqoMUS7gRMM9RhxXiGLUoeZeD8mWPa6VYGWrchqHR
/wsI3ioDgyfJQEmF/575+lS1RUbL8V5lcTGwwfi7l7V4DeJehL0TxZbkFGwd1TlO
O6OcQQSYGlYKw4vNJPxYCEm9Pp9eUFVT9+QTS3PjxZlQD48ljwQ=
-----END RSA PRIVATE KEY-----`

// 小坑 首先如果用rsa 加密 用 rs256 加密方法
// 然后需保证data数据顺序一致。
func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("create: parse key: %v", err)
	}

	tk := NewJWTokenGen("zsj", 2*time.Second, key)
	tk.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}

	token, err := tk.GenerateToken("zsj", 123, 0)
	if err != nil {
		t.Errorf("can not GenerateToken err: %v", err)
	}
	want := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJuaWNrbmFtZSI6InpzaiIsInVzZXJfaWQiOjEyMywiZXhwIjoxNTE2MjM5MDI0LCJpYXQiOjE1MTYyMzkwMjIsImlzcyI6InpzaiIsIm5iZiI6MTUxNjIzOTAyMn0.HiXgzd_ZiNziIExmXD3sUxZNv6BV0JLj7IGwTECT07zdaqng1eZ-iXLeNjO-J9snlH9eybmPn8NQkLt7go8HeC0IzxM63vz6cXbnRgf3kGrRRycm06V-vVOb9X4bZWsylCbumXuQdjP70qG9K2qfoacvN3nL6wYQCG8GxKg0n1ii1c_CgUajv0hU2SAFS9TkRWyKZNChPgxe3vV0bgV_AA89_N_TrwAzj3he-JG-hcVS5nFsXiy5VmcP65JocjM0rPDQcn0cMq8x24IY_8jyRMQuPItVf1ccsRO67BgaXKlUvWE9Yf54y-LC60RWIM2MuAu_I4iMZOvDl25LF0WPQw"
	if token != want {
		t.Errorf("work generate token error\n want:%q \b but:%q", want, token)
	}
}
