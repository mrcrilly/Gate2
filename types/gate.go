
package types

import (
    "encoding/base32"
    "encoding/base64"
    "io/ioutil"
    "fmt"

    otp "github.com/dgryski/dgoogauth"
    "code.google.com/p/rsc/qr"
    "github.com/AutoLogicTechnology/Gate2/helpers"
)

const (
    GATE_WINDOW_SIZE = 0
    GATE_HOTP_COUNTER = 0
)

type Gate struct {
    UserID string 
    UserSecret string // raw version, for printing/debugging

    OTP *otp.OTPConfig
    QRCode string 
}

func NewGate (userid string) (g *Gate) {
    usersecret := helpers.RandomString()
    b32 := base32.StdEncoding.EncodeToString([]byte(usersecret))

    g = &Gate{
        UserID: userid,
        UserSecret: usersecret,
        OTP: &otp.OTPConfig{
            Secret: b32,
            WindowSize: GATE_WINDOW_SIZE,
            HotpCounter: GATE_HOTP_COUNTER,
        },
    }
    
    code, _ := qr.Encode(g.OTP.ProvisionURI(g.UserID), qr.Q)
    g.QRCode = base64.StdEncoding.EncodeToString(code.PNG())

    return g
}

func (g *Gate) WritePng () {
    q, _ := base64.StdEncoding.DecodeString(g.QRCode)
    ioutil.WriteFile(fmt.Sprintf("%s.png", g.UserID), q, 0644)
}

func (g *Gate) CheckCode (password string) (result bool, err error) {
    result, err = g.OTP.Authenticate(password)

    if err != nil {
        return false, err 
    }

    return result, nil
}