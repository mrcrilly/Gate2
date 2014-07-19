
package gate

import (
    "encoding/base32"
    "encoding/base64"
    "io/ioutil"
    "regexp"

    otp "github.com/dgryski/dgoogauth"
    "code.google.com/p/rsc/qr"
    // "github.com/AutoLogicTechnology/Gate2/helpers"
)

const (
    GATE_WINDOW_SIZE = 0
    GATE_HOTP_COUNTER = 0
)

type Gate struct {
    UserID string 
    UserSecret string 
    ScratchCodes []string 

    OTP *otp.OTPConfig
    QRCode string 
}

func NewGate (userid string) (g *Gate) {
    usersecret := NewSecretCode()

    g = &Gate{
        UserID: userid,
        UserSecret: usersecret,
        OTP: &otp.OTPConfig{
            Secret: base32.StdEncoding.EncodeToString([]byte(usersecret)),
            WindowSize: GATE_WINDOW_SIZE,
            HotpCounter: GATE_HOTP_COUNTER,
        },
    }

    for i := 1; i <= 3; i ++ {
        g.ScratchCodes = append(g.ScratchCodes, NewScratchCode())
    }

    return g
}

func NewGateAndQRCode (userid string) (g *Gate) {
    g = NewGate(userid)

    code, _ := qr.Encode(g.OTP.ProvisionURI(g.UserID), qr.Q)
    g.QRCode = base64.StdEncoding.EncodeToString(code.PNG())

    return g 
}

func NewGateWithCustomSecret (userid, usersecret string) (g *Gate) {
    g = NewGateAndQRCode(userid)
    g.UserSecret = usersecret

    b32 := base32.StdEncoding.EncodeToString([]byte(usersecret))
    g.OTP.Secret = b32 

    return g 
}

func IsValidUserId (userid string) (bool) {
    r := regexp.MustCompile("^[a-zA-Z0-9._@-]+$")
    return r.MatchString(userid)
}

func IsValidTOTPCode (totpcode string) (bool) {
    r := regexp.MustCompile("^[0-9]{6}$")
    return r.MatchString(totpcode)
}

func (g *Gate) WritePngToFile (filename string) (err error) {
    q, err := base64.StdEncoding.DecodeString(g.QRCode)
    if err != nil {
        return err
    }

    ioutil.WriteFile(filename, q, 0644)
    return nil 
}

func (g *Gate) CheckCode (password string) (result bool, err error) {
    result, err = g.OTP.Authenticate(password)

    if err != nil {
        return false, err 
    }

    return result, nil
}