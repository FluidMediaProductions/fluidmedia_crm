package main

import (
	"github.com/fluidmediaproductions/fluidmedia_crm/model"
	"net/http"
	"log"
	"github.com/pquerna/otp/totp"
	"bytes"
	"image/png"
	"encoding/base64"
	"fmt"
	"html/template"
)

type ProfileContext struct {
	User *model.User
}

type Profile2FAContext struct {
	User *model.User
	Secret string
	QR template.URL
}

func handleProfile(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		displayWithContext(w, "profile", page, user, &ProfileContext{User: user})
	} else if r.Method == "POST" {
		r.ParseForm()
		newUser := &model.User{
			ID: user.ID,
			Name: r.Form.Get("name"),
			Email: r.Form.Get("email"),
			Phone: r.Form.Get("phone"),
			IsAdmin: user.IsAdmin,
			Login: user.Login,
			Pass: r.Form.Get("pass"),
			Disabled: user.Disabled,
			TotpSecret: user.TotpSecret,
		}
		err := m.SaveUser(newUser)
		if err != nil {
			log.Printf("Error updating user: %v", err)
			display500(w, err)
			return
		}
		http.Redirect(w, r, "/profile", 302)
	}
}

func handleProfile2FA(m *model.Model, page *Page, user *model.User, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if user.TotpSecret == "" {
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "FluidMedia CRM",
				AccountName: user.Login,
			})
			if err != nil {
				display500(w, err)
				return
			}

			var buf bytes.Buffer
			img, err := key.Image(200, 200)
			png.Encode(&buf, img)

			byt := buf.Bytes()
			e64 := base64.StdEncoding
			maxEncLen := e64.EncodedLen(len(byt))
			enc := make([]byte, maxEncLen)
			e64.Encode(enc, byt)

			context := Profile2FAContext{
				QR:     template.URL(fmt.Sprintf("data:image/png;base64,%s", enc)),
				Secret: key.Secret(),
				User:   user,
			}

			displayWithContext(w, "profile2fa", page, user, &context)
		} else {
			newUser := &model.User{
				ID: user.ID,
				Name: user.Name,
				Email: user.Email,
				Phone: user.Phone,
				IsAdmin: user.IsAdmin,
				Login: user.Login,
				Pass: "",
				Disabled: user.Disabled,
				TotpSecret: "",
			}
			err := m.SaveUser(newUser)
			if err != nil {
				log.Printf("Error updating user: %v", err)
				display500(w, err)
				return
			}
			http.Redirect(w, r, "/profile", 302)
		}
	} else if r.Method == "POST" {
		r.ParseForm()
		secret := r.Form.Get("secret")
		verify := r.Form.Get("verify")

		valid := totp.Validate(verify, secret)

		if valid {
			newUser := &model.User{
				ID: user.ID,
				Name: user.Name,
				Email: user.Email,
				Phone: user.Phone,
				IsAdmin: user.IsAdmin,
				Login: user.Login,
				Pass: "",
				Disabled: user.Disabled,
				TotpSecret: secret,
			}
			err := m.SaveUser(newUser)
			if err != nil {
				log.Printf("Error updating user: %v", err)
				display500(w, err)
				return
			}
			http.Redirect(w, r, "/profile", 302)
		} else {
			http.Redirect(w, r, "/profile2fa", 302)
		}
	}
}
