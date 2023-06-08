package auth

func verificationEmailContent(email, token string) string {
	return `
	<p>Selamat datang di Gimsak</p>
	<p style="text-align: justify">Kami telah menerima pendaftaran akun Kamu, segera <b>verifikasi e-mail</b> dengan menekan tombol di bawah.</p>
	<table role="presentation" border="0" cellpadding="0" cellspacing="0" class="btn btn-primary">
	  <tbody>
		<tr>
		  <td align="center">
			<table role="presentation" border="0" cellpadding="0" cellspacing="0">
			  <tbody>
				<tr>
				  <td> <a href="https://gimsak.com/auth/verify?email=` + email + `&token=` + token + `" target="_blank">Verifikasi Email</a> </td>
				</tr>
			  </tbody>
			</table>
		  </td>
		</tr>
	  </tbody>
	</table>
	<p style="text-align: justify">Apabila verifikasi bermasalah, Kamu bisa meminta ulang e-mail konfirmasi <a href="https://gimsak.com/auth/resend-verify?email=` + email + `">di sini</a>.</p>
	<p style="text-align: justify">Segala bentuk informasi seperti nomor kontak, alamat e-mail, atau password Anda bersifat rahasia. Jangan menginformasikan data-data tersebut kepada siapapun, termasuk kepada pihak yang mengatasnamakan Gimsak.</p>	
	`
}
