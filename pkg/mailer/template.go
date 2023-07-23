package mailer

type Mailer int

const (
	LOGIN Mailer = iota + 1
)

var mapTemplate = map[Mailer]string{
	LOGIN: loginTemplate,
}

var mapSubject = map[Mailer]string{
	LOGIN: "Login Verification",
}

var (
	loginTemplate = `<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
  <div style="margin:50px auto;width:70%;padding:20px 0">
    <div style="border-bottom:1px solid #eee">
      <p style="font-size:1.4em;color: #267adc;text-decoration:none;font-weight:600">Login Verification</p>
    </div>
    <p style="font-size:1.1em">Hi there, thank you for choosing dating apps,<br /> Below is the verification code for login, this verification code is valid for 5 minutes</p>
    <h2 style="background: #267adc;margin: 20px 10px 20px 0px;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;">{{ .Code}}</h2>
    <p style="font-size:0.9em;">Regards,<br />dating apps</p>
    <hr style="border:none;border-top:1px solid #eee" />
  </div>
</div>`
)
